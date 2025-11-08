# cheenv - Environment Variable Utilities

Type-safe environment variable utilities for Go with defaults, type conversions, and list support.

## Features

- **Type-Safe Conversions**: Automatic parsing to int, int64, float64, bool, duration
- **Default Values**: Graceful fallback when variables are missing or invalid
- **List Support**: Parse comma-separated or custom-delimited lists
- **Flexible Boolean Parsing**: Accepts true, false, yes, no, on, off, 1, 0, y, n, t, f
- **Prefix Filtering**: Get all variables matching a prefix
- **Must Variants**: Panic on missing/invalid values for required config
- **Zero Dependencies**: Only uses Go standard library

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/cheenv
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"

    "github.com/comfortablynumb/che/pkg/cheenv"
)

func main() {
    // String values
    dbHost := cheenv.Get("DB_HOST", "localhost")
    dbName := cheenv.MustGet("DB_NAME") // panics if not set

    // Numeric values
    dbPort := cheenv.GetInt("DB_PORT", 5432)
    maxConns := cheenv.GetInt64("MAX_CONNECTIONS", 100)
    timeout := cheenv.GetFloat("TIMEOUT", 30.0)

    // Boolean values (accepts: true, false, yes, no, on, off, 1, 0, y, n, t, f)
    debug := cheenv.GetBool("DEBUG", false)
    enableCache := cheenv.MustGetBool("ENABLE_CACHE")

    // Duration values (e.g., "5s", "10m", "2h")
    requestTimeout := cheenv.GetDuration("REQUEST_TIMEOUT", 30*time.Second)

    // List values
    allowedHosts := cheenv.GetStringList("ALLOWED_HOSTS", ",", []string{"localhost"})
    ports := cheenv.GetIntList("PORTS", ",", []int{8080, 8081})

    fmt.Printf("Connecting to %s:%d/%s\n", dbHost, dbPort, dbName)
}
```

## Usage

### String Values

#### Get

Retrieves a string environment variable with a default value:

```go
dbHost := cheenv.Get("DB_HOST", "localhost")
dbUser := cheenv.Get("DB_USER", "postgres")

// Returns default if variable is not set
apiKey := cheenv.Get("API_KEY", "")
```

#### MustGet

Retrieves a string environment variable and panics if not set:

```go
// Panics if DB_NAME is not set
dbName := cheenv.MustGet("DB_NAME")

// Use for required configuration
apiEndpoint := cheenv.MustGet("API_ENDPOINT")
```

### Numeric Values

#### GetInt

Retrieves an integer environment variable:

```go
port := cheenv.GetInt("PORT", 8080)
maxRetries := cheenv.GetInt("MAX_RETRIES", 3)

// Returns default if variable is not set or invalid
workers := cheenv.GetInt("WORKERS", 10)
```

#### MustGetInt

Retrieves an integer and panics if not set or invalid:

```go
// Panics if PORT is not set or not a valid integer
port := cheenv.MustGetInt("PORT")
```

#### GetInt64

Retrieves a 64-bit integer:

```go
maxFileSize := cheenv.GetInt64("MAX_FILE_SIZE", 10485760) // 10MB default
userID := cheenv.GetInt64("USER_ID", 0)
```

#### GetFloat

Retrieves a floating-point value:

```go
threshold := cheenv.GetFloat("THRESHOLD", 0.95)
price := cheenv.GetFloat("PRICE", 99.99)
```

### Boolean Values

#### GetBool

Retrieves a boolean with flexible parsing:

```go
debug := cheenv.GetBool("DEBUG", false)
enableSSL := cheenv.GetBool("ENABLE_SSL", true)

// Accepts (case-insensitive):
// - true: "true", "1", "yes", "on", "y", "t"
// - false: "false", "0", "no", "off", "n", "f"
```

#### MustGetBool

Retrieves a boolean and panics if not set or invalid:

```go
// Panics if ENABLE_FEATURE is not set or not a valid boolean
enableFeature := cheenv.MustGetBool("ENABLE_FEATURE")
```

### Duration Values

#### GetDuration

Retrieves a time.Duration from strings like "5s", "10m", "2h":

```go
timeout := cheenv.GetDuration("TIMEOUT", 30*time.Second)
cacheExpiry := cheenv.GetDuration("CACHE_EXPIRY", 5*time.Minute)

