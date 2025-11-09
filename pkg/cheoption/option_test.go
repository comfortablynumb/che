package cheoption

import (
	"errors"
	"testing"
)

func TestOptional_Some(t *testing.T) {
	opt := Some(42)
	if !opt.IsPresent() {
		t.Error("Some should be present")
	}
	if opt.Get() != 42 {
		t.Errorf("expected 42, got %d", opt.Get())
	}
}

func TestOptional_None(t *testing.T) {
	opt := None[int]()
	if opt.IsPresent() {
		t.Error("None should not be present")
	}
	if !opt.IsEmpty() {
		t.Error("None should be empty")
	}
}

func TestOptional_GetOr(t *testing.T) {
	some := Some(42)
	none := None[int]()

	if some.GetOr(0) != 42 {
		t.Error("Some.GetOr should return value")
	}
	if none.GetOr(99) != 99 {
		t.Error("None.GetOr should return default")
	}
}

func TestOptional_Map(t *testing.T) {
	opt := Some(5).Map(func(x int) int { return x * 2 })
	if opt.Get() != 10 {
		t.Errorf("expected 10, got %d", opt.Get())
	}

	none := None[int]().Map(func(x int) int { return x * 2 })
	if none.IsPresent() {
		t.Error("Map on None should return None")
	}
}

func TestOptional_Filter(t *testing.T) {
	opt := Some(5).Filter(func(x int) bool { return x > 3 })
	if !opt.IsPresent() {
		t.Error("Filter should keep value > 3")
	}

	opt = Some(2).Filter(func(x int) bool { return x > 3 })
	if opt.IsPresent() {
		t.Error("Filter should remove value <= 3")
	}
}

func TestResult_Ok(t *testing.T) {
	result := Ok(42)
	if !result.IsOk() {
		t.Error("Ok should be ok")
	}
	if result.Unwrap() != 42 {
		t.Errorf("expected 42, got %d", result.Unwrap())
	}
}

func TestResult_Err(t *testing.T) {
	err := errors.New("test error")
	result := Err[int](err)
	if !result.IsErr() {
		t.Error("Err should be error")
	}
	if result.Error() != err {
		t.Error("Error() should return the error")
	}
}

func TestResult_UnwrapOr(t *testing.T) {
	ok := Ok(42)
	err := Err[int](errors.New("test"))

	if ok.UnwrapOr(0) != 42 {
		t.Error("Ok.UnwrapOr should return value")
	}
	if err.UnwrapOr(99) != 99 {
		t.Error("Err.UnwrapOr should return default")
	}
}

func TestResult_Map(t *testing.T) {
	result := Ok(5).Map(func(x int) int { return x * 2 })
	if result.Unwrap() != 10 {
		t.Errorf("expected 10, got %d", result.Unwrap())
	}

	err := Err[int](errors.New("test")).Map(func(x int) int { return x * 2 })
	if err.IsOk() {
		t.Error("Map on Err should return Err")
	}
}

func TestResult_FlatMap(t *testing.T) {
	result := Ok(5).FlatMap(func(x int) Result[int] {
		if x > 0 {
			return Ok(x * 2)
		}
		return Err[int](errors.New("negative"))
	})

	if result.Unwrap() != 10 {
		t.Errorf("expected 10, got %d", result.Unwrap())
	}
}
