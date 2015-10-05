package main

import (
	"fmt"
	"log"
	"net/http"
)

func notfound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func render(w http.ResponseWriter, r *http.Request) {
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

	address, err := publish(token.Value)
	if err != nil {
		log.Printf("Failed to publish public folder: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, address)
}

func setup(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Printf("Could not get token: %v", err)
		http.Error(w, "Could not authorize", http.StatusUnauthorized)
		return
	}

	_, err = getFileList("/blau.io", token.Value)
	if err == nil {
		err = recursiveDelete("/blau.io", token.Value)
		if err != nil {
			log.Printf("Could not delete folders: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	err = mkdir("/blau.io", token.Value)
	if err != nil {
		log.Printf("Could not create folder blau.io: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = mkdir("/blau.io/configuration", token.Value)
	if err != nil {
		log.Printf("Could not create folder configuration: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = mkdir("/blau.io/content", token.Value)
	if err != nil {
		log.Printf("Could not create folder content: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = mkdir("/blau.io/PUBLIC", token.Value)
	if err != nil {
		log.Printf("Could not create folder PUBLIC: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = upload("html/index.html", "/blau.io/configuration/index.html",
		token.Value)
	if err != nil {
		log.Printf("Could not upload index.html: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = upload("html/main.css", "/blau.io/configuration/main.css",
		token.Value)
	if err != nil {
		log.Printf("Could not upload index.html: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = upload("example/hello.md", "/blau.io/content/hello.md", token.Value)
	if err != nil {
		log.Printf("Could not upload hello.md: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
