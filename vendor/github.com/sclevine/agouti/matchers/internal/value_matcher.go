package internal

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
)

type ValueMatcher struct {
	Method      string
	Property    string
	Expected    interface{}
	actualValue interface{}
}

func (m *ValueMatcher) Match(actual interface{}) (success bool, err error) {
	method := reflect.ValueOf(actual).MethodByName(m.Method)
	if !method.IsValid() {
		return false, fmt.Errorf("Have%s matcher requires a *Selection.  Got:\n%s", m.Method, format.Object(actual, 1))
	}

	results := method.Call(nil)
	propertyValue, errValue := results[0], results[1]
	if !errValue.IsNil() {
		return false, errValue.Interface().(error)
	}

	m.actualValue = propertyValue.Interface()

	return reflect.DeepEqual(m.actualValue, m.Expected), nil
}

func (m *ValueMatcher) FailureMessage(actual interface{}) (message string) {
	return valueMessage(actual, fmt.Sprintf("to have %s equaling", m.Property), m.Expected, m.actualValue)
}

func (m *ValueMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return valueMessage(actual, fmt.Sprintf("not to have %s equaling", m.Property), m.Expected, m.actualValue)
}
