# chehttp - Ergonomic HTTP Client

A simple, ergonomic HTTP client for Go that makes HTTP requests easier and more convenient.

## Features

- **Builder Pattern**: Fluent API for creating HTTP clients
- **Interface-Based**: Hide implementation details behind clean interfaces
- **Automatic JSON Marshalling**: Automatically marshal request bodies to JSON
- **Automatic JSON Unmarshalling**: Automatically unmarshal response bodies (success and error)
- **Flexible Options**: Configure each request with headers, timeouts, and more
- **Request Lifecycle Hooks**: Pre-request, post-request, success, error, and complete hooks
- **Timeout Control**: Separate connection and request timeouts for fine-grained control
- **Type-Safe**: Fully generic with Go 1.20+ generics
- **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chehttp
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"

    "github.com/comfortablynumb/che/pkg/chehttp"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // Create a client with the builder
    client := chehttp.NewBuilder().
        WithBaseURL("https://api.example.com").
        WithDefaultHeader("Authorization", "Bearer token").
        WithRequestTimeout(30 * time.Second).
        Build()

    // Perform a GET request
    resp, err := client.Get("/users/1")
    if err != nil {
        panic(err)
    }

    fmt.Println("Status:", resp.StatusCode())
    fmt.Println("Body:", resp.String())
}
```

## Usage

### Creating a Client

Use the builder pattern to create HTTP clients:

```go
client := chehttp.NewBuilder().
    WithBaseURL("https://api.example.com").
    WithDefaultHeader("User-Agent", "MyApp/1.0").
    WithDefaultHeader("Accept", "application/json").
    WithRequestTimeout(30 * time.Second).
    WithConnectionTimeout(5 * time.Second).
    WithMaxIdleConns(100).
    WithMaxIdleConnsPerHost(10).
    Build()
```

### Making Requests

The client supports all standard HTTP methods:

```go
// GET request
resp, err := client.Get("/users")

// POST request with JSON body
resp, err := client.Post("/users", chehttp.WithJSONBody(user))

// PUT request
resp, err := client.Put("/users/1", chehttp.WithJSONBody(user))

// PATCH request
resp, err := client.Patch("/users/1", chehttp.WithJSONBody(updates))

// DELETE request
resp, err := client.Delete("/users/1")
```

### Request Options

Configure individual requests with options:

#### Custom Headers

```go
resp, err := client.Get("/users",
    chehttp.WithHeader("X-Request-ID", "12345"),
    chehttp.WithHeader("X-API-Key", "secret"),
)
```

#### Request Timeout

```go
resp, err := client.Get("/users",
    chehttp.WithTimeout(5 * time.Second),
)
```

#### JSON Request Body

```go
user := User{Name: "John"}
resp, err := client.Post("/users",
    chehttp.WithJSONBody(user),
)
```

### Automatic JSON Unmarshalling

Automatically unmarshal success and error responses:

```go
type SuccessResponse struct {
    Data    User   `json:"data"`
    Message string `json:"message"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code"`
}

var success SuccessResponse
var errResp ErrorResponse

resp, err := client.Get("/users/1",
    chehttp.WithSuccess(&success),
    chehttp.WithError(&errResp),
)

if resp.IsSuccess() {
    fmt.Println("User:", success.Data)
} else if resp.IsError() {
    fmt.Println("Error:", errResp.Error)
}
```

### Timeout Control

The client supports two types of timeouts for fine-grained control:

#### Request Timeout

The total time allowed for the entire request (including connection establishment, sending data, and receiving response):

```go
client := chehttp.NewBuilder().
    WithRequestTimeout(30 * time.Second).  // Total request time
    Build()

// Override per request
resp, err := client.Get("/users",
    chehttp.WithTimeout(5 * time.Second),
)
```

#### Connection Timeout

The time allowed to establish the connection:

```go
client := chehttp.NewBuilder().
    WithConnectionTimeout(5 * time.Second).  // Connection establishment time
    WithRequestTimeout(30 * time.Second).     // Total request time
    Build()
