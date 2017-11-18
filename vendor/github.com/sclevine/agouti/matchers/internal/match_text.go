package internal

import (
	"fmt"
	"regexp"

	"github.com/onsi/gomega/format"
)

type MatchTextMatcher struct {
	Regexp     string
	actualText string
}

func (m *MatchTextMatcher) Match(actual interface{}) (success bool, err error) {
	actualSelection, ok := actual.(interface {
		Text() (string, error)
	})

	if !ok {
		return false, fmt.Errorf("MatchText matcher requires a *Selection.  Got:\n%s", format.Object(actual, 1))
	}

	m.actualText, err = actualSelection.Text()
	if err != nil {
		return false, err
	}

	return regexp.MatchString(m.Regexp, m.actualText)
}

func (m *MatchTextMatcher) FailureMessage(actual interface{}) (message string) {
	return valueMessage(actual, "to have text matching", m.Regexp, m.actualText)
}

func (m *MatchTextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return valueMessage(actual, "not to have text matching", m.Regexp, m.actualText)
}
