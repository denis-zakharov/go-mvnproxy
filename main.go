package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	mirror    = os.Getenv("MIRROR")
	cachePath = os.Getenv("CACHE_PATH")
)

func init() {
	// parameters
	if mirror == "" {
		log.Fatal("MIRROR is not set")
	}
	if cachePath == "" {
		log.Fatal("MIRROR is not set")
	}
	if !strings.HasSuffix(mirror, "/") {
		mirror += "/"
	}
}

func checked(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func cached(mirrorURL *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		fn := http.FileServer(http.Dir(cachePath))

		switch {
		case isCached(path):
			fn.ServeHTTP(w, req)
		case strings.HasSuffix(path, ".sha1"):
			computeSHA1(path)
			fn.ServeHTTP(w, req)
		default:
			download(req, mirrorURL)
			fn.ServeHTTP(w, req)
		}
	}
}

func main() {
	// handlers
	mirrorURL, err := url.Parse(mirror)
	checked(err)
	proxyMirror := httputil.NewSingleHostReverseProxy(mirrorURL)
	http.Handle("HEAD /", proxyMirror)
	http.Handle("PUT /", proxyMirror)
	http.HandleFunc("GET /", cached(mirrorURL))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
