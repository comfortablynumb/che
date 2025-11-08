package chehttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an HTTP client with convenient methods
type Client interface {
	// Get performs an HTTP GET request
	Get(url string, opts ...RequestOption) (Response, error)

	// Post performs an HTTP POST request
	Post(url string, opts ...RequestOption) (Response, error)

	// Put performs an HTTP PUT request
	Put(url string, opts ...RequestOption) (Response, error)

	// Patch performs an HTTP PATCH request
	Patch(url string, opts ...RequestOption) (Response, error)

	// Delete performs an HTTP DELETE request
	Delete(url string, opts ...RequestOption) (Response, error)

	// Do performs an HTTP request with the given method
	Do(method, url string, opts ...RequestOption) (Response, error)
}

// client implements the Client interface
type client struct {
	httpClient        *http.Client
	baseURL           string
	defaultHeaders    map[string]string
	requestTimeout    time.Duration
	connectionTimeout time.Duration
	hooks             *Hooks
}

// Get performs an HTTP GET request
func (c *client) Get(url string, opts ...RequestOption) (Response, error) {
	return c.Do(http.MethodGet, url, opts...)
}

// Post performs an HTTP POST request
func (c *client) Post(url string, opts ...RequestOption) (Response, error) {
	return c.Do(http.MethodPost, url, opts...)
}

// Put performs an HTTP PUT request
func (c *client) Put(url string, opts ...RequestOption) (Response, error) {
	return c.Do(http.MethodPut, url, opts...)
}

// Patch performs an HTTP PATCH request
func (c *client) Patch(url string, opts ...RequestOption) (Response, error) {
	return c.Do(http.MethodPatch, url, opts...)
}

// Delete performs an HTTP DELETE request
func (c *client) Delete(url string, opts ...RequestOption) (Response, error) {
	return c.Do(http.MethodDelete, url, opts...)
}

// Do performs an HTTP request with the given method
func (c *client) Do(method, url string, opts ...RequestOption) (Response, error) {
	startTime := time.Now()

	// Apply request configuration
	config := &requestConfig{
		headers: make(map[string]string),
	}

	// Apply default headers
	for k, v := range c.defaultHeaders {
		config.headers[k] = v
	}

	// Apply request options
	for _, opt := range opts {
		opt(config)
	}

	// Build full URL
	fullURL := c.buildURL(url)

	// Create request
	var body io.Reader
	if config.body != nil {
		body = config.body
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range config.headers {
		req.Header.Set(k, v)
	}

	// Create hook context
	hookCtx := &HookContext{
		Method:    method,
		URL:       fullURL,
		Headers:   req.Header,
		StartTime: startTime,
	}

	// Call pre-request hooks
	if c.hooks != nil {
		for _, hook := range c.hooks.PreRequest {
			if err := hook(hookCtx); err != nil {
				hookCtx.Error = err
				hookCtx.Duration = time.Since(startTime)
				// Call complete hooks even on pre-request error
				for _, completeHook := range c.hooks.OnComplete {
					completeHook(hookCtx)
				}
				return nil, fmt.Errorf("pre-request hook failed: %w", err)
			}
		}
	}

	// Set request timeout (total time for the request)
	timeout := c.requestTimeout
	if config.timeout != nil {
		timeout = *config.timeout
	}

	ctx := req.Context()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	req = req.WithContext(ctx)

	// Perform request
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		hookCtx.Error = err
		hookCtx.Duration = time.Since(startTime)

		// Call complete hooks on error
		if c.hooks != nil {
			for _, completeHook := range c.hooks.OnComplete {
				completeHook(hookCtx)
			}
		}

		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Create response
	resp, err := newResponse(httpResp)
	if err != nil {
		hookCtx.Error = err
		hookCtx.Duration = time.Since(startTime)
		hookCtx.StatusCode = httpResp.StatusCode

		// Call complete hooks on error
		if c.hooks != nil {
			for _, completeHook := range c.hooks.OnComplete {
				completeHook(hookCtx)
			}
		}

		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Update hook context with response info
	hookCtx.StatusCode = resp.StatusCode()
	hookCtx.Response = resp
	hookCtx.Duration = time.Since(startTime)

	// Call post-request hooks
	if c.hooks != nil {
		for _, hook := range c.hooks.PostRequest {
			hook(hookCtx)
		}

		// Call success/error hooks
		if resp.IsSuccess() {
			for _, hook := range c.hooks.OnSuccess {
				hook(hookCtx)
			}
		} else if resp.IsError() {
			for _, hook := range c.hooks.OnError {
				hook(hookCtx)
			}
		}

		// Call complete hooks
		for _, completeHook := range c.hooks.OnComplete {
			completeHook(hookCtx)
		}
	}

	// Auto-unmarshal if configured
	if config.autoUnmarshal {
		if resp.IsSuccess() && config.successTarget != nil {
			if err := resp.UnmarshalJSON(config.successTarget); err != nil {
				return resp, fmt.Errorf("failed to unmarshal success response: %w", err)
			}
		} else if resp.IsError() && config.errorTarget != nil {
			if err := resp.UnmarshalJSON(config.errorTarget); err != nil {
				return resp, fmt.Errorf("failed to unmarshal error response: %w", err)
			}
		}
	}

	return resp, nil
}

// buildURL builds the full URL from base URL and path
func (c *client) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}

	// Simple URL joining - in production you might want to use url.Parse
	if len(path) > 0 && path[0] == '/' {
		return c.baseURL + path
	}
	return c.baseURL + "/" + path
}
