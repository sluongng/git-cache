package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestProxyHandler(t *testing.T) {
	is := is.New(t)
	s := newServer()

	r := httptest.NewRequest("GET", "/git/git.git/info/refs?service=git-upload-pack", nil)
	r.Header.Add("Git-Protocol", "version=2")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusOK)
}

func TestProxyHeader(t *testing.T) {
	is := is.New(t)
	s := newServer()

	r := httptest.NewRequest("GET", "/gitlab-org/gitlab.git/info/refs?service=git-upload-pack", nil)
	r.Header.Add("Git-Protocol", "version=2")
	r.Header.Add(UpstreamHeader, GitlabUpstream)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
}
