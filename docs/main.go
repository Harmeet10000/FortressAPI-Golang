// TLS HTTP/2 server with mutual TLS authentication in Go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

func loadClientCAs() *x509.CertPool {
	clientCAs := x509.NewCertPool()
	caCert, err := os.ReadFile("certs/cert.pem") // Load
	if err != nil {
		log.Fatalf("Failed to read CA cert: %v", err)
	}
	clientCAs.AppendCertsFromPEM(caCert) // Append CA cert to pool
	return clientCAs
}

func main() {
	port := 8000
	cert := "certs/cert.pem"
	key := "certs/key.pem"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Golang Auth API"))
	})

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert, // enforce mTLS
		ClientCAs:  loadClientCAs(),                // load CA cert for client verification
	}
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	http2.ConfigureServer(server, &http2.Server{})
	fmt.Printf("Starting server on port %d...\n", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Failed to start server:", err)
	}
	fmt.Println("Server started successfully")
}
