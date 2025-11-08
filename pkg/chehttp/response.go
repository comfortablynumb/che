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

	// Body returns the raw response body (reads and caches if not already read)
	Body() []byte

	// BodyReader returns the underlying response body reader
	// Note: Once you use BodyReader, you cannot use Body() or UnmarshalJSON()
	// as the reader will be consumed
	BodyReader() io.ReadCloser

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
	bodyReader io.ReadCloser
	headers    http.Header
	bodyRead   bool
}

// StatusCode returns the HTTP status code
func (r *response) StatusCode() int {
	return r.statusCode
}

// Body returns the raw response body (reads and caches if not already read)
func (r *response) Body() []byte {
	if !r.bodyRead && r.bodyReader != nil {
		r.readBody()
	}
	return r.body
}

// BodyReader returns the underlying response body reader
func (r *response) BodyReader() io.ReadCloser {
	return r.bodyReader
}

// readBody reads the body from the reader and caches it
func (r *response) readBody() {
	if r.bodyReader == nil {
		return
	}

	body, err := io.ReadAll(r.bodyReader)
	if err == nil {
		r.body = body
	}
	r.bodyRead = true
	r.bodyReader.Close()
}

// Headers returns the response headers
func (r *response) Headers() http.Header {
	return r.headers
}

// UnmarshalJSON unmarshals the response body into the provided value
func (r *response) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(r.Body(), v)
}

// String returns the response body as a string
func (r *response) String() string {
	return string(r.Body())
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
// If readBody is true, the body is read immediately and cached
// Otherwise, the body reader is kept for later use
func newResponse(httpResp *http.Response, readBody bool) (Response, error) {
	resp := &response{
		statusCode: httpResp.StatusCode,
		bodyReader: httpResp.Body,
		headers:    httpResp.Header,
		bodyRead:   false,
	}

	if readBody {
		resp.readBody()
	}

	return resp, nil
}
