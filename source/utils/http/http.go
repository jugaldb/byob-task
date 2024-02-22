package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func Post(url string, headers map[string]string, body any) (*map[string]any, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resbody, err := io.ReadAll(resp.Body)
	m := make(map[string]any)
	err = json.Unmarshal(resbody, &m)
	return &m, err
}

func Get(url string, headers map[string]string) (*map[string]any, error) {

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(make([]byte, 0)))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resbody, err := io.ReadAll(resp.Body)
	m := make(map[string]any)
	err = json.Unmarshal(resbody, &m)
	return &m, err
}

func GetWithBody(url string, headers map[string]string, body any) (*map[string]any, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(payload))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resbody, err := io.ReadAll(resp.Body)
	m := make(map[string]any)
	err = json.Unmarshal(resbody, &m)
	return &m, err
}
