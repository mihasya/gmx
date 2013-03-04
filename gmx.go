/*
	Exposes arbitrary functions over a JSON HTTP interface. Inspired by
	JMX on the JVM platform, minus the proprietary protocol

	GMXLISTENON environment variable is used to determine the interface
	and port to listen on; default is localhost:41441.
*/
package gmx

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"fmt"
	"strings"
)

const GMX_VERSION = 0

var (
	nr = &nestedRegistry{
		regs:	make(map[string]map[string]func() interface{}),
	}
)

func init() {
	listenOn := os.Getenv("GMXLISTENON")
	if listenOn == "" {
		listenOn = "localhost:41441"
	}
	go func() {
		err := http.ListenAndServe(listenOn, nr)
		if err != nil {
			panic(err)
		}
	}()
}

type nestedRegistry struct {
	sync.Mutex
	regs		map[string]map[string]func() interface{}
}

func (nr *nestedRegistry) registry(name string) func(key string, f func() interface{}) {
	nr.Lock()
	defer nr.Unlock()
	r, ok := nr.regs[name]
	if !ok {
		r = make(map[string]func() interface{})
		nr.regs[name] = r
	}
	return func(key string, f func() interface{}) {
		r[key] = f
	}
}

// TODO: abstract this out, break it up
// TODO: allow some sort of protocol versioning; perhaps take advantage of
// ServeMux's inexact matching to direct to sub-handlers
func (nr *nestedRegistry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header()["Content-type"] = []string{"application/json"}
	path := req.URL.Path
	if path[0] == '/' {
		path = path[1:]
	}
	parts := strings.Split(path, "/")
	if len(path) == 0 {
		keys := make([]string, 0)
		for k := range nr.regs {
			keys = append(keys, k)
		}
		content, _ := json.Marshal(keys)
		w.Write(content)
		return
	} else if len(parts) == 1 {
		name := parts[0]
		r, ok := nr.regs[name]
		if !ok {
			http.Error(w, fmt.Sprintf("No gmx registry found with name '%v'", name), 404)
			return
		}
		keys := make([]string, 0)
		for k := range r {
			keys = append(keys, k)
		}
		content, _ := json.Marshal(keys)
		w.Write(content)
		return

	} else if len(parts) == 2 {
		name := parts[0]
		r, ok := nr.regs[name]
		if !ok {
			http.Error(w, fmt.Sprintf("No gmx registry found with name '%v'", name), 404)
			return
		}
		switch req.Method {
			case "GET": {
				f, ok := r[parts[1]]
				if !ok {
					http.Error(w, fmt.Sprintf("No function registered at '%v'", path), 404)
					return
				}
				result := f()
				content, _ := json.Marshal(result)
				w.Write([]byte(string(content)))
				return
			}
			default: {
				http.Error(w, "Only GETs supported for now", 400)
				return
			}
		}
	} else {
		http.Error(w, fmt.Sprintf("Invalid path '%v'", path), 400)
	}

}

func Registry(name string) func(key string, f func() interface{}) {
	return nr.registry(name)
}

// keep for now to avoid having to figure out how instrument should work
func Publish(name string, f func()interface{}) {}
