package gossh_test

import (
	"fmt"
	. "github.com/onsi/gomega"
	"reflect"
)

// pulled out of matchers code
func lengthOf(a interface{}) (int, bool) {
	if a == nil {
		return 0, false
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Map, reflect.Array, reflect.String, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(a).Len(), true
	case reflect.Int:
		data, ok := a.(int)
		if ok {
			return data, true
		}
		return 0, false
	default:
		return 0, false
	}
}
func formatMessage(actual interface{}, message string, expected ...interface{}) string {
	if len(expected) == 0 {
		return fmt.Sprintf("Expected%s\n%s", formatObject(actual), message)
	} else {
		return fmt.Sprintf("Expected%s\n%s%s", formatObject(actual), message, formatObject(expected[0]))
	}
}

func formatObject(object interface{}) string {
	hasLength := false
	length := 0
	if !isString(object) {
		length, hasLength = lengthOf(object)
	}

	if hasLength {
		return fmt.Sprintf("\n\t<%T> of length %d: %#v", object, length, object)
	} else {
		return fmt.Sprintf("\n\t<%T>: %#v", object, object)
	}
}

func isString(a interface{}) bool {
	if a == nil {
		return false
	}
	return reflect.TypeOf(a).Kind() == reflect.String
}

func BeGreaterThan(num int) OmegaMatcher {
	return &GreaterThanMatcher{
		Num: num,
	}
}

type GreaterThanMatcher struct {
	Num int
}

func (matcher *GreaterThanMatcher) Match(actual interface{}) (success bool, message string, err error) {
	length, ok := lengthOf(actual)
	if ok {
		if length > matcher.Num {
			return true, fmt.Sprintf("Expected%s\n (length: %d) not to be greater than length %d", formatObject(actual), length, matcher.Num), nil
		} else {
			return false, fmt.Sprintf("Expected%s\n (length: %d) to be greater than length %d", formatObject(actual), length, matcher.Num), nil
		}
	} else {
		return false, "", fmt.Errorf("GreaterThan matcher expects a int/string/array/map/channel/slice.  Got:%s", formatObject(actual))
	}
}
