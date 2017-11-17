package internal_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers/internal"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("BeFoundMatcher", func() {
	var (
		matcher   *BeFoundMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.StringCall.ReturnString = "selection 'CSS: #selector'"
		matcher = &BeFoundMatcher{}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			Context("when the element is found", func() {
				It("should successfully return true", func() {
					selection.CountCall.ReturnCount = 1
					Expect(matcher.Match(selection)).To(BeTrue())
				})
			})

			Context("when the element is not found", func() {
				It("should successfully return false", func() {
					selection.CountCall.ReturnCount = 0
					Expect(matcher.Match(selection)).To(BeFalse())
				})
			})
		})

		Context("when the actual object is not a selection", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a selection")
				Expect(err).To(MatchError("BeFound matcher requires a *Selection.  Got:\n    <string>: not a selection"))
			})
		})

		Context("when there is an error retrieving the count", func() {
			Context("when the error is an 'element not found' error", func() {
				It("should successfully return false", func() {
					selection.CountCall.Err = errors.New("some error: element not found")
					Expect(matcher.Match(selection)).To(BeFalse())
				})
			})

			Context("when the error is an 'element index out of range' error", func() {
				It("should successfully return false", func() {
					selection.CountCall.Err = errors.New("some error: element index out of range")
					Expect(matcher.Match(selection)).To(BeFalse())
				})
			})

			Context("when the error is any other error", func() {
				It("should return an error", func() {
					selection.CountCall.Err = errors.New("some error")
					_, err := matcher.Match(selection)
					Expect(err).To(MatchError("some error"))
				})
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			selection.CountCall.ReturnCount = 0
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' to be found"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			selection.CountCall.ReturnCount = 1
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(Equal("Expected selection 'CSS: #selector' not to be found"))
		})
	})
})
