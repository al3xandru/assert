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

// Equal asserts that got is equal to want
func Equal[T any](tb testing.TB, got, want T) {
	tb.Helper()

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
		//fmt.Println("[]byte")
		bBytes := any(b).([]byte)
		return bytes.Equal(aBytes, bBytes)
	}

	// compare using the type Equal function
	if eq, ok := any(a).(equaler[T]); ok {
		return eq.Equal(b)
	}

	if aStr, ok := any(a).(string); ok {
		bStr := any(b).(string)
		return aStr == bStr
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

// Error asserts got is an
// If want is nil, this is equivalent to [AssertNil].
// If want is a string, then the assertion compares the error messages.
// If want is an error,
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

func True(tb testing.TB, got bool, messageAndArgs ...any) {
	tb.Helper()
	if !got {
		msg := message(messageAndArgs...)
		fmt.Println(msg)
		if msg == "" {
			tb.Errorf("got: false; want: true")
		} else {
			tb.Errorf("got: false; want: true\n%s\n", msg)
		}
	}
}

func False(tb testing.TB, got bool, messageAndArgs ...any) {
	tb.Helper()
	if got {
		msg := message(messageAndArgs...)
		fmt.Println(msg)
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
