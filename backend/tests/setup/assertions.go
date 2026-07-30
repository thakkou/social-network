package setup

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// AssertEqual fails the test if expected and actual are not equal.
func AssertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\texpected: %v\n\tactual:   %v", msg, expected, actual)
	}
}

// AssertNotEqual fails the test if expected and actual are equal.
func AssertNotEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\tunexpected value: %v", msg, expected)
	}
}

// AssertTrue fails the test if condition is false.
func AssertTrue(t *testing.T, condition bool, msgAndArgs ...interface{}) {
	t.Helper()
	if !condition {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\texpected true, got false", msg)
	}
}

// AssertFalse fails the test if condition is true.
func AssertFalse(t *testing.T, condition bool, msgAndArgs ...interface{}) {
	t.Helper()
	if condition {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\texpected false, got true", msg)
	}
}

// AssertNil fails the test if value is not nil.
func AssertNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value != nil && !reflect.ValueOf(value).IsNil() {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\texpected nil, got: %v", msg, value)
	}
}

// AssertNotNil fails the test if value is nil.
func AssertNotNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value == nil || reflect.ValueOf(value).IsNil() {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\texpected non-nil value", msg)
	}
}

// AssertResponseOK checks that the response wrapper indicates success (2xx).
func AssertResponseOK(t *testing.T, wrapper *ResponseWrapper, msgAndArgs ...interface{}) {
	t.Helper()
	if wrapper.StatusCode < 200 || wrapper.StatusCode >= 300 {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\texpected successful response, got status %d: %s", msg, wrapper.StatusCode, wrapper.Message)
	}
}

// AssertResponseError checks that the response wrapper indicates an error (4xx/5xx).
func AssertResponseError(t *testing.T, wrapper *ResponseWrapper, expectedStatus int, msgAndArgs ...interface{}) {
	t.Helper()
	if wrapper.StatusCode != expectedStatus {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\texpected status %d, got %d: %s", msg, expectedStatus, wrapper.StatusCode, wrapper.Message)
	}
}

// AssertMessage checks that the response message matches.
func AssertMessage(t *testing.T, wrapper *ResponseWrapper, expectedMsg string, msgAndArgs ...interface{}) {
	t.Helper()
	if wrapper.Message != expectedMsg {
		msg := formatMessage(msgAndArgs...)
		t.Fatalf("%s\n\texpected message %q, got %q", msg, expectedMsg, wrapper.Message)
	}
}

// MustMarshalJSON marshals the value to JSON or fails the test.
func MustMarshalJSON(t *testing.T, v interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}
	return b
}

// MustUnmarshalJSON unmarshals JSON into the given target or fails the test.
func MustUnmarshalJSON(t *testing.T, data []byte, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("failed to unmarshal JSON %s: %v", string(data), err)
	}
}

// formatMessage formats optional message and args similar to testing.T.Logf.
func formatMessage(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return "assertion failed"
	}
	if len(msgAndArgs) == 1 {
		return fmt.Sprintf("%v", msgAndArgs[0])
	}
	msg, ok := msgAndArgs[0].(string)
	if !ok {
		return fmt.Sprintf("%v", msgAndArgs)
	}
	return fmt.Sprintf(msg, msgAndArgs[1:]...)
}
