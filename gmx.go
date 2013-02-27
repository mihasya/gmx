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
)

const GMX_VERSION = 0

var (
	r = &registry{
		entries: make([]string, 0),
		hdl:     http.NewServeMux(),
	}
)

func init() {
	r.register("/", func() interface{} {
		return r.entries
	})
	listenOn := os.Getenv("GMXLISTENON")
	if listenOn == "" {
		listenOn = "localhost:41441"
	}
	go func() {
		err := http.ListenAndServe(listenOn, r.hdl)
		if err != nil {
			panic(err)
		}
	}()
}

// Publish registers the function f with the supplied key.
func Publish(key string, f func() interface{}) {
	r.register(key, f)
}

type registry struct {
	sync.Mutex // protects entries from concurrent mutation
	entries    []string
	hdl        *http.ServeMux
}

func (r *registry) register(key string, f func() interface{}) {
	r.Lock()
	defer r.Unlock()
	if key[0] != '/' {
		key = "/" + key
	}
	r.entries = append(r.entries, key)
	r.hdl.HandleFunc(key, func(response http.ResponseWriter, req *http.Request) {
		content, _ := json.Marshal(f())
		response.Write([]byte(string(content)))
	})
}
