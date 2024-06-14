package engines

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type HttpRequestEngine struct {
	host string
	port int
}

func NewHttpRequestEngine(host string, port int) *HttpRequestEngine {
	return &HttpRequestEngine{
		host: host,
		port: port,
	}
}

func (h HttpRequestEngine) DoRequest(method string, requestBody string) []byte {
	url, err := url.Parse(h.host)
	if err != nil {
		log.Fatalf("Failed to parsing url: %v", err)
	}

	if h.port > 0 {
		url, err = url.Parse(fmt.Sprintf("%s:%d", h.host, h.port))
		if err != nil {
			log.Fatalf("Failed to parsing url: %v", err)
		}
	}

	var body io.Reader = nil

	if method != http.MethodGet && len(requestBody) > 0 {
		body = bytes.NewReader([]byte(requestBody))
	}

	if method == http.MethodGet && len(requestBody) > 0 {
		url.RawQuery = requestBody
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	defer client.CloseIdleConnections()
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	return response
}
