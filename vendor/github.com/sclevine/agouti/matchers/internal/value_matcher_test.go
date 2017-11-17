package internal_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers/internal"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("ValueMatcher", func() {
	var (
		matcher   *ValueMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.StringCall.ReturnString = "selection 'CSS: #selector'"
		matcher = &ValueMatcher{Method: "Text", Property: "text", Expected: "some text"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			Context("when the expected text matches the actual text", func() {
				It("should successfully return true", func() {
					selection.TextCall.ReturnText = "some text"
					Expect(matcher.Match(selection)).To(BeTrue())
				})
			})

			Context("when the expected text does not equal the actual text", func() {
				It("should successfully return false", func() {
					selection.TextCall.ReturnText = "some other text"
					Expect(matcher.Match(selection)).To(BeFalse())
				})
			})

			Context("when retrieving the text fails", func() {
				It("should return an error", func() {
					selection.TextCall.Err = errors.New("some error")
					_, err := matcher.Match(selection)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a selection", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a selection")
				Expect(err).To(MatchError("HaveText matcher requires a *Selection.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message with the provided property name", func() {
			selection.TextCall.ReturnText = "some other text"
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' to have text equaling\n    some text"))
			Expect(message).To(ContainSubstring("but found\n    some other text"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message with the provided property name", func() {
			selection.TextCall.ReturnText = "some text"
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' not to have text equaling\n    some text"))
			Expect(message).To(ContainSubstring("but found\n    some text"))
		})
	})
})
