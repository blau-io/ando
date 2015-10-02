package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"
	"strconv"
)

var (
	client      *http.Client
	globalFlags struct {
		port        int
		remoteDrive string
	}
)

func init() {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	client = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}}

	flag.IntVar(&globalFlags.port, "port", 80, "The port to listen on")
	flag.StringVar(&globalFlags.remoteDrive, "remoteDrive",
		"http://drive.internal.blau.io",
		"The URL of the Google Drive Service")
	flag.Parse()
}

func main() {
	log.Printf("Listening on Port %d", globalFlags.port)
	log.Println("Remote Drive URL: ", globalFlags.remoteDrive)

	http.HandleFunc("/favicon.ico", notfound) //dirty, dirty
	http.HandleFunc("/", start)
	http.ListenAndServe(":"+strconv.Itoa(globalFlags.port), nil)
}

func notfound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func start(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Printf("Could not get token: %v", err)
		http.Error(w, "Could not authorize", http.StatusUnauthorized)
		return
	}

	s := &Site{
		Token: token.Value,
	}
	s.Build()
}
