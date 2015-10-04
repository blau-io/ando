package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/a8m/mark"
)

type Site struct {
	Token string
}

func (s *Site) Build() error {
	list, err := getFileList("/blau.io/content", s.Token)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(list))

	for _, name := range list {
		go s.CreatePage(name, wg)
	}

	wg.Wait()
	return nil
}

func (s *Site) CreatePage(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	remote := strings.TrimSuffix(globalFlags.remoteDrive, "/") +
		"/read/blau.io/content/" + name
	req, err := http.NewRequest("GET", remote, nil)
	if err != nil {
		log.Printf("Error while creating request: %v", err)
		return
	}

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: s.Token,
	})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error while making request: %v", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading response: %v", err)
		return
	}

	header, err := ioutil.ReadFile("html/header.html")
	if err != nil {
		log.Printf("Error loading header.html: %v", err)
		return
	}

	footer, err := ioutil.ReadFile("html/footer.html")
	if err != nil {
		log.Printf("Error loading footer.html: %v", err)
		return
	}

	reqBody := string(header) + mark.Render(string(body)) + string(footer)

	name = strings.Replace(name, ".md", ".html", -1)
	remote = strings.TrimSuffix(globalFlags.remoteDrive, "/") +
		"/add/blau.io/PUBLIC/" + name
	req, err = http.NewRequest("POST", remote, strings.NewReader(reqBody))
	if err != nil {
		log.Printf("Error while creating request: %v", err)
		return
	}

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: s.Token,
	})

	_, err = client.Do(req)
	if err != nil {
		log.Printf("Error while making request: %v", err)
		return
	}
}