```

**Note**: `WithDefaultTimeout()` is an alias for `WithRequestTimeout()` for backward compatibility.

### Request Lifecycle Hooks

Hooks allow you to intercept and observe requests at different stages of their lifecycle. This is useful for logging, metrics, error handling, and more.

#### Available Hooks

- **PreRequest**: Called before sending the request (can cancel by returning error)
- **PostRequest**: Called after receiving the response
- **OnSuccess**: Called when response status is 2xx
- **OnError**: Called when response status is 4xx or 5xx
- **OnComplete**: Called after request completes (success or failure)

#### Hook Context

All hooks receive a `HookContext` with:

```go
type HookContext struct {
    Method     string        // HTTP method (GET, POST, etc.)
    URL        string        // Full URL
    Headers    http.Header   // Request headers
    StatusCode int           // Response status code (0 if not available)
    Response   Response      // Response object (nil if not available)
    Error      error         // Error if request failed (nil if successful)
    StartTime  time.Time     // Request start time
    Duration   time.Duration // Request duration (0 if not complete)
}
```

#### Logging Example

```go
client := chehttp.NewBuilder().
    WithBaseURL("https://api.example.com").
    WithPreRequestHook(func(ctx *chehttp.HookContext) error {
        fmt.Printf("[%s] %s %s\n", ctx.StartTime.Format(time.RFC3339), ctx.Method, ctx.URL)
        return nil
    }).
    WithCompleteHook(func(ctx *chehttp.HookContext) {
        status := ctx.StatusCode
        if ctx.Error != nil {
            status = 0
        }
        fmt.Printf("[%s] %s %s - %d (%v)\n",
            ctx.StartTime.Add(ctx.Duration).Format(time.RFC3339),
            ctx.Method, ctx.URL, status, ctx.Duration)
    }).
    Build()
```

#### Metrics Collection Example

```go
var totalRequests int64
var totalErrors int64
var totalDuration time.Duration

client := chehttp.NewBuilder().
    WithCompleteHook(func(ctx *chehttp.HookContext) {
        atomic.AddInt64(&totalRequests, 1)
        if ctx.Error != nil || ctx.Response.IsError() {
            atomic.AddInt64(&totalErrors, 1)
        }
        // Note: Use mutex for duration in production
        totalDuration += ctx.Duration
    }).
    Build()
```

#### Error Logging Example

```go
client := chehttp.NewBuilder().
    WithErrorHook(func(ctx *chehttp.HookContext) {
        log.Printf("HTTP error: %s %s - status %d, duration %v",
            ctx.Method, ctx.URL, ctx.StatusCode, ctx.Duration)
        if ctx.Response != nil {
            log.Printf("Response body: %s", ctx.Response.String())
        }
    }).
    Build()
```

#### Request Validation Example

```go
client := chehttp.NewBuilder().
    WithPreRequestHook(func(ctx *chehttp.HookContext) error {
        // Validate that API key is present
        if ctx.Headers.Get("X-API-Key") == "" {
            return fmt.Errorf("missing required X-API-Key header")
        }
        return nil
    }).
    Build()

// This request will fail with "pre-request hook failed: missing required X-API-Key header"
resp, err := client.Get("/users")
```

#### Multiple Hooks

You can add multiple hooks of the same type - they will be called in order:

```go
client := chehttp.NewBuilder().
    WithCompleteHook(func(ctx *chehttp.HookContext) {
        // First hook - metrics
        recordMetrics(ctx)
    }).
    WithCompleteHook(func(ctx *chehttp.HookContext) {
        // Second hook - logging
        logRequest(ctx)
    }).
    Build()
```

### Response Methods

The Response interface provides convenient methods:

```go
resp, err := client.Get("/users")

// Get status code
statusCode := resp.StatusCode()

// Get raw body
body := resp.Body()

// Get body as string
str := resp.String()

// Get headers
headers := resp.Headers()
contentType := headers.Get("Content-Type")

// Check if successful (2xx)
if resp.IsSuccess() {
    fmt.Println("Request succeeded!")
}

// Check if error (4xx or 5xx)
if resp.IsError() {
    fmt.Println("Request failed!")
}

// Manual JSON unmarshalling
var data map[string]interface{}
err = resp.UnmarshalJSON(&data)
```

## Examples

### REST API Client

```go
type APIClient struct {
    client chehttp.Client
}

func NewAPIClient(baseURL, token string) *APIClient {
    client := chehttp.NewBuilder().
        WithBaseURL(baseURL).
        WithDefaultHeader("Authorization", "Bearer "+token).
        WithDefaultHeader("Content-Type", "application/json").
        WithRequestTimeout(30 * time.Second).
        Build()

    return &APIClient{client: client}
}

func (c *APIClient) GetUser(id int) (*User, error) {
    var user User
    resp, err := c.client.Get(
        fmt.Sprintf("/users/%d", id),
        chehttp.WithSuccess(&user),
    )
    if err != nil {
        return nil, err
    }

    if !resp.IsSuccess() {
        return nil, fmt.Errorf("failed to get user: %d", resp.StatusCode())
    }

    return &user, nil
}

