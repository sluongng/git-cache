package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	UpstreamHeader = "Git-Cache-Upstream"
)

// InfoRefHandler handles the incoming GET /info/refs requests which often initiated by
// all git activities that interacts with remotes.
// The operations is meant so that the client can retrieve the capabilities of remote
// to decide what could/should be used in the follow up operations.
//
// Current implementation is to forward all of these requests upstream.
func (s *server) InfoRefHandler() http.HandlerFunc {
	log.Println("InfoRefHandler")
	return s.proxyHandler()
}

// ReceivePackHandler handles the incoming POST requests which often initiated by
// git operations such as:
//   - git-push
func (s *server) ReceivePackHandler() http.HandlerFunc {
	log.Println("ReceivePackHandler")
	return s.proxyHandler()
}

// UploadPackHandler handles the incoming POST requests which often initiated by
// git operations such as:
//   - git-fetch
//   - git-clone
//   - git-pull
//   - git-archive
func (s *server) UploadPackHandler() http.HandlerFunc {
	log.Println("UploadPackHandler")
	return s.proxyHandler()
}

func (s *server) proxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, false)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(b))

		upstream := r.Header.Get(UpstreamHeader)
		if upstream == "" {
			upstream = s.defaultUpstreamHost
		}
		upstreamURL, err := url.Parse(upstream)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(upstreamURL)

		// change req.Host so github would not fail
		// and response with a redirect
		defaultDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			defaultDirector(req)
			req.Host = req.URL.Host
		}

		proxy.ServeHTTP(w, r)
	}
}
