package chehttp

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Builder is used to build HTTP clients with custom configuration
type Builder struct {
	httpClient        *http.Client
	baseURL           string
	defaultHeaders    map[string]string
	requestTimeout    time.Duration
	connectionTimeout time.Duration
	hooks             *Hooks
	retryConfig       *RetryConfig
}

// NewBuilder creates a new HTTP client builder
func NewBuilder() *Builder {
	return &Builder{
		httpClient:     &http.Client{},
		defaultHeaders: make(map[string]string),
		hooks:          &Hooks{},
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

// WithRequestTimeout sets the default request timeout (total time for the entire request)
func (b *Builder) WithRequestTimeout(timeout time.Duration) *Builder {
	b.requestTimeout = timeout
	return b
}

// WithDefaultTimeout is an alias for WithRequestTimeout for backward compatibility
func (b *Builder) WithDefaultTimeout(timeout time.Duration) *Builder {
	return b.WithRequestTimeout(timeout)
}

// WithConnectionTimeout sets the connection timeout (time to establish connection)
func (b *Builder) WithConnectionTimeout(timeout time.Duration) *Builder {
	b.connectionTimeout = timeout
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

// WithPreRequestHook adds a pre-request hook
func (b *Builder) WithPreRequestHook(hook PreRequestHook) *Builder {
	if b.hooks == nil {
		b.hooks = &Hooks{}
	}
	b.hooks.PreRequest = append(b.hooks.PreRequest, hook)
	return b
}

// WithPostRequestHook adds a post-request hook
func (b *Builder) WithPostRequestHook(hook PostRequestHook) *Builder {
	if b.hooks == nil {
		b.hooks = &Hooks{}
	}
	b.hooks.PostRequest = append(b.hooks.PostRequest, hook)
	return b
}

// WithSuccessHook adds a success hook (called on 2xx responses)
func (b *Builder) WithSuccessHook(hook SuccessHook) *Builder {
	if b.hooks == nil {
		b.hooks = &Hooks{}
	}
	b.hooks.OnSuccess = append(b.hooks.OnSuccess, hook)
	return b
}

// WithErrorHook adds an error hook (called on 4xx/5xx responses)
func (b *Builder) WithErrorHook(hook ErrorHook) *Builder {
	if b.hooks == nil {
		b.hooks = &Hooks{}
	}
	b.hooks.OnError = append(b.hooks.OnError, hook)
	return b
}

// WithCompleteHook adds a complete hook (always called after request)
func (b *Builder) WithCompleteHook(hook CompleteHook) *Builder {
	if b.hooks == nil {
		b.hooks = &Hooks{}
	}
	b.hooks.OnComplete = append(b.hooks.OnComplete, hook)
	return b
}

// WithRetryConfig sets the retry configuration
func (b *Builder) WithRetryConfig(config *RetryConfig) *Builder {
	b.retryConfig = config
	return b
}

// WithRetries enables retries with the specified maximum number of retry attempts
func (b *Builder) WithRetries(maxRetries int) *Builder {
	if b.retryConfig == nil {
		b.retryConfig = DefaultRetryConfig()
	}
	b.retryConfig.MaxRetries = maxRetries
	return b
}

// WithRetryBackoff sets the backoff strategy for retries
func (b *Builder) WithRetryBackoff(strategy BackoffStrategy) *Builder {
	if b.retryConfig == nil {
		b.retryConfig = DefaultRetryConfig()
	}
	b.retryConfig.BackoffStrategy = strategy
	return b
}

// WithRetryableStatusCodes sets the HTTP status codes that should trigger retries
func (b *Builder) WithRetryableStatusCodes(ranges ...StatusCodeRange) *Builder {
	if b.retryConfig == nil {
		b.retryConfig = DefaultRetryConfig()
	}
	b.retryConfig.RetryableStatusCodes = ranges
	return b
}

// WithNonRetryableStatusCodes sets the HTTP status codes that should never be retried
func (b *Builder) WithNonRetryableStatusCodes(codes ...int) *Builder {
	if b.retryConfig == nil {
		b.retryConfig = DefaultRetryConfig()
	}
	b.retryConfig.NonRetryableStatusCodes = codes
	return b
}

// Build creates the HTTP client
func (b *Builder) Build() Client {
	// Configure connection timeout via transport
	if b.connectionTimeout > 0 {
		if b.httpClient.Transport == nil {
			b.httpClient.Transport = &http.Transport{}
		}
		if transport, ok := b.httpClient.Transport.(*http.Transport); ok {
			transport.DialContext = (&http.Transport{
				DialContext: (&http.Transport{}).DialContext,
			}).DialContext
			// Note: DialContext timeout is set in the transport, but we need to use net.Dialer
			// For simplicity, we'll set it in the client Do method
		}
	}

	return &client{
		httpClient:        b.httpClient,
		baseURL:           b.baseURL,
		defaultHeaders:    b.defaultHeaders,
		requestTimeout:    b.requestTimeout,
		connectionTimeout: b.connectionTimeout,
		hooks:             b.hooks.Clone(),
		retryConfig:       b.retryConfig,
	}
}
