# chefile

File utilities for Go with support for common operations and structured data formats.

## Features

- File copy and move operations with permission preservation
- Atomic file writes (write to temp file, then rename)
- File existence and type checking
- Human-readable file size formatting
- JSON, YAML, and CSV file operations
- Zero external dependencies (except gopkg.in/yaml.v3 for YAML support)

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chefile
```

## Usage

### Basic File Operations

#### Copy Files

```go
package main

import (
    "github.com/comfortablynumb/che/pkg/chefile"
)

func main() {
    // Copy file with permission preservation
    err := chefile.Copy("source.txt", "dest.txt")
    if err != nil {
        panic(err)
    }
}
```

#### Move Files

```go
// Move file (tries os.Rename, falls back to copy+delete)
err := chefile.Move("old.txt", "new.txt")
if err != nil {
    panic(err)
}
```

#### Atomic Writes

```go
// Write data atomically (temp file + rename)
data := []byte("important data")
err := chefile.AtomicWrite("config.txt", data, 0644)
if err != nil {
    panic(err)
}
```

### File Information

#### Check Existence

```go
if chefile.Exists("file.txt") {
    fmt.Println("File exists")
}

if chefile.IsDir("mydir") {
    fmt.Println("It's a directory")
}

if chefile.IsFile("file.txt") {
    fmt.Println("It's a regular file")
}
```

#### Get File Size

```go
size, err := chefile.Size("largefile.bin")
if err != nil {
    panic(err)
}

// Format size in human-readable form
formatted := chefile.FormatSize(size)
fmt.Println(formatted) // e.g., "1.5 MB"
```

### JSON Operations

#### Read JSON

```go
type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

var config Config
err := chefile.ReadJSON("config.json", &config)
if err != nil {
    panic(err)
}

fmt.Printf("Server: %s:%d\n", config.Host, config.Port)
```

#### Write JSON

```go
config := Config{
    Host: "localhost",
    Port: 8080,
}

// Write compact JSON
err := chefile.WriteJSON("config.json", config, 0644)

// Write indented JSON (pretty-printed)
err = chefile.WriteJSONIndent("config.json", config, 0644)
```

### YAML Operations

#### Read YAML

```go
type Config struct {
    Database struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"database"`
}

var config Config
err := chefile.ReadYAML("config.yaml", &config)
if err != nil {
    panic(err)
}

fmt.Printf("DB: %s:%d\n", config.Database.Host, config.Database.Port)
```

#### Write YAML

```go
type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
}

config := Config{}
config.Server.Host = "localhost"
config.Server.Port = 8080

err := chefile.WriteYAML("config.yaml", config, 0644)
if err != nil {
    panic(err)
}
```

### CSV Operations

#### Read CSV

```go
records, err := chefile.ReadCSV("data.csv")
if err != nil {
    panic(err)
}

// Process records
for i, record := range records {
    if i == 0 {
        // Header row
        fmt.Println("Headers:", record)
        continue
    }
    fmt.Printf("Row %d: %v\n", i, record)
}
```

#### Write CSV

```go
records := [][]string{
    {"Name", "Age", "City"},
    {"Alice", "25", "NYC"},
    {"Bob", "30", "LA"},
    {"Charlie", "35", "SF"},
}

err := chefile.WriteCSV("output.csv", records, 0644)
if err != nil {
    panic(err)
}
```

### Complete Examples

#### Config File Manager

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chefile"
)

type AppConfig struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Features map[string]bool `yaml:"features"`
}

type ServerConfig struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

type DatabaseConfig struct {
    URL      string `yaml:"url"`
    MaxConns int    `yaml:"max_conns"`
}

func main() {
    configFile := "app.yaml"

    // Read config
    var config AppConfig
    if chefile.Exists(configFile) {
        if err := chefile.ReadYAML(configFile, &config); err != nil {
            panic(err)
        }
    } else {
        // Create default config
        config = AppConfig{
            Server: ServerConfig{
                Host: "localhost",
                Port: 8080,
            },
            Database: DatabaseConfig{
                URL:      "postgres://localhost/myapp",
                MaxConns: 10,
            },
            Features: map[string]bool{
                "feature_a": true,
                "feature_b": false,
            },
        }

        // Save default config
        if err := chefile.WriteYAML(configFile, config, 0644); err != nil {
            panic(err)
        }
    }

    fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)
}
```

#### CSV Data Processor

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/comfortablynumb/che/pkg/chefile"
)

func main() {
    // Read input CSV
    records, err := chefile.ReadCSV("sales.csv")
    if err != nil {
        panic(err)
    }

    // Process data
    var total float64
    output := [][]string{{"Product", "Quantity", "Price", "Total"}}

    for i, record := range records {
        if i == 0 {
            continue // Skip header
        }

        product := record[0]
        qty, _ := strconv.ParseFloat(record[1], 64)
        price, _ := strconv.ParseFloat(record[2], 64)
        lineTotal := qty * price
        total += lineTotal

        output = append(output, []string{
            product,
            record[1],
            record[2],
            fmt.Sprintf("%.2f", lineTotal),
        })
    }

    // Add total row
    output = append(output, []string{"", "", "TOTAL", fmt.Sprintf("%.2f", total)})

    // Write output
    if err := chefile.WriteCSV("sales_report.csv", output, 0644); err != nil {
        panic(err)
    }

    fmt.Printf("Processed %d records. Total: $%.2f\n", len(records)-1, total)
}
```

