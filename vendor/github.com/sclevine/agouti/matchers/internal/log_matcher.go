package internal

import (
	"fmt"
	"strings"

	"github.com/onsi/gomega/format"
	"github.com/sclevine/agouti"
)

type LogMatcher struct {
	ExpectedMessages []string
	Levels           []string
	Name             string
	Type             string
}

func (m *LogMatcher) Match(actual interface{}) (success bool, err error) {
	actualPage, ok := actual.(interface {
		ReadAllLogs(logType string) ([]agouti.Log, error)
	})

	if !ok {
		return false, fmt.Errorf("HaveLogged%s matcher requires a Page.  Got:\n%s", firstToUpper(m.Name), format.Object(actual, 1))
	}

	logs, err := actualPage.ReadAllLogs(m.Type)
	if err != nil {
		return false, err
	}

	if len(m.ExpectedMessages) == 0 {
		return m.anyLogsForLevels(logs), nil
	}

	for _, message := range m.ExpectedMessages {
		if !m.messageInLogs(message, logs) {
			return false, nil
		}
	}

	return true, nil
}

func (m *LogMatcher) FailureMessage(actual interface{}) (message string) {
	if len(m.ExpectedMessages) == 0 {
		return booleanMessage(actual, fmt.Sprintf("to have logged %s logs", m.Name))

	}
	messages := strings.Join(m.ExpectedMessages, "\n"+tab)
	return equalityMessage(actual, fmt.Sprintf("to have %s logs matching", m.Name), messages)
}

func (m *LogMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	if len(m.ExpectedMessages) == 0 {
		return booleanMessage(actual, fmt.Sprintf("not to have logged %s logs", m.Name))

	}
	messages := strings.Join(m.ExpectedMessages, "\n"+tab)
	return equalityMessage(actual, fmt.Sprintf("not to have %s logs matching", m.Name), messages)
}

func (m *LogMatcher) messageInLogs(message string, logs []agouti.Log) bool {
	for _, log := range logs {
		if m.logInLevels(log) && log.Message == message {
			return true
		}
	}
	return false
}

func (m *LogMatcher) anyLogsForLevels(logs []agouti.Log) bool {
	for _, log := range logs {
		if m.logInLevels(log) {
			return true
		}
	}
	return false
}

func (m *LogMatcher) logInLevels(log agouti.Log) bool {
	for _, level := range m.Levels {
		if level == log.Level {
			return true
		}
	}
	return false
}

func firstToUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
