package chehttp

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

// RequestOption is a function that configures a request
type RequestOption func(*requestConfig)

// requestConfig holds the configuration for a request
type requestConfig struct {
	headers        map[string]string
	timeout        *time.Duration
	body           io.Reader
	successTarget  interface{}
	errorTarget    interface{}
	autoUnmarshal  bool
}

// WithHeaders sets custom headers for the request
func WithHeaders(headers map[string]string) RequestOption {
	return func(rc *requestConfig) {
		if rc.headers == nil {
			rc.headers = make(map[string]string)
		}
		for k, v := range headers {
			rc.headers[k] = v
		}
	}
}

// WithHeader sets a single header for the request
func WithHeader(key, value string) RequestOption {
	return func(rc *requestConfig) {
		if rc.headers == nil {
			rc.headers = make(map[string]string)
		}
		rc.headers[key] = value
	}
}

// WithTimeout sets a timeout for the request
func WithTimeout(timeout time.Duration) RequestOption {
	return func(rc *requestConfig) {
		rc.timeout = &timeout
	}
}

// WithBody sets the request body
func WithBody(body io.Reader) RequestOption {
	return func(rc *requestConfig) {
		rc.body = body
	}
}

// WithJSONBody marshals the provided value to JSON and sets it as the request body
func WithJSONBody(v interface{}) RequestOption {
	return func(rc *requestConfig) {
		data, err := json.Marshal(v)
		if err != nil {
			// Store error in a way that can be checked later
			rc.body = &errorReader{err: err}
			return
		}
		rc.body = bytes.NewReader(data)
		if rc.headers == nil {
			rc.headers = make(map[string]string)
		}
		rc.headers["Content-Type"] = "application/json"
	}
}

// WithSuccess sets the target for unmarshaling successful responses (2xx)
func WithSuccess(target interface{}) RequestOption {
	return func(rc *requestConfig) {
		rc.successTarget = target
		rc.autoUnmarshal = true
	}
}

// WithError sets the target for unmarshaling error responses (4xx, 5xx)
func WithError(target interface{}) RequestOption {
	return func(rc *requestConfig) {
		rc.errorTarget = target
		rc.autoUnmarshal = true
	}
}

// errorReader is used to defer JSON marshal errors
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}
