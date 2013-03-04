package main

import "github.com/davecheney/gmx"

func init() {
	gmx.Registry("gmx.example")("hello", func() interface{} {
		return "world"
	})
}

func main() {
	// sleep forever
	select {}
}
