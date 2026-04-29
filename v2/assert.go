// assert module provides just 6 assertion functions
// that cover the most common types of assertions.
// This module is based on the articles:
// https://antonz.org/do-not-testify/
// https://www.alexedwards.net/blog/the-9-go-test-assertions-i-use
package assert

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type assertions struct {
	tb testing.TB
}

func New(t testing.TB) assertions {
	return assertions{tb: t}
}

func (a assertions) Equal(got, want any) {
	a.tb.Helper()

	if isFunction(got) || isFunction(want) {
		a.tb.Errorf("cannot compare func: got %#v want %#v", got, want)
		return
	}
	if areEqual(got, want) {
		return
	}

	errorMessage := fmt.Sprintf("got: %#v; want: %#v", got, want)
	if len(errorMessage) > 80 {
		a.tb.Errorf("not equal\ngot:\n%#v\nwant:\n%#v", got, want)
	} else {
		a.tb.Error(errorMessage)
	}
}

func (a assertions) Error(got error, want any) {
	a.tb.Helper()

	if want != nil && got == nil {
		a.tb.Errorf("got: <nil>; want: %v", want)
		return
	}
	if want == nil && got != nil {
		a.tb.Fatalf("unexpected error: %v", got)
		return
	}

	switch w := want.(type) {
	case nil:
		if got != nil {
			a.tb.Fatalf("unexpected error: %v", got)
			return
		}
	case error:
		if !errors.Is(got, w) {
			a.tb.Errorf("got: %T(%v); want: %T(%v)", got, got, w, want)
		}
	case string:
		if !strings.Contains(got.Error(), w) {
			a.tb.Errorf("got: %q; want: %q", got.Error(), w)
		}
	case reflect.Type:
		target := reflect.New(w).Interface()
		if !errors.As(got, target) {
			a.tb.Errorf("got: %T; want: %s", got, w)
		}
	default:
		a.tb.Errorf("unsupported want type: %T", want)
	}
}

func (a assertions) Nil(got any, msgAndArgs ...any) {
	a.tb.Helper()
	if !isNil(got) {
		msg := message(msgAndArgs...)
		if msg == "" {
			a.tb.Errorf("got: %#v; want <nil>", got)
		} else {
			a.tb.Errorf("got: %#v; want <nil>\n%s\n", got, msg)
		}
	}
}

func (a assertions) NotNil(got any, msgAndArgs ...any) {
	a.tb.Helper()
	if isNil(got) {
		msg := message(msgAndArgs...)
		if msg == "" {
			a.tb.Errorf("got: <nil>; want non-nil")
		} else {
			a.tb.Errorf("got: <nil>; want non-nil\n%s\n", msg)
		}
	}
}

func (a assertions) True(got bool, msgAndArgs ...any) {
	a.tb.Helper()
	if !got {
		msg := message(msgAndArgs...)
		if msg == "" {
			a.tb.Errorf("got: false; want: true")
		} else {
			a.tb.Errorf("got: false; want: true\n%s\n", msg)
		}
	}
}

func (a assertions) False(got bool, msgAndArgs ...any) {
	a.tb.Helper()
	if got {
		msg := message(msgAndArgs...)
		if msg == "" {
			a.tb.Errorf("got: true; want: false")
		} else {
			a.tb.Errorf("got: true; want: false\n%s\n", msg)
		}
	}
}

// Equal asserts that got is equal to want
func Equal[T any](tb testing.TB, got, want T) {
	tb.Helper()

	if isFunction(got) || isFunction(want) {
		tb.Errorf("cannot compare func: got %#v want %#v", got, want)
		return
	}
	if areEqual(got, want) {
		return
	}

	errorMessage := fmt.Sprintf("got: %#v; want: %#v", got, want)
	if len(errorMessage) > 80 {
		tb.Errorf("not equal\ngot:\n%#v\nwant:\n%#v", got, want)
	} else {
		tb.Error(errorMessage)
	}

}

