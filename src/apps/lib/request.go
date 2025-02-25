package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HTTPRequestMethodType string

const (
	HTTPRequestMethodGet    HTTPRequestMethodType = "GET"
	HTTPRequestMethodPost   HTTPRequestMethodType = "POST"
	HTTPRequestMethodPut    HTTPRequestMethodType = "PUT"
	HTTPRequestMethodDelete HTTPRequestMethodType = "DELETE"
)

type HTTPRequestOptions struct {
	Method   HTTPRequestMethodType
	Endpoint string
	Headers  map[string]string
	Query    map[string]string
	Body     map[string]any
}

// HTTPRequest sends an HTTP request with given method, URL, headers, query parameters, and body as map[string]any.
func HTTPRequest(options HTTPRequestOptions) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Build the full URL with query parameters
	reqURL, err := url.Parse(options.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	query := reqURL.Query()
	for key, value := range options.Query {
		query.Set(key, value)
	}
	reqURL.RawQuery = query.Encode()

	// Convert body map to JSON if provided
	var reqBody io.Reader
	if options.Body != nil {
		jsonBody, err := json.Marshal(options.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create request with body
	req, err := http.NewRequest(string(options.Method), reqURL.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// Ensure Content-Type is set when sending a JSON body
	if options.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
