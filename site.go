package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/a8m/mark"
)

type Site struct {
	BaseURL string
	Token   string
}

func (s *Site) Build() error {
	list, err := s.GetFileList()
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

	log.Println("Getting file content")

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

	name = strings.Replace(name, ".md", ".html", -1)
	remote = strings.TrimSuffix(globalFlags.remoteDrive, "/") +
		"/add/blau.io/PUBLIC/" + name
	reqBody := strings.NewReader(mark.Render(string(body)))
	req, err = http.NewRequest("PUT", remote, reqBody)
	if err != nil {
		log.Printf("Error while creating request: %v", err)
		return
	}

	log.Println("Uploading file content")

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

func (s *Site) GetFileList() ([]string, error) {
	remote := strings.TrimSuffix(globalFlags.remoteDrive, "/") +
		"/browse/blau.io/content"
	req, err := http.NewRequest("GET", remote, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: s.Token,
	})

	log.Println("Getting file list")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var list []string
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&list)
	return list, err
}
