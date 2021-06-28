package main

import "net/http"

type server struct {
	router *http.ServeMux
}

func newServer() *server {
	s := &server{}
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
