package chehttp

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Builder is used to build HTTP clients with custom configuration
type Builder struct {
	httpClient     *http.Client
	baseURL        string
	defaultHeaders map[string]string
	defaultTimeout time.Duration
}

// NewBuilder creates a new HTTP client builder
func NewBuilder() *Builder {
	return &Builder{
		httpClient:     &http.Client{},
		defaultHeaders: make(map[string]string),
	}
}

// WithHTTPClient sets a custom http.Client
func (b *Builder) WithHTTPClient(httpClient *http.Client) *Builder {
	b.httpClient = httpClient
	return b
}

// WithBaseURL sets the base URL for all requests
func (b *Builder) WithBaseURL(baseURL string) *Builder {
	b.baseURL = baseURL
	return b
}

// WithDefaultHeader sets a default header that will be included in all requests
func (b *Builder) WithDefaultHeader(key, value string) *Builder {
	if b.defaultHeaders == nil {
		b.defaultHeaders = make(map[string]string)
	}
	b.defaultHeaders[key] = value
	return b
}

// WithDefaultHeaders sets multiple default headers
func (b *Builder) WithDefaultHeaders(headers map[string]string) *Builder {
	if b.defaultHeaders == nil {
		b.defaultHeaders = make(map[string]string)
	}
	for k, v := range headers {
		b.defaultHeaders[k] = v
	}
	return b
}

// WithDefaultTimeout sets the default timeout for all requests
func (b *Builder) WithDefaultTimeout(timeout time.Duration) *Builder {
	b.defaultTimeout = timeout
	return b
}

// WithTransport sets a custom transport for the HTTP client
func (b *Builder) WithTransport(transport http.RoundTripper) *Builder {
	if b.httpClient == nil {
		b.httpClient = &http.Client{}
	}
	b.httpClient.Transport = transport
	return b
}

// WithMaxIdleConns sets the maximum number of idle connections
func (b *Builder) WithMaxIdleConns(n int) *Builder {
	if b.httpClient == nil {
		b.httpClient = &http.Client{}
	}
	if b.httpClient.Transport == nil {
		b.httpClient.Transport = &http.Transport{}
	}
	if transport, ok := b.httpClient.Transport.(*http.Transport); ok {
		transport.MaxIdleConns = n
	}
	return b
}

// WithMaxIdleConnsPerHost sets the maximum number of idle connections per host
func (b *Builder) WithMaxIdleConnsPerHost(n int) *Builder {
	if b.httpClient == nil {
		b.httpClient = &http.Client{}
	}
	if b.httpClient.Transport == nil {
		b.httpClient.Transport = &http.Transport{}
	}
	if transport, ok := b.httpClient.Transport.(*http.Transport); ok {
		transport.MaxIdleConnsPerHost = n
	}
	return b
}

// WithInsecureSkipVerify disables TLS certificate verification (use with caution!)
func (b *Builder) WithInsecureSkipVerify() *Builder {
	if b.httpClient == nil {
		b.httpClient = &http.Client{}
	}
	if b.httpClient.Transport == nil {
		b.httpClient.Transport = &http.Transport{}
	}
	if transport, ok := b.httpClient.Transport.(*http.Transport); ok {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = true
	}
	return b
}

// Build creates the HTTP client
func (b *Builder) Build() Client {
	return &client{
		httpClient:     b.httpClient,
		baseURL:        b.baseURL,
		defaultHeaders: b.defaultHeaders,
		defaultTimeout: b.defaultTimeout,
	}
}
