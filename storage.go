package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func isCached(path string) bool {
	_, err := os.Stat(resolveCachePath(path))
	return err == nil
}

func resolveCachePath(path string) string {
	return filepath.Join(cachePath, path)
}

func computeSHA1(path string) {
	shacp := resolveCachePath(path) // file.sha1
	cp, _ := strings.CutSuffix(shacp, ".sha1")
	data, err := os.ReadFile(cp)
	checked(err)
	bs := sha1.Sum(data)

	f, err := os.Create(shacp)
	checked(err)
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%x", bs))
	checked(err)
}

func download(req *http.Request, mirrorURL *url.URL) {
	downloadPath := mirrorURL.Path + req.URL.Path
	req.URL = mirrorURL
	req.URL.Path = downloadPath
	out, err := os.Create(resolveCachePath(req.URL.Path))
	checked(err)
	defer out.Close()
	resp, err := http.Get(req.URL.String())
	checked(err)
	defer resp.Body.Close()
	for {
		_, err := io.Copy(out, resp.Body)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}
	}
}
