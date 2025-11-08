package chectx

import (
	"context"
	"testing"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestKey(t *testing.T) {
	key1 := Key[string]("user")
	key2 := Key[string]("user")

	// Keys with the same name should be different instances
	chetest.RequireEqual(t, key1 != key2, true)
	chetest.RequireEqual(t, key1.name, "user")
	chetest.RequireEqual(t, key2.name, "user")
}

func TestWithValue_And_Value(t *testing.T) {
	key := Key[string]("username")
	ctx := context.Background()

	// Store value
	ctx = WithValue(ctx, key, "john_doe")

	// Retrieve value
	val, ok := Value(ctx, key)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, "john_doe")
}

func TestValue_NotFound(t *testing.T) {
	key := Key[string]("username")
	ctx := context.Background()

	// Try to retrieve value that doesn't exist
	val, ok := Value(ctx, key)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, val, "")
}

func TestValue_DifferentTypes(t *testing.T) {
	stringKey := Key[string]("value")
	intKey := Key[int]("value")

	ctx := context.Background()
	ctx = WithValue(ctx, stringKey, "hello")
	ctx = WithValue(ctx, intKey, 42)

	// Retrieve string value
	strVal, ok := Value(ctx, stringKey)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, strVal, "hello")

	// Retrieve int value
	intVal, ok := Value(ctx, intKey)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, intVal, 42)
}

func TestValue_ComplexTypes(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	key := Key[*User]("user")
	ctx := context.Background()

	user := &User{ID: 1, Name: "John"}
	ctx = WithValue(ctx, key, user)

	// Retrieve complex type
	val, ok := Value(ctx, key)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val != nil, true)
	chetest.RequireEqual(t, val.ID, 1)
	chetest.RequireEqual(t, val.Name, "John")
}

func TestMustValue_Success(t *testing.T) {
	key := Key[string]("username")
	ctx := context.Background()
	ctx = WithValue(ctx, key, "jane_doe")

	// Should not panic
	val := MustValue(ctx, key)
	chetest.RequireEqual(t, val, "jane_doe")
}

func TestMustValue_Panic(t *testing.T) {
	key := Key[string]("username")
	ctx := context.Background()

	// Should panic
	defer func() {
		r := recover()
		chetest.RequireEqual(t, r != nil, true)
	}()

	_ = MustValue(ctx, key)
}

func TestGetOrDefault_Found(t *testing.T) {
	key := Key[int]("count")
	ctx := context.Background()
	ctx = WithValue(ctx, key, 42)

	val := GetOrDefault(ctx, key, 10)
	chetest.RequireEqual(t, val, 42)
}

func TestGetOrDefault_NotFound(t *testing.T) {
	key := Key[int]("count")
	ctx := context.Background()

	val := GetOrDefault(ctx, key, 10)
	chetest.RequireEqual(t, val, 10)
}

func TestGetOrDefault_ZeroValue(t *testing.T) {
	key := Key[int]("count")
	ctx := context.Background()
	ctx = WithValue(ctx, key, 0)

	// Should return the stored zero value, not the default
	val := GetOrDefault(ctx, key, 10)
	chetest.RequireEqual(t, val, 0)
}

func TestContextChaining(t *testing.T) {
	userKey := Key[string]("user")
	roleKey := Key[string]("role")
	idKey := Key[int]("id")

	ctx := context.Background()
	ctx = WithValue(ctx, userKey, "john")
	ctx = WithValue(ctx, roleKey, "admin")
	ctx = WithValue(ctx, idKey, 123)

	// All values should be retrievable
	user, ok := Value(ctx, userKey)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, user, "john")

	role, ok := Value(ctx, roleKey)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, role, "admin")

	id, ok := Value(ctx, idKey)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, id, 123)
}

func TestKeyIsolation(t *testing.T) {
	key1 := Key[string]("test")
	key2 := Key[string]("test")

	ctx := context.Background()
	ctx = WithValue(ctx, key1, "value1")

	// key2 should not retrieve value1 even though they have the same name
	val, ok := Value(ctx, key2)
	chetest.RequireEqual(t, ok, false)
	chetest.RequireEqual(t, val, "")

	// key1 should still work
	val, ok = Value(ctx, key1)
	chetest.RequireEqual(t, ok, true)
	chetest.RequireEqual(t, val, "value1")
}