// Accepts any valid Go duration string:
// - "300ms"
// - "1.5s"
// - "5m"
// - "2h30m"
```

### List Values

#### GetStringList

Retrieves a list of strings from a delimited environment variable:

```go
// Comma-separated
hosts := cheenv.GetStringList("ALLOWED_HOSTS", ",", []string{"localhost"})
// Example: ALLOWED_HOSTS="api.example.com, cdn.example.com"
// Result: []string{"api.example.com", "cdn.example.com"}

// Custom separator
paths := cheenv.GetStringList("SEARCH_PATHS", ":", []string{"/usr/local"})
// Example: SEARCH_PATHS="/opt/bin:/usr/local/bin:/usr/bin"
// Result: []string{"/opt/bin", "/usr/local/bin", "/usr/bin"}

// Whitespace is automatically trimmed
tags := cheenv.GetStringList("TAGS", ",", []string{"default"})
// Example: TAGS=" production , v2 , stable "
// Result: []string{"production", "v2", "stable"}
```

#### GetIntList

Retrieves a list of integers from a delimited environment variable:

```go
ports := cheenv.GetIntList("PORTS", ",", []int{8080})
// Example: PORTS="8080, 8081, 8082"
// Result: []int{8080, 8081, 8082}

priorityLevels := cheenv.GetIntList("PRIORITY_LEVELS", ":", []int{1, 2, 3})
// Example: PRIORITY_LEVELS="1:5:10:20"
// Result: []int{1, 5, 10, 20}

// Returns default if any value is invalid
ids := cheenv.GetIntList("IDS", ",", []int{})
// Example: IDS="1, 2, invalid, 4"
// Result: []int{} (default, because "invalid" can't be parsed)
```

### Variable Management

#### Set

Sets an environment variable:

```go
err := cheenv.Set("API_KEY", "secret-key")
if err != nil {
    log.Fatal(err)
}
```

#### Unset

Unsets an environment variable:

```go
err := cheenv.Unset("TEMP_VAR")
if err != nil {
    log.Fatal(err)
}
```

#### Has

Checks if an environment variable is set (even if empty):

```go
if cheenv.Has("API_KEY") {
    fmt.Println("API_KEY is configured")
}

// Distinguishes between unset and empty
// API_KEY="" returns true
// API_KEY not set returns false
```

### Batch Operations

#### GetAll

Returns all environment variables as a map:

```go
allVars := cheenv.GetAll()
for key, value := range allVars {
    fmt.Printf("%s=%s\n", key, value)
}
```

#### GetWithPrefix

Returns all environment variables with a specific prefix:

```go
// Get all variables starting with "APP_"
appConfig := cheenv.GetWithPrefix("APP_")

// Example environment:
// APP_NAME=MyApp
// APP_VERSION=1.0.0
// APP_DEBUG=true
// DB_HOST=localhost

// Result map:
// {
//     "NAME": "MyApp",
//     "VERSION": "1.0.0",
//     "DEBUG": "true"
// }
// (prefix "APP_" is removed from keys)

for key, value := range appConfig {
    fmt.Printf("APP_%s = %s\n", key, value)
}
```

## Examples

### Application Configuration

```go
package main

import (
    "fmt"
    "time"

    "github.com/comfortablynumb/che/pkg/cheenv"
)

type Config struct {
    // Server
    Host string
    Port int

    // Database
    DBHost     string
    DBPort     int
    DBName     string
    DBUser     string
    DBPassword string

    // Features
    EnableCache  bool
    EnableMetrics bool
    Debug        bool

    // Timeouts
    ReadTimeout  time.Duration
    WriteTimeout time.Duration

    // Lists
    AllowedOrigins []string
    TrustedProxies []string
}

func LoadConfig() *Config {
    return &Config{
        // Server
        Host: cheenv.Get("HOST", "0.0.0.0"),
        Port: cheenv.GetInt("PORT", 8080),

        // Database
        DBHost:     cheenv.Get("DB_HOST", "localhost"),
        DBPort:     cheenv.GetInt("DB_PORT", 5432),
        DBName:     cheenv.MustGet("DB_NAME"),
        DBUser:     cheenv.Get("DB_USER", "postgres"),
        DBPassword: cheenv.MustGet("DB_PASSWORD"),

        // Features
        EnableCache:   cheenv.GetBool("ENABLE_CACHE", true),
        EnableMetrics: cheenv.GetBool("ENABLE_METRICS", false),
        Debug:         cheenv.GetBool("DEBUG", false),

        // Timeouts
        ReadTimeout:  cheenv.GetDuration("READ_TIMEOUT", 30*time.Second),
        WriteTimeout: cheenv.GetDuration("WRITE_TIMEOUT", 30*time.Second),

        // Lists
        AllowedOrigins: cheenv.GetStringList("ALLOWED_ORIGINS", ",", []string{"*"}),
        TrustedProxies: cheenv.GetStringList("TRUSTED_PROXIES", ",", []string{}),
    }
}

