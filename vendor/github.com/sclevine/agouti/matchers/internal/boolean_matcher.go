package internal

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
)

type BooleanMatcher struct {
	Method   string
	Property string
}

func (m *BooleanMatcher) Match(actual interface{}) (success bool, err error) {
	method := reflect.ValueOf(actual).MethodByName(m.Method)
	if !method.IsValid() {
		return false, fmt.Errorf("Be%s matcher requires a *Selection.  Got:\n%s", m.Method, format.Object(actual, 1))
	}

	results := method.Call(nil)
	propertyValue, errValue := results[0], results[1]
	if !errValue.IsNil() {
		return false, errValue.Interface().(error)
	}

	return propertyValue.Bool(), nil
}

func (m *BooleanMatcher) FailureMessage(actual interface{}) (message string) {
	return booleanMessage(actual, "to be "+m.Property)
}

func (m *BooleanMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return booleanMessage(actual, "not to be "+m.Property)
}
