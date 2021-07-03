package main

import (
	"net/http"
	"strings"
)

func (s *server) routes() {
	s.router.HandleFunc("/", s.RoutingHandler())
}

// RoutingHandler routes requests to the correct handler
// we use this instead of the usual http.RouteMux because of the unique
// requirements of git HTTP URL convention that matches action using
// URL's suffix.
//
// Examples:
//
//   https://git.mydomain.com/parent-group/child-group/repo.git/info/refs
//   https://github.com/org-a/my-team/project.git/git-receive-pack
//   https://gitlab.com/org-b/team-10/foo-bar.git/git-upload-pack
//
func (s *server) RoutingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var handler http.HandlerFunc

		switch {
		case strings.HasSuffix(r.URL.Path, "/info/refs"):
			handler = s.InfoRefHandler()
		case strings.HasSuffix(r.URL.Path, "/git-receive-pack"):
			handler = s.ReceivePackHandler()
		case strings.HasSuffix(r.URL.Path, "/git-upload-pack"):
			handler = s.UploadPackHandler()
		default:
			handler = http.NotFound
		}

		handler(w, r)
	}
}