func main() {
    config := LoadConfig()
    fmt.Printf("Starting server on %s:%d\n", config.Host, config.Port)
}
```

### Database Connection String

```go
package main

import (
    "fmt"

    "github.com/comfortablynumb/che/pkg/cheenv"
)

func BuildDatabaseURL() string {
    host := cheenv.Get("DB_HOST", "localhost")
    port := cheenv.GetInt("DB_PORT", 5432)
    user := cheenv.Get("DB_USER", "postgres")
    password := cheenv.MustGet("DB_PASSWORD")
    dbname := cheenv.MustGet("DB_NAME")
    sslMode := cheenv.Get("DB_SSLMODE", "disable")

    return fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=%s",
        user, password, host, port, dbname, sslMode,
    )
}
```

### Feature Flags

```go
package main

import (
    "github.com/comfortablynumb/che/pkg/cheenv"
)

type FeatureFlags struct {
    NewUI           bool
    BetaFeatures    bool
    ExperimentalAPI bool
    MaintenanceMode bool
}

func LoadFeatureFlags() *FeatureFlags {
    return &FeatureFlags{
        NewUI:           cheenv.GetBool("FEATURE_NEW_UI", false),
        BetaFeatures:    cheenv.GetBool("FEATURE_BETA", false),
        ExperimentalAPI: cheenv.GetBool("FEATURE_EXPERIMENTAL_API", false),
        MaintenanceMode: cheenv.GetBool("MAINTENANCE_MODE", false),
    }
}

func main() {
    flags := LoadFeatureFlags()

    if flags.MaintenanceMode {
        // Show maintenance page
    }
}
```

### Service Discovery

```go
package main

import (
    "github.com/comfortablynumb/che/pkg/cheenv"
)

func DiscoverServices() map[string]string {
    // Get all SERVICE_* environment variables
    services := cheenv.GetWithPrefix("SERVICE_")

    // Example environment:
    // SERVICE_AUTH=http://auth.internal:8080
    // SERVICE_PAYMENT=http://payment.internal:8081
    // SERVICE_NOTIFICATION=http://notification.internal:8082

    // Result:
    // {
    //     "AUTH": "http://auth.internal:8080",
    //     "PAYMENT": "http://payment.internal:8081",
    //     "NOTIFICATION": "http://notification.internal:8082"
    // }

    return services
}
```

### Rate Limiting Configuration

```go
package main

import (
    "time"

    "github.com/comfortablynumb/che/pkg/cheenv"
)

type RateLimitConfig struct {
    RequestsPerSecond int
    BurstSize         int
    WindowDuration    time.Duration
    WhitelistedIPs    []string
}

func LoadRateLimitConfig() *RateLimitConfig {
    return &RateLimitConfig{
        RequestsPerSecond: cheenv.GetInt("RATE_LIMIT_RPS", 100),
        BurstSize:         cheenv.GetInt("RATE_LIMIT_BURST", 200),
        WindowDuration:    cheenv.GetDuration("RATE_LIMIT_WINDOW", 1*time.Minute),
        WhitelistedIPs:    cheenv.GetStringList("RATE_LIMIT_WHITELIST", ",", []string{}),
    }
}
```

### Multi-Environment Configuration

```go
package main

import (
    "github.com/comfortablynumb/che/pkg/cheenv"
)

type Environment struct {
    Name       string
    LogLevel   string
    MetricsURL string
}

func DetectEnvironment() *Environment {
    env := cheenv.Get("ENV", "development")

    var config *Environment

    switch env {
    case "production":
        config = &Environment{
            Name:       "production",
            LogLevel:   cheenv.Get("LOG_LEVEL", "info"),
            MetricsURL: cheenv.MustGet("METRICS_URL"),
        }
    case "staging":
        config = &Environment{
            Name:       "staging",
            LogLevel:   cheenv.Get("LOG_LEVEL", "debug"),
            MetricsURL: cheenv.Get("METRICS_URL", "http://metrics.staging:8080"),
        }
    default:
        config = &Environment{
            Name:       "development",
            LogLevel:   cheenv.Get("LOG_LEVEL", "debug"),
            MetricsURL: cheenv.Get("METRICS_URL", "http://localhost:8080"),
        }
    }

    return config
}
```

### Worker Pool Configuration

```go
package main