type equaler[T any] interface {
	Equal(T) bool
}

func areEqual[T any](a, b T) bool {
	if isNil(a) && isNil(b) {
		return true
	}

	// slices of bytes
	if aBytes, ok := any(a).([]byte); ok {
		bBytes, ok := any(b).([]byte)
		return ok && bytes.Equal(aBytes, bBytes)
	}

	if aStr, ok := any(a).(string); ok {
		bStr, ok := any(b).(string)
		return ok && aStr == bStr
	}

	// compare using the type Equal function
	if eq, ok := any(a).(equaler[T]); ok {
		return eq.Equal(b)
	}

	return reflect.DeepEqual(a, b)
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	default:
		return false
	}
}

func isFunction(arg interface{}) bool {
	if arg == nil {
		return false
	}

	return reflect.TypeOf(arg).Kind() == reflect.Func
}

// Nil asserts if got is nil
func Nil(tb testing.TB, got any, msgAndArgs ...any) {
	tb.Helper()
	if !isNil(got) {
		msg := message(msgAndArgs...)
		if msg == "" {
			tb.Errorf("got: %#v; want <nil>", got)
		} else {
			tb.Errorf("got: %#v; want <nil>\n%s\n", got, msg)
		}
	}
}

// NotNil asserts got is not nil
func NotNil(tb testing.TB, got any, msgAndArgs ...any) {
	tb.Helper()
	if isNil(got) {
		msg := message(msgAndArgs...)
		if msg == "" {
			tb.Errorf("got: <nil>; want non-nil")
		} else {
			tb.Errorf("got: <nil>; want non-nil\n%s\n", msg)
		}
	}
}

// Error asserts the following cases:
// 1. If want is nil, this is equivalent to [Nil].
// 2. If want is a string, then the assertion compares the error messages.
// 3. If want is an error, then the assertion is based on [errors.Is]
// 4. If want is a type, then the assertion uses [errors.As]
func Error(tb testing.TB, got error, want any) {
	tb.Helper()

	if want != nil && got == nil {
		tb.Errorf("got: <nil>; want: %v", want)
		return
	}
	if want == nil && got != nil {
		tb.Fatalf("unexpected error: %v", got)
		return
	}

	switch w := want.(type) {
	case nil:
		if got != nil {
			tb.Fatalf("unexpected error: %v", got)
			return
		}
	case error:
		if !errors.Is(got, w) {
			tb.Errorf("got: %T(%v); want: %T(%v)", got, got, w, want)
		}
	case string:
		if !strings.Contains(got.Error(), w) {
			tb.Errorf("got: %q; want: %q", got.Error(), w)
		}
	case reflect.Type:
		target := reflect.New(w).Interface()
		if !errors.As(got, target) {
			tb.Errorf("got: %T; want: %s", got, w)
		}
	default:
		tb.Errorf("unsupported want type: %T", want)
	}
}

// True asserts got is true.
// messageAndArgs accepts a formating string and arguments to be displayed if the assertion fails
func True(tb testing.TB, got bool, messageAndArgs ...any) {
	tb.Helper()
	if !got {
		msg := message(messageAndArgs...)
		if msg == "" {
			tb.Errorf("got: false; want: true")
		} else {
			tb.Errorf("got: false; want: true\n%s\n", msg)
		}
	}
}

// False asserts got is false
// messageAndArgs accepts a formating string and arguments to be displayed if the assertion fails
func False(tb testing.TB, got bool, messageAndArgs ...any) {
	tb.Helper()
	if got {
		msg := message(messageAndArgs...)
		if msg == "" {
			tb.Errorf("got: true; want: false")
		} else {
			tb.Errorf("got: true; want: false\n%s\n", msg)
		}
	}
}

func message(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		if msg, ok := msgAndArgs[0].(string); ok {
			return msg
		}
		return fmt.Sprintf("%+v", msgAndArgs[0])
	}

	return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
}
