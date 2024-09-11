package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient interface uses the native http.Client Do method
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RateLimiter struct to handle rate limiting
type RateLimiter struct {
	ticker *time.Ticker
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		ticker: time.NewTicker(interval),
	}
}

// Wait ensures the rate limit is respected
func (rl *RateLimiter) Wait() {
	<-rl.ticker.C
}

// Client struct that wraps the native http.Client and allows middleware to be used
type Client struct {
	httpClient  HTTPClient
	rateLimiter *RateLimiter
}

// NewClient initializes a new Client with a given http.Client or the default one
func NewClient(httpClient HTTPClient, rateLimitInterval time.Duration) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		httpClient:  httpClient,
		rateLimiter: NewRateLimiter(rateLimitInterval),
	}
}

// CreateRequest abstracts the logic for creating an HTTP request with context support.
func (c *Client) CreateRequest(ctx context.Context, methodType, url string, body []byte) (*http.Request, error) {
	var req *http.Request
	var err error

	// Create the request based on method type with context
	if methodType == http.MethodPost && body != nil {
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	}

	if err != nil {
		return nil, err
	}

	return req, nil
}

// HandleResponse abstracts the logic for handling the HTTP response.
func (c *Client) HandleResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	// Check if the status code is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch external data: %s", resp.Status)
	}

	// Read the response body
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ApiCall performs the HTTP request using the native http.Client Do method and returns the response, with context support.
func (c *Client) ApiCall(ctx context.Context, methodType, url string, body []byte) ([]byte, error) {
	// Wait for the rate limiter before making the request
	c.rateLimiter.Wait()

	// Create the request using the abstracted function
	req, err := c.CreateRequest(ctx, methodType, url, body)
	if err != nil {
		return nil, err
	}

	// Perform the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Handle the response using the abstracted function
	return c.HandleResponse(resp)
}