import (
    "time"

    "github.com/comfortablynumb/che/pkg/cheenv"
)

type WorkerPoolConfig struct {
    NumWorkers      int
    QueueSize       int
    WorkerTimeout   time.Duration
    ShutdownTimeout time.Duration
    EnableMetrics   bool
}

func LoadWorkerPoolConfig() *WorkerPoolConfig {
    return &WorkerPoolConfig{
        NumWorkers:      cheenv.GetInt("WORKER_POOL_SIZE", 10),
        QueueSize:       cheenv.GetInt("WORKER_QUEUE_SIZE", 1000),
        WorkerTimeout:   cheenv.GetDuration("WORKER_TIMEOUT", 30*time.Second),
        ShutdownTimeout: cheenv.GetDuration("WORKER_SHUTDOWN_TIMEOUT", 10*time.Second),
        EnableMetrics:   cheenv.GetBool("WORKER_METRICS", true),
    }
}
```

## API Reference

### String Values
- `Get(key, defaultValue string) string` - Get string with default
- `MustGet(key string) string` - Get string or panic

### Numeric Values
- `GetInt(key string, defaultValue int) int` - Get int with default
- `MustGetInt(key string) int` - Get int or panic
- `GetInt64(key string, defaultValue int64) int64` - Get int64 with default
- `GetFloat(key string, defaultValue float64) float64` - Get float64 with default

### Boolean Values
- `GetBool(key string, defaultValue bool) bool` - Get bool with default
- `MustGetBool(key string) bool` - Get bool or panic

### Duration Values
- `GetDuration(key string, defaultValue time.Duration) time.Duration` - Get duration with default

### List Values
- `GetStringList(key, separator string, defaultValue []string) []string` - Get string list
- `GetIntList(key, separator string, defaultValue []int) []int` - Get int list

### Variable Management
- `Set(key, value string) error` - Set environment variable
- `Unset(key string) error` - Unset environment variable
- `Has(key string) bool` - Check if variable is set

### Batch Operations
- `GetAll() map[string]string` - Get all environment variables
- `GetWithPrefix(prefix string) map[string]string` - Get all variables with prefix

## Boolean Value Parsing

The `GetBool` and `MustGetBool` functions accept the following values (case-insensitive):

**True values**: `true`, `1`, `yes`, `on`, `y`, `t`

**False values**: `false`, `0`, `no`, `off`, `n`, `f`

Examples:
```go
// All of these evaluate to true:
cheenv.GetBool("FLAG", false) // when FLAG="true"
cheenv.GetBool("FLAG", false) // when FLAG="YES"
cheenv.GetBool("FLAG", false) // when FLAG="1"
cheenv.GetBool("FLAG", false) // when FLAG="on"

// All of these evaluate to false:
cheenv.GetBool("FLAG", true) // when FLAG="false"
cheenv.GetBool("FLAG", true) // when FLAG="NO"
cheenv.GetBool("FLAG", true) // when FLAG="0"
cheenv.GetBool("FLAG", true) // when FLAG="off"
```

## Best Practices

1. **Use Must* for Required Values**: Use `MustGet`, `MustGetInt`, etc. for configuration that is absolutely required
2. **Provide Sensible Defaults**: Always provide reasonable default values for optional configuration
3. **Validate After Loading**: Validate configuration values after loading (e.g., port ranges, positive numbers)
4. **Use Prefix Grouping**: Group related configuration with prefixes (e.g., `DB_*`, `REDIS_*`, `AWS_*`)
5. **Document Expected Variables**: Maintain a list of expected environment variables and their defaults
6. **Use .env Files in Development**: Use tools like `godotenv` to load `.env` files in development
7. **Type Safety**: Use typed getters (GetInt, GetBool, etc.) instead of parsing strings manually

## Related Packages

- **[chestring](../chestring)** - String manipulation utilities
- **[chectx](../chectx)** - Type-safe context utilities
- **[chesignal](../chesignal)** - Graceful shutdown utilities

## License

This package is part of the Che library and shares the same license.
