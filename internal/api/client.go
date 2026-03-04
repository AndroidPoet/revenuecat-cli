package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AndroidPoet/revenuecat-cli/internal/config"
)

// BaseURL is the RevenueCat API v2 base URL
const BaseURL = "https://api.revenuecat.com/v2"

// Client wraps the RevenueCat API client
type Client struct {
	httpClient *http.Client
	apiKey     string
	projectID  string
	baseURL    string
	timeout    time.Duration
	debug      bool
}

// APIError represents a structured error from the RevenueCat API
type APIError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	Param      string `json:"param,omitempty"`
	DocURL     string `json:"doc_url,omitempty"`
}

func (e *APIError) Error() string {
	if e.Type != "" {
		return fmt.Sprintf("%s: %s", e.Type, e.Message)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// debugTransport wraps http.RoundTripper to log requests
type debugTransport struct {
	base http.RoundTripper
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf("DEBUG: %s %s\n", req.Method, req.URL)
	for k, v := range req.Header {
		if k == "Authorization" {
			fmt.Printf("DEBUG:   %s: %s...%s\n", k, v[0][:15], v[0][len(v[0])-4:])
		} else {
			fmt.Printf("DEBUG:   %s: %s\n", k, v)
		}
	}
	resp, err := t.base.RoundTrip(req)
	if err == nil {
		fmt.Printf("DEBUG: Response: %d %s\n", resp.StatusCode, resp.Status)
	}
	return resp, err
}

// NewClient creates a new API client
func NewClient(projectID string, timeout time.Duration) (*Client, error) {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, err
	}

	transport := http.DefaultTransport
	if config.IsDebug() {
		transport = &debugTransport{base: transport}
	}

	return &Client{
		httpClient: &http.Client{Transport: transport},
		apiKey:     apiKey,
		projectID:  projectID,
		baseURL:    BaseURL,
		timeout:    timeout,
		debug:      config.IsDebug(),
	}, nil
}

// GetProjectID returns the project ID
func (c *Client) GetProjectID() string {
	return c.projectID
}

// Context returns a context with timeout
func (c *Client) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

// Do executes an API request
func (c *Client) Do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Handle non-2xx responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if err := json.Unmarshal(respBody, apiErr); err != nil {
			apiErr.Message = string(respBody)
		}
		return apiErr
	}

	// Unmarshal response if result is provided
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.Do(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	return c.Do(ctx, http.MethodPost, path, body, result)
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body, result interface{}) error {
	return c.Do(ctx, http.MethodPatch, path, body, result)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) error {
	return c.Do(ctx, http.MethodDelete, path, nil, nil)
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Items      json.RawMessage `json:"items"`
	NextPage   string          `json:"next_page,omitempty"`
	URL        string          `json:"url,omitempty"`
	TotalCount int             `json:"total_count,omitempty"`
}

// ListAll fetches all pages and collects items
func (c *Client) ListAll(ctx context.Context, path string, limit int, collector func(json.RawMessage) error) error {
	cursor := ""
	for {
		pagePath := path
		sep := "?"
		if limit > 0 {
			pagePath += fmt.Sprintf("%slimit=%d", sep, limit)
			sep = "&"
		}
		if cursor != "" {
			pagePath += fmt.Sprintf("%sstarting_after=%s", sep, cursor)
		}

		var resp ListResponse
		if err := c.Get(ctx, pagePath, &resp); err != nil {
			return err
		}

		if err := collector(resp.Items); err != nil {
			return err
		}

		if resp.NextPage == "" {
			break
		}
		cursor = resp.NextPage
	}
	return nil
}
