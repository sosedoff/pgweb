package internal_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers/internal"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("BooleanMatcher", func() {
	var (
		matcher   *BooleanMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.StringCall.ReturnString = "selection 'CSS: #selector'"
		matcher = &BooleanMatcher{Method: "Visible", Property: "visible"}
	})

	Describe("#Match", func() {
		Context("when the actual object has a corresponding method", func() {
			Context("when the provided method returns true", func() {
				It("should successfully return true", func() {
					selection.VisibleCall.ReturnVisible = true
					Expect(matcher.Match(selection)).To(BeTrue())
				})
			})

			Context("when the provided method returns false", func() {
				It("should successfully return false", func() {
					selection.VisibleCall.ReturnVisible = false
					Expect(matcher.Match(selection)).To(BeFalse())
				})
			})

			Context("when the provided method returns an error", func() {
				It("should return an error", func() {
					selection.VisibleCall.Err = errors.New("some error")
					_, err := matcher.Match(selection)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object does not have the corresponding method", func() {
			It("should return an error", func() {
				_, err := matcher.Match("missing method")
				Expect(err).To(MatchError("BeVisible matcher requires a *Selection.  Got:\n    <string>: missing method"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message with the provided property name", func() {
			selection.VisibleCall.ReturnVisible = false
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' to be visible"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message with the provided property name", func() {
			selection.VisibleCall.ReturnVisible = true
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' not to be visible"))
		})
	})
})
