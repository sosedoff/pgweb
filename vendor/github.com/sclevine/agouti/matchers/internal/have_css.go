package internal

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti/matchers/internal/colorparser"
)

type HaveCSSMatcher struct {
	ExpectedProperty string
	ExpectedValue    string
	actualValue      string

	isColorComparison bool
	expectedColor     colorparser.Color
	actualColor       colorparser.Color
}

func (m *HaveCSSMatcher) Match(actual interface{}) (success bool, err error) {
	actualSelection, ok := actual.(interface {
		CSS(property string) (string, error)
	})

	if !ok {
		return false, fmt.Errorf("HaveCSS matcher requires a *Selection.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualValue, err = actualSelection.CSS(m.ExpectedProperty)
	if err != nil {
		return false, err
	}

	expectedColor, err := colorparser.ParseCSSColor(m.ExpectedValue)
	if err != nil {
		return m.actualValue == m.ExpectedValue, nil
	}

	actualColor, err := colorparser.ParseCSSColor(m.actualValue)
	if err != nil {
		return false, errors.New(expectedColorMessage(m.ExpectedValue, expectedColor, m.actualValue))
	}

	m.isColorComparison = true
	m.expectedColor = expectedColor
	m.actualColor = actualColor
	return reflect.DeepEqual(actualColor, expectedColor), nil
}

func (m *HaveCSSMatcher) FailureMessage(actual interface{}) (message string) {
	var expectedValue, actualValue string
	if m.isColorComparison {
		expectedValue, actualValue = m.styleColor(m.expectedColor), m.styleColor(m.actualColor)
	} else {
		expectedValue, actualValue = m.style(m.ExpectedValue), m.style(m.actualValue)
	}
	return valueMessage(actual, "to have CSS matching", expectedValue, actualValue)
}

func (m *HaveCSSMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	var expectedValue, actualValue string
	if m.isColorComparison {
		expectedValue, actualValue = m.styleColor(m.expectedColor), m.styleColor(m.actualColor)
	} else {
		expectedValue, actualValue = m.style(m.ExpectedValue), m.style(m.actualValue)
	}
	return valueMessage(actual, "not to have CSS matching", expectedValue, actualValue)
}

func (m *HaveCSSMatcher) style(value string) string {
	return fmt.Sprintf(`%s: "%s"`, m.ExpectedProperty, value)
}

func (m *HaveCSSMatcher) styleColor(value colorparser.Color) string {
	return fmt.Sprintf(`%s: %s`, m.ExpectedProperty, value)
}
