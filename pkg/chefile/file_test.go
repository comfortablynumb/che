package chefile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestCopy(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.txt")
	dst := filepath.Join(tmpDir, "dest.txt")

	// Create source file
	err := os.WriteFile(src, []byte("test content"), 0644)
	chetest.RequireEqual(t, err, nil)

	// Copy file
	err = Copy(src, dst)
	chetest.RequireEqual(t, err, nil)

	// Verify content
	data, err := os.ReadFile(dst)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, string(data), "test content")

	// Verify permissions
	srcInfo, _ := os.Stat(src)
	dstInfo, _ := os.Stat(dst)
	chetest.RequireEqual(t, dstInfo.Mode(), srcInfo.Mode())
}

func TestCopy_NonExistentSource(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "nonexistent.txt")
	dst := filepath.Join(tmpDir, "dest.txt")

	err := Copy(src, dst)
	chetest.RequireEqual(t, err != nil, true)
}

func TestMove(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.txt")
	dst := filepath.Join(tmpDir, "dest.txt")

	// Create source file
	err := os.WriteFile(src, []byte("test content"), 0644)
	chetest.RequireEqual(t, err, nil)

	// Move file
	err = Move(src, dst)
	chetest.RequireEqual(t, err, nil)

	// Verify destination exists
	chetest.RequireEqual(t, Exists(dst), true)

	// Verify source no longer exists
	chetest.RequireEqual(t, Exists(src), false)

	// Verify content
	data, err := os.ReadFile(dst)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, string(data), "test content")
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.txt")
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")

	// Create file
	err := os.WriteFile(existingFile, []byte("test"), 0644)
	chetest.RequireEqual(t, err, nil)

	chetest.RequireEqual(t, Exists(existingFile), true)
	chetest.RequireEqual(t, Exists(nonExistentFile), false)
	chetest.RequireEqual(t, Exists(tmpDir), true)
}

func TestIsDir(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "file.txt")
	dir := filepath.Join(tmpDir, "subdir")

	// Create file and directory
	err := os.WriteFile(file, []byte("test"), 0644)
	chetest.RequireEqual(t, err, nil)
	err = os.Mkdir(dir, 0755)
	chetest.RequireEqual(t, err, nil)

	chetest.RequireEqual(t, IsDir(tmpDir), true)
	chetest.RequireEqual(t, IsDir(dir), true)
	chetest.RequireEqual(t, IsDir(file), false)
	chetest.RequireEqual(t, IsDir("/nonexistent"), false)
}

func TestIsFile(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "file.txt")
	dir := filepath.Join(tmpDir, "subdir")

	// Create file and directory
	err := os.WriteFile(file, []byte("test"), 0644)
	chetest.RequireEqual(t, err, nil)
	err = os.Mkdir(dir, 0755)
	chetest.RequireEqual(t, err, nil)

	chetest.RequireEqual(t, IsFile(file), true)
	chetest.RequireEqual(t, IsFile(dir), false)
	chetest.RequireEqual(t, IsFile(tmpDir), false)
	chetest.RequireEqual(t, IsFile("/nonexistent"), false)
}

func TestSize(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "file.txt")

	content := "test content"
	err := os.WriteFile(file, []byte(content), 0644)
	chetest.RequireEqual(t, err, nil)

	size, err := Size(file)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, size, int64(len(content)))
}

func TestSize_NonExistent(t *testing.T) {
	_, err := Size("/nonexistent/file.txt")
	chetest.RequireEqual(t, err != nil, true)
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tt := range tests {
		result := FormatSize(tt.bytes)
		chetest.RequireEqual(t, result, tt.expected)
	}
}

func TestAtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "atomic.txt")

	data := []byte("atomic content")
	err := AtomicWrite(file, data, 0644)
	chetest.RequireEqual(t, err, nil)

	// Verify content
	readData, err := os.ReadFile(file)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, string(readData), string(data))

	// Verify permissions
	info, err := os.Stat(file)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, info.Mode().Perm(), os.FileMode(0644))
}

func TestReadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.json")

	// Write test JSON
	jsonData := `{"name":"John","age":30,"active":true}`
	err := os.WriteFile(file, []byte(jsonData), 0644)
	chetest.RequireEqual(t, err, nil)

	// Read JSON
	type Person struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Active bool   `json:"active"`
	}

	var person Person
	err = ReadJSON(file, &person)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, person.Name, "John")
	chetest.RequireEqual(t, person.Age, 30)
	chetest.RequireEqual(t, person.Active, true)
}

func TestReadJSON_InvalidFile(t *testing.T) {
	var data map[string]interface{}
	err := ReadJSON("/nonexistent/file.json", &data)
	chetest.RequireEqual(t, err != nil, true)
}

