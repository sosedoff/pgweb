package internal

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

var tab = format.Indent

func valueMessage(actual interface{}, message string, expected, actualValue interface{}) string {
	failureMessage := "Expected %s %s\n%s%s\nbut found\n%s%s"
	return fmt.Sprintf(failureMessage, actual, message, tab, expected, tab, actualValue)
}

func booleanMessage(actual interface{}, message string) string {
	failureMessage := "Expected %s %s"
	return fmt.Sprintf(failureMessage, actual, message)
}

func equalityMessage(actual interface{}, message string, expected interface{}) string {
	failureMessage := "Expected %s %s\n%s%s"
	return fmt.Sprintf(failureMessage, actual, message, tab, expected)
}

func expectedColorMessage(expectedValue string, expectedColor, actualColor interface{}) string {
	failureMessage := "The expected value:\n%s%s\nis a color:\n%s%s\nBut the actual value:\n%s%s\nis not.\n"
	return fmt.Sprintf(failureMessage, tab, expectedValue, tab, expectedColor, tab, actualColor)
}
