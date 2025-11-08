package chehttp

import (
	"encoding/json"
	"io"
	"net/http"
)

// Response represents an HTTP response with convenient methods
type Response interface {
	// StatusCode returns the HTTP status code
	StatusCode() int

	// Body returns the raw response body
	Body() []byte

	// Headers returns the response headers
	Headers() http.Header

	// UnmarshalJSON unmarshals the response body into the provided value
	UnmarshalJSON(v interface{}) error

	// String returns the response body as a string
	String() string

	// IsSuccess returns true if the status code is 2xx
	IsSuccess() bool

	// IsError returns true if the status code is 4xx or 5xx
	IsError() bool
}

// response implements the Response interface
type response struct {
	statusCode int
	body       []byte
	headers    http.Header
}

// StatusCode returns the HTTP status code
func (r *response) StatusCode() int {
	return r.statusCode
}

// Body returns the raw response body
func (r *response) Body() []byte {
	return r.body
}

// Headers returns the response headers
func (r *response) Headers() http.Header {
	return r.headers
}

// UnmarshalJSON unmarshals the response body into the provided value
func (r *response) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(r.body, v)
}

// String returns the response body as a string
func (r *response) String() string {
	return string(r.body)
}

// IsSuccess returns true if the status code is 2xx
func (r *response) IsSuccess() bool {
	return r.statusCode >= 200 && r.statusCode < 300
}

// IsError returns true if the status code is 4xx or 5xx
func (r *response) IsError() bool {
	return r.statusCode >= 400
}

// newResponse creates a new Response from an http.Response
func newResponse(httpResp *http.Response) (Response, error) {
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return &response{
		statusCode: httpResp.StatusCode,
		body:       body,
		headers:    httpResp.Header,
	}, nil
}