func (c *APIClient) CreateUser(user *User) error {
    var created User
    var errResp ErrorResponse

    resp, err := c.client.Post("/users",
        chehttp.WithJSONBody(user),
        chehttp.WithSuccess(&created),
        chehttp.WithError(&errResp),
    )
    if err != nil {
        return err
    }

    if resp.IsError() {
        return fmt.Errorf("API error: %s", errResp.Error)
    }

    *user = created
    return nil
}
```

### With Custom HTTP Client

```go
import (
    "crypto/tls"
    "net/http"
)

customHTTPClient := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true, // Don't do this in production!
        },
    },
}

client := chehttp.NewBuilder().
    WithHTTPClient(customHTTPClient).
    WithBaseURL("https://api.example.com").
    Build()
```

### Multiple Headers

```go
headers := map[string]string{
    "X-Request-ID":  "12345",
    "X-Correlation": "abc-def",
    "X-API-Version": "v1",
}

resp, err := client.Get("/users",
    chehttp.WithHeaders(headers),
)
```

## API Reference

### Builder

**Basic Configuration:**
- `NewBuilder()` - Creates a new builder
- `WithHTTPClient(client)` - Sets custom http.Client
- `WithBaseURL(url)` - Sets base URL for all requests
- `WithDefaultHeader(key, value)` - Adds a default header
- `WithDefaultHeaders(headers)` - Adds multiple default headers
- `Build()` - Builds the client

**Timeout Configuration:**
- `WithRequestTimeout(duration)` - Sets total request timeout
- `WithConnectionTimeout(duration)` - Sets connection establishment timeout
- `WithDefaultTimeout(duration)` - Alias for WithRequestTimeout (backward compatibility)

**Transport Configuration:**
- `WithTransport(transport)` - Sets custom transport
- `WithMaxIdleConns(n)` - Sets max idle connections
- `WithMaxIdleConnsPerHost(n)` - Sets max idle connections per host
- `WithInsecureSkipVerify()` - Disables TLS certificate verification (use with caution!)

**Lifecycle Hooks:**
- `WithPreRequestHook(hook)` - Adds a pre-request hook
- `WithPostRequestHook(hook)` - Adds a post-request hook
- `WithSuccessHook(hook)` - Adds a success hook (2xx responses)
- `WithErrorHook(hook)` - Adds an error hook (4xx/5xx responses)
- `WithCompleteHook(hook)` - Adds a complete hook (always called)

### Client Interface

- `Get(url, ...opts)` - Performs GET request
- `Post(url, ...opts)` - Performs POST request
- `Put(url, ...opts)` - Performs PUT request
- `Patch(url, ...opts)` - Performs PATCH request
- `Delete(url, ...opts)` - Performs DELETE request
- `Do(method, url, ...opts)` - Performs request with custom method

### Request Options

- `WithHeader(key, value)` - Adds a header
- `WithHeaders(headers)` - Adds multiple headers
- `WithTimeout(duration)` - Sets request timeout
- `WithBody(reader)` - Sets raw request body
- `WithJSONBody(v)` - Marshals v to JSON and sets as body
- `WithSuccess(target)` - Auto-unmarshal success response to target
- `WithError(target)` - Auto-unmarshal error response to target

### Response Interface

- `StatusCode()` - Returns HTTP status code
- `Body()` - Returns raw body bytes
- `String()` - Returns body as string
- `Headers()` - Returns response headers
- `UnmarshalJSON(v)` - Unmarshals body to v
- `IsSuccess()` - Returns true if 2xx status
- `IsError()` - Returns true if 4xx or 5xx status

## Performance Considerations

- Connection pooling is enabled by default
- Use `WithMaxIdleConns` and `WithMaxIdleConnsPerHost` for high-throughput scenarios
- The base URL is concatenated on each request - consider caching if performance critical
- Response bodies are fully read into memory - be cautious with large responses

## Best Practices

1. **Reuse Clients**: Create one client per API and reuse it
2. **Set Timeouts**: Always set appropriate timeouts
3. **Handle Errors**: Always check both err and response status
4. **Use Auto-Unmarshal**: Use `WithSuccess` and `WithError` for cleaner code
5. **Base URL**: Use base URL for API clients
6. **Default Headers**: Set common headers like User-Agent, Authorization as defaults

## License

This package is part of the Che library and shares the same license.
