package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func MakeHTTPRequest(method string, url string, body []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	r := bytes.NewReader(body)
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return []byte(""), err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte(""), err
	}

	return data, nil
}
