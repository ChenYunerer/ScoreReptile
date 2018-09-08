package net

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func GetRequest(url string) (string, error) {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", errors.New("response code is " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

var client = &http.Client{}

func GetRequestForReader(url string) (io.Reader, error) {
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("response code is " + response.Status)
	}
	return response.Body, nil
}
