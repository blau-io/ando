package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func getFileList(path, token string) ([]string, error) {
	remote := strings.TrimSuffix(globalFlags.remoteDrive, "/") +
		"/browse" + path
	req, err := http.NewRequest("GET", remote, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var list []string
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&list)
	return list, err
}

func mkdir(path, token string) error {
	remote := strings.TrimSuffix(globalFlags.remoteDrive, "/") +
		"/add" + path
	req, err := http.NewRequest("POST", remote, nil)
	if err != nil {
		return err
	}

	req.Header.Add("folder", "true")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func recursiveDelete(root, token string) error {
	list, err := getFileList(root, token)
	if err != nil {
		return err
	}

	for _, v := range list {
		err = recursiveDelete(root+"/"+v, token)
		if err != nil {
			return err
		}
	}

	remote := strings.TrimSuffix(globalFlags.remoteDrive, "/") +
		"/delete" + root
	req, err := http.NewRequest("DELETE", remote, nil)
	if err != nil {
		return err
	}

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func upload(filepath, remotepath, token string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	remote := strings.TrimSuffix(globalFlags.remoteDrive, "/") +
		"/add" + remotepath
	req, err := http.NewRequest("POST", remote, file)
	if err != nil {
		return err
	}

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
