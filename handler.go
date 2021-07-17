package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/google/gitprotocolio"
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
	return func(w http.ResponseWriter, r *http.Request) {
		proxyHandler := s.proxyHandler()

		// Middleware layer analyzes the request body to determine the command
		// and body chunk to determine:
		//
		//   - Should the response be cached?
		//   - What would be the caching key?
		//   - Does a valid cache entry exists in the store?
		//
		// For now this is accomplished by consuming the entire request body
		// into memory.  In the future, we should be smart about when to stop
		// and/or be more selective on which requests to analyze the body.

		bodyReader, err := extractRequestBody(r)
		if err != nil {
			log.Fatalln("could not extract request body: %w\n", err)
		}

		// Parse and handle git protocol command and content
		scanner := gitprotocolio.NewProtocolV2Request(bodyReader)
		for {
			if !scanner.Scan() {
				if scanner.Err() != nil {
					log.Printf("Unable to scan request: %s\n", scanner.Err())
				}

				break
			}

			c := ConvertChunk(scanner.Chunk())

			data, err := json.MarshalIndent(c, "", " ")
			if err != nil {
				log.Fatalln("Unable to unmarshal chunk: %w", err)
			}
			log.Printf("chunk: %s", string(data))
		}

		// Pass the original request to proxy to upstream
		proxyHandler(w, r)
	}
}

type ProtocolV2RequestConvertedChunk struct {
	*gitprotocolio.ProtocolV2RequestChunk
	ArgumentString string
}

func ConvertChunk(c *gitprotocolio.ProtocolV2RequestChunk) *ProtocolV2RequestConvertedChunk {
	return &ProtocolV2RequestConvertedChunk{
		ProtocolV2RequestChunk: c,
		ArgumentString:         string(c.Argument),
	}
}

// extractRequestBody extracts the request body into a new io.Reader without
// affecting the request state.  The returned reader is decompressed for easier
// processing.
func extractRequestBody(r *http.Request) (io.Reader, error) {
	// Read from request Body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("Unable to read body: %w", err)
	}
	// Rewrite the read bytes into the body to upstream request
	// this is needed as the request body buffer was closed after
	// previous read.
	r.Body = io.NopCloser(bytes.NewReader(b))

	// Handle compression
	var bodyReader io.Reader
	if r.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, fmt.Errorf("cannot create gzip reader: %w", err)
		}

		bodyReader = gzipReader
	} else {
		bodyReader = bytes.NewReader(b)
	}

	return bodyReader, nil
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