func TestReadJSON_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "invalid.json")

	err := os.WriteFile(file, []byte("not valid json"), 0644)
	chetest.RequireEqual(t, err, nil)

	var data map[string]interface{}
	err = ReadJSON(file, &data)
	chetest.RequireEqual(t, err != nil, true)
}

func TestWriteJSON(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.json")

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	person := Person{Name: "Alice", Age: 25}
	err := WriteJSON(file, person, 0644)
	chetest.RequireEqual(t, err, nil)

	// Read back and verify
	var result Person
	err = ReadJSON(file, &result)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, result.Name, "Alice")
	chetest.RequireEqual(t, result.Age, 25)
}

func TestWriteJSONIndent(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.json")

	data := map[string]interface{}{
		"name": "Bob",
		"age":  30,
	}

	err := WriteJSONIndent(file, data, 0644)
	chetest.RequireEqual(t, err, nil)

	// Read raw content to verify indentation
	content, err := os.ReadFile(file)
	chetest.RequireEqual(t, err, nil)

	// Should contain newlines and spaces for indentation
	contentStr := string(content)
	chetest.RequireEqual(t, contentStr != "{\"age\":30,\"name\":\"Bob\"}", true)

	// Verify it's still valid JSON
	var result map[string]interface{}
	err = ReadJSON(file, &result)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, result["name"], "Bob")
}

func TestReadYAML(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.yaml")

	yamlData := `name: John
age: 30
active: true
`
	err := os.WriteFile(file, []byte(yamlData), 0644)
	chetest.RequireEqual(t, err, nil)

	type Person struct {
		Name   string `yaml:"name"`
		Age    int    `yaml:"age"`
		Active bool   `yaml:"active"`
	}

	var person Person
	err = ReadYAML(file, &person)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, person.Name, "John")
	chetest.RequireEqual(t, person.Age, 30)
	chetest.RequireEqual(t, person.Active, true)
}

func TestReadYAML_InvalidFile(t *testing.T) {
	var data map[string]interface{}
	err := ReadYAML("/nonexistent/file.yaml", &data)
	chetest.RequireEqual(t, err != nil, true)
}

func TestReadYAML_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "invalid.yaml")

	err := os.WriteFile(file, []byte(":\ninvalid: yaml: content:"), 0644)
	chetest.RequireEqual(t, err, nil)

	var data map[string]interface{}
	err = ReadYAML(file, &data)
	chetest.RequireEqual(t, err != nil, true)
}

func TestWriteYAML(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.yaml")

	type Person struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	person := Person{Name: "Alice", Age: 25}
	err := WriteYAML(file, person, 0644)
	chetest.RequireEqual(t, err, nil)

	// Read back and verify
	var result Person
	err = ReadYAML(file, &result)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, result.Name, "Alice")
	chetest.RequireEqual(t, result.Age, 25)
}

func TestReadCSV(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.csv")

	csvData := `name,age,city
Alice,25,NYC
Bob,30,LA
Charlie,35,SF
`
	err := os.WriteFile(file, []byte(csvData), 0644)
	chetest.RequireEqual(t, err, nil)

	records, err := ReadCSV(file)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, len(records), 4)
	chetest.RequireEqual(t, len(records[0]), 3)
	chetest.RequireEqual(t, records[0][0], "name")
	chetest.RequireEqual(t, records[1][0], "Alice")
	chetest.RequireEqual(t, records[1][1], "25")
	chetest.RequireEqual(t, records[2][0], "Bob")
}

func TestReadCSV_InvalidFile(t *testing.T) {
	_, err := ReadCSV("/nonexistent/file.csv")
	chetest.RequireEqual(t, err != nil, true)
}

func TestReadCSV_InvalidCSV(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "invalid.csv")

	// CSV with inconsistent columns
	csvData := `a,b,c
1,2
3,4,5,6
`
	err := os.WriteFile(file, []byte(csvData), 0644)
	chetest.RequireEqual(t, err, nil)

	_, err = ReadCSV(file)
	chetest.RequireEqual(t, err != nil, true)
}

func TestWriteCSV(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.csv")

	records := [][]string{
		{"name", "age", "city"},
		{"Alice", "25", "NYC"},
		{"Bob", "30", "LA"},
	}

	err := WriteCSV(file, records, 0644)
	chetest.RequireEqual(t, err, nil)

	// Read back and verify
	result, err := ReadCSV(file)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, len(result), 3)
	chetest.RequireEqual(t, result[0][0], "name")
	chetest.RequireEqual(t, result[1][0], "Alice")
	chetest.RequireEqual(t, result[2][0], "Bob")
}

func TestWriteCSV_EmptyRecords(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "empty.csv")

	records := [][]string{}
	err := WriteCSV(file, records, 0644)
	chetest.RequireEqual(t, err, nil)

	result, err := ReadCSV(file)
	chetest.RequireEqual(t, err, nil)
	chetest.RequireEqual(t, len(result), 0)
}
