# chestrsim

String similarity and fuzzy matching algorithms for Go.

## Features

- Multiple similarity algorithms
- Unicode-aware (handles UTF-8 properly)
- Fuzzy matching for search
- Efficient implementations
- Zero external dependencies

## Installation

```bash
go get github.com/comfortablynumb/che/pkg/chestrsim
```

## Algorithms

### Levenshtein Distance

Edit distance between two strings (insertions, deletions, substitutions).

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chestrsim"
)

func main() {
    distance := chestrsim.Levenshtein("kitten", "sitting")
    fmt.Println(distance) // 3

    similarity := chestrsim.LevenshteinSimilarity("kitten", "sitting")
    fmt.Printf("Similarity: %.2f\n", similarity) // 0.57
}
```

### Hamming Distance

Number of positions at which characters differ (equal length strings only).

```go
distance := chestrsim.Hamming("1011101", "1001001")
fmt.Println(distance) // 2

// Returns -1 for different lengths
distance = chestrsim.Hamming("abc", "ab")
fmt.Println(distance) // -1

similarity := chestrsim.HammingSimilarity("abc", "axc")
fmt.Printf("Similarity: %.2f\n", similarity) // 0.67
```

### Jaro-Winkler Similarity

Good for short strings like names. Gives more weight to common prefixes.

```go
similarity := chestrsim.JaroWinkler("martha", "marhta")
fmt.Printf("Similarity: %.3f\n", similarity) // 0.961

similarity = chestrsim.JaroWinkler("dwayne", "duane")
fmt.Printf("Similarity: %.3f\n", similarity) // 0.840

// Plain Jaro (without prefix bonus)
similarity = chestrsim.Jaro("martha", "marhta")
fmt.Printf("Similarity: %.3f\n", similarity) // 0.944
```

### Cosine Similarity

Based on character bigrams (pairs of consecutive characters).

```go
similarity := chestrsim.Cosine("hello", "hallo")
fmt.Printf("Similarity: %.2f\n", similarity) // ~0.50

similarity = chestrsim.Cosine("data", "date")
fmt.Printf("Similarity: %.2f\n", similarity) // ~0.67
```

### Jaccard Similarity

Ratio of intersection to union of character bigrams.

```go
similarity := chestrsim.Jaccard("hello", "hallo")
fmt.Printf("Similarity: %.2f\n", similarity) // ~0.40

similarity = chestrsim.Jaccard("data", "data")
fmt.Printf("Similarity: %.2f\n", similarity) // 1.00
```

### Fuzzy Matching

Check if query characters appear in order within target string.

```go
// Basic fuzzy match
matches := chestrsim.FuzzyMatch("fb", "FooBar")
fmt.Println(matches) // true (case-insensitive)

matches = chestrsim.FuzzyMatch("fb", "foobar")
fmt.Println(matches) // true

matches = chestrsim.FuzzyMatch("abc", "axbxcx")
fmt.Println(matches) // true

matches = chestrsim.FuzzyMatch("abc", "acb")
fmt.Println(matches) // false (order matters)

// Fuzzy score (higher is better match)
score := chestrsim.FuzzyScore("fb", "foobar")
fmt.Printf("Score: %.2f\n", score) // Higher for consecutive matches
```

## Use Cases

### Finding Similar Names

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chestrsim"
)

func main() {
    users := []string{"John Smith", "Jon Smyth", "Jane Doe", "John Smythe"}
    query := "John Smith"

    for _, user := range users {
        sim := chestrsim.JaroWinkler(query, user)
        if sim > 0.8 {
            fmt.Printf("%s: %.3f\n", user, sim)
        }
    }
}
// Output:
// John Smith: 1.000
// Jon Smyth: 0.933
// John Smythe: 0.933
```

### Fuzzy Search

```go
package main

import (
    "fmt"
    "github.com/comfortablynumb/che/pkg/chestrsim"
)

type SearchResult struct {
    Text  string
    Score float64
}

func FuzzySearch(query string, items []string) []SearchResult {
    var results []SearchResult

    for _, item := range items {
        score := chestrsim.FuzzyScore(query, item)
        if score > 0 {
            results = append(results, SearchResult{item, score})
        }
    }

    // Sort by score (descending)
    sort.Slice(results, func(i, j int) bool {
        return results[i].Score > results[j].Score
    })

    return results
}

func main() {
    files := []string{
        "config.yaml",
        "controller.go",
        "component.tsx",
        "constants.go",
    }

    results := FuzzySearch("cg", files)
    for _, r := range results {
        fmt.Printf("%s: %.2f\n", r.Text, r.Score)
    }
}
// Output:
// config.yaml: 1.00
// controller.go: 0.50
// component.tsx: 0.33
```

### Deduplication

```go
func FindDuplicates(items []string, threshold float64) [][]string {
    var groups [][]string
    used := make(map[int]bool)

    for i, item1 := range items {
        if used[i] {
            continue
        }

        group := []string{item1}
        used[i] = true

        for j, item2 := range items[i+1:] {
            if used[i+1+j] {
                continue
            }

            sim := chestrsim.LevenshteinSimilarity(item1, item2)
            if sim >= threshold {
                group = append(group, item2)
                used[i+1+j] = true
            }
        }

        if len(group) > 1 {
            groups = append(groups, group)
        }
    }

    return groups
}
```

### Spell Checking

```go
func SuggestCorrection(word string, dictionary []string, maxSuggestions int) []string {
    type suggestion struct {
        word     string
        distance int
    }

    var suggestions []suggestion

    for _, dictWord := range dictionary {
        dist := chestrsim.Levenshtein(word, dictWord)
        if dist <= 2 { // Only suggest if close enough
            suggestions = append(suggestions, suggestion{dictWord, dist})
        }
    }

    // Sort by distance
    sort.Slice(suggestions, func(i, j int) bool {
        return suggestions[i].distance < suggestions[j].distance
    })

    // Return top N
    result := []string{}
    for i := 0; i < len(suggestions) && i < maxSuggestions; i++ {
        result = append(result, suggestions[i].word)
    }

    return result
}
```

## Algorithm Comparison

| Algorithm | Use Case | Strength | Time Complexity |
|-----------|----------|----------|----------------|
| Levenshtein | General similarity | Most accurate edit distance | O(m×n) |
| Hamming | Fixed-length comparison | Fast, simple | O(n) |
| Jaro-Winkler | Short strings, names | Good for typos with common prefix | O(m×n) |
| Cosine | Document similarity | Character distribution | O(m+n) |
| Jaccard | Set similarity | Simple, intuitive | O(m+n) |
| Fuzzy Match | Search/autocomplete | Fast substring matching | O(m+n) |

## Performance

All algorithms are optimized for:
- Unicode/UTF-8 strings
- Minimal memory allocation
- Efficient computation

Benchmarks included in test files.

## License

MIT
