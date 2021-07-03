package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	UpstreamHeader = "GIT-CACHE-UPSTREAM"
)

func (s *server) InfoRefHandler() http.HandlerFunc {
	log.Println("InfoRefHandler")
	return s.proxyHandler()
}

func (s *server) ReceivePackHandler() http.HandlerFunc {
	log.Println("ReceivePackHandler")
	return s.proxyHandler()
}

func (s *server) UploadPackHandler() http.HandlerFunc {
	log.Println("UploadPackHandler")
	return s.proxyHandler()
}

func (s *server) proxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, true)
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
