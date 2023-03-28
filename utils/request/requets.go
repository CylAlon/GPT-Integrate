package request

import (
	"bytes"
	"io"
	"net/http"
	"time"
)
func httpRequest(url string, headers map[string]string, body []byte,timeout int,port string) ([]byte, error) {
	req, err := http.NewRequest(port, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	tr := &http.Transport{
		MaxIdleConns: 100,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rebody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return rebody, nil
}

func HttpGet(url string, headers map[string]string,timeout int) ([]byte, error) {
	return httpRequest(url, headers, nil,timeout,"GET")
}

func HttpPost(url string, headers map[string]string, body []byte,timeout int) ([]byte, error) {
	return httpRequest(url, headers, body,timeout,"POST")
}
