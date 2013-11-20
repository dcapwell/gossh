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

func isInteger(a interface{}) bool {
	kind := reflect.TypeOf(a).Kind()
	return reflect.Int <= kind && kind <= reflect.Int64
}

func toInteger(a interface{}) int64 {
	if isInteger(a) {
		return reflect.ValueOf(a).Int()
	} else if isUnsignedInteger(a) {
		return int64(reflect.ValueOf(a).Uint())
	} else if isFloat(a) {
		return int64(reflect.ValueOf(a).Float())
	} else {
		panic(fmt.Sprintf("Expected a number!  Got <%T> %#v", a, a))
	}
}

func isUnsignedInteger(a interface{}) bool {
	kind := reflect.TypeOf(a).Kind()
	return reflect.Uint <= kind && kind <= reflect.Uint64
}

func isFloat(a interface{}) bool {
	kind := reflect.TypeOf(a).Kind()
	return reflect.Float32 <= kind && kind <= reflect.Float64
}

func BeGreaterThan(num int64) OmegaMatcher {
	return &GreaterThanMatcher{
		Num: num,
	}
}

type GreaterThanMatcher struct {
	Num int64
}

func (matcher *GreaterThanMatcher) Match(actual interface{}) (success bool, message string, err error) {
	//TODO support floating point as well
	ok := isInteger(actual)
	if ok {
		num := toInteger(actual)
		if num > matcher.Num {
			return true, fmt.Sprintf("Expected\n (%d) not to be greater than %d", num, matcher.Num), nil
		} else {
			return false, fmt.Sprintf("Expected\n (%d) to be greater than %d", num, matcher.Num), nil
		}
	} else {
		return false, "", fmt.Errorf("GreaterThan matcher expects a int based type.  Got:%s", formatObject(actual))
	}
}