#### Backup Manager

```go
package main

import (
    "fmt"
    "time"
    "github.com/comfortablynumb/che/pkg/chefile"
)

func backupFile(source string) error {
    if !chefile.Exists(source) {
        return fmt.Errorf("source file does not exist: %s", source)
    }

    // Create backup filename with timestamp
    timestamp := time.Now().Format("20060102-150405")
    backup := fmt.Sprintf("%s.backup-%s", source, timestamp)

    // Copy to backup
    if err := chefile.Copy(source, backup); err != nil {
        return fmt.Errorf("failed to create backup: %w", err)
    }

    size, _ := chefile.Size(backup)
    fmt.Printf("Created backup: %s (%s)\n", backup, chefile.FormatSize(size))
    return nil
}

func main() {
    files := []string{"config.yaml", "data.json", "users.csv"}

    for _, file := range files {
        if err := backupFile(file); err != nil {
            fmt.Printf("Error backing up %s: %v\n", file, err)
        }
    }
}
```

## API Reference

### File Operations

- `Copy(src, dst string) error` - Copy file with permission preservation
- `Move(src, dst string) error` - Move file (rename or copy+delete)
- `AtomicWrite(path string, data []byte, perm os.FileMode) error` - Atomic write via temp file

### File Information

- `Exists(path string) bool` - Check if file or directory exists
- `IsDir(path string) bool` - Check if path is a directory
- `IsFile(path string) bool` - Check if path is a regular file
- `Size(path string) (int64, error)` - Get file size in bytes
- `FormatSize(bytes int64) string` - Format bytes as human-readable string

### JSON Operations

- `ReadJSON(path string, v interface{}) error` - Read and unmarshal JSON file
- `WriteJSON(path string, v interface{}, perm os.FileMode) error` - Marshal and write JSON (compact)
- `WriteJSONIndent(path string, v interface{}, perm os.FileMode) error` - Marshal and write JSON (indented)

### YAML Operations

- `ReadYAML(path string, v interface{}) error` - Read and unmarshal YAML file
- `WriteYAML(path string, v interface{}, perm os.FileMode) error` - Marshal and write YAML file

### CSV Operations

- `ReadCSV(path string) ([][]string, error)` - Read all CSV records
- `WriteCSV(path string, records [][]string, perm os.FileMode) error` - Write CSV records atomically

## Features

### Atomic Writes

All write operations (JSON, YAML, CSV) use atomic writes by default:
1. Create temporary file in same directory
2. Write data to temp file
3. Set correct permissions
4. Rename temp file to target (atomic operation)
5. Cleanup temp file on error

This ensures your files are never partially written or corrupted.

### Permission Preservation

Copy operations preserve file permissions from the source file.

### Error Handling

All functions return descriptive errors with context:
- `failed to open source file: <error>`
- `failed to unmarshal JSON: <error>`
- `failed to write CSV: <error>`

This makes debugging easier by providing clear error messages.

## Size Formatting

The `FormatSize` function converts bytes to human-readable format:

```go
chefile.FormatSize(500)           // "500 B"
chefile.FormatSize(1024)          // "1.0 KB"
chefile.FormatSize(1536)          // "1.5 KB"
chefile.FormatSize(1048576)       // "1.0 MB"
chefile.FormatSize(1073741824)    // "1.0 GB"
chefile.FormatSize(1099511627776) // "1.0 TB"
```

## Dependencies

- Standard library packages: `os`, `io`, `encoding/json`, `encoding/csv`, `fmt`, `path/filepath`
- External: `gopkg.in/yaml.v3` (for YAML support only)

## Testing

The package includes comprehensive tests covering:
- All file operations
- Error cases (missing files, invalid data)
- Permission handling
- Atomic write behavior
- JSON/YAML/CSV round-trip tests

Run tests with:
```bash
go test -v github.com/comfortablynumb/che/pkg/chefile
```

## License

MIT
