package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s := newServer()
	http.HandleFunc("/", s.ProxyHandler())

	return http.ListenAndServe("localhost:8080", nil)
}
