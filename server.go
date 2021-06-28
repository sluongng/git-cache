package main

import "net/http"

type server struct {
	router              http.ServeMux
	defaultUpstreamHost string
}

func newServer() *server {
	s := &server{}
	s.routes()
	s.defaultUpstreamHost = GithubUpstream
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
