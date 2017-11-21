package matchers

import (
	"github.com/onsi/gomega/types"
	"github.com/sclevine/agouti/matchers/internal"
)

// HaveTitle passes when the expected title is equivalent to the
// title of the provided page.
func HaveTitle(title string) types.GomegaMatcher {
	return &internal.ValueMatcher{Method: "Title", Property: "title", Expected: title}
}

// HaveURL passes when the expected URL is equivalent to the
// current URL of the provided page.
func HaveURL(url string) types.GomegaMatcher {
	return &internal.ValueMatcher{Method: "URL", Property: "URL", Expected: url}
}

// HavePopupText passes when the expected text is equivalent to the
// text contents of an open alert, confirm, or prompt popup.
func HavePopupText(text string) types.GomegaMatcher {
	return &internal.ValueMatcher{Method: "PopupText", Property: "popup text", Expected: text}
}

// HaveWindowCount passes when the expected window count is equivalent
// to the number of open windows.
func HaveWindowCount(count int) types.GomegaMatcher {
	return &internal.ValueMatcher{Method: "WindowCount", Property: "window count", Expected: count}
}

// HaveLoggedError passes when all of the expected log messages are logged as
// errors in the browser console. If no message is provided, this matcher will
// pass if any error message has been logged. When negated, this matcher will
// only fail if all of the provided messages are logged.
func HaveLoggedError(messages ...string) types.GomegaMatcher {
	return &internal.LogMatcher{
		ExpectedMessages: messages,
		Levels:           []string{"WARNING", "SEVERE"},
		Name:             "error",
		Type:             "browser",
	}
}

// HaveLoggedInfo passes when all of the expected log messages are logged in
// the browser console. If no messages are provided, this matcher will pass if
// any message has been logged. When negated, this matcher will only fail if
// all of the provided messages are logged. Error logs are not considered in
// any of these cases.
func HaveLoggedInfo(messages ...string) types.GomegaMatcher {
	return &internal.LogMatcher{
		ExpectedMessages: messages,
		Levels:           []string{"INFO"},
		Name:             "info",
		Type:             "browser",
	}
}
