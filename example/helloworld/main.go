package main

import "gmx"

func init() {
	gmx.Publish("hello", func() interface{} {
		return "world"
	})
}

func main() {
	// sleep forever
	select {}
}
