//go:build !solution

package testequal

import (
	"bytes"
	"fmt"
	"reflect"
)

func equal(expected, actual any) bool {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return false
	}

	if expected == nil || actual == nil {
		return expected == nil && actual == nil
	}

	switch e := expected.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return expected == actual

	case struct{}:
		return false

	case map[string]string:
		actualMap, err := actual.(map[string]string)

		if !err {
			return false
		}

		if len(e) != len(actualMap) {
			return false
		}

		if len(actualMap) == 0 {
			return false
		}

		for key, expectedValue := range e {
			var actualValue string
			actualValue, err = actualMap[key]

			if !err || expectedValue != actualValue {
				return false
			}
		}

		return true
	}

	return equalForBytes(expected, actual)
}

func equalForBytes(expected, actual any) bool {
	ExpectedBytes, errForExpected := expected.([]byte)
	ActualBytes, okActual := actual.([]byte)

	if !errForExpected || !okActual {
		return reflect.DeepEqual(expected, actual)
	}

	if ExpectedBytes == nil || ActualBytes == nil {
		return ExpectedBytes == nil && ActualBytes == nil
	}

	return bytes.Equal(ExpectedBytes, ActualBytes)
}

func errorf(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	format :=
		`
		expected: %v
        actual  : %v
        message : %v`

	msg := ""

	switch len(msgAndArgs) {
	case 0:
		break
	case 1:
		msg = msgAndArgs[0].(string)

	default:
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	t.Errorf(format, expected, actual, msg)
}

func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if equal(expected, actual) {
		return true
	}

	errorf(t, expected, actual, msgAndArgs...)
	return false
}

func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !equal(expected, actual) {
		return true
	}

	errorf(t, expected, actual, msgAndArgs...)
	return false
}

func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if equal(expected, actual) {
		return
	}

	errorf(t, expected, actual, msgAndArgs...)
	t.FailNow()
}

func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !equal(expected, actual) {
		return
	}

	errorf(t, expected, actual, msgAndArgs...)
	t.FailNow()
}
