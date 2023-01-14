package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// setup a simple handler which sends a HTHS header for six months (!)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
		fmt.Fprintf(w, "Hello, HTTPS world!")
	})

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// Replace with your domain
		HostPolicy: func(ctx context.Context, host string) error {
			// * Always allow domain.
			return nil
		},
		// Folder to store the certificates
		Cache: autocert.DirCache("./certs"),
	}

	// create the server itself
	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	log.Printf("Serving http/https for domains: *")
	go func() {
		// serve HTTP, which will redirect automatically to HTTPS
		h := certManager.HTTPHandler(nil)
		log.Fatal(http.ListenAndServe(":http", h))
	}()

	// serve HTTPS!
	log.Fatal(server.ListenAndServeTLS("", ""))
}

// cacheDir makes a consistent cache directory inside /tmp. Returns "" on error.
func cacheDir() (dir string) {
	if u, _ := user.Current(); u != nil {
		dir = filepath.Join(os.TempDir(), "cache-golang-autocert-"+u.Username)
		if err := os.MkdirAll(dir, 0700); err == nil {
			return dir
		}
	}
	return ""
}
