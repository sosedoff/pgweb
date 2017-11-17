package internal_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers/internal"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("HaveCSS", func() {
	var (
		matcher   *HaveCSSMatcher
		selection *mocks.Selection
	)

	BeforeEach(func() {
		selection = &mocks.Selection{}
		selection.StringCall.ReturnString = "selection 'CSS: #selector'"
		matcher = &HaveCSSMatcher{ExpectedProperty: "some-property", ExpectedValue: "some value"}
	})

	Describe("#Match", func() {
		Context("when the actual object is a selection", func() {
			It("should request the provided page property", func() {
				matcher.Match(selection)
				Expect(selection.CSSCall.Property).To(Equal("some-property"))
			})

			Context("when the expected property value matches the actual property value", func() {
				It("should successfully return true", func() {
					selection.CSSCall.ReturnValue = "some value"
					Expect(matcher.Match(selection)).To(BeTrue())
				})
			})

			Context("when the expected property value does not match the actual property value", func() {
				It("should successfully return false", func() {
					selection.CSSCall.ReturnValue = "some other value"
					Expect(matcher.Match(selection)).To(BeFalse())

				})
			})

			Context("when retrieving the CSS property value fails", func() {
				It("should return an error", func() {
					selection.CSSCall.Err = errors.New("some error")
					_, err := matcher.Match(selection)
					Expect(err).To(MatchError("some error"))
				})
			})

			Context("when the expected value is a color", func() {
				BeforeEach(func() {
					matcher = &HaveCSSMatcher{ExpectedProperty: "some-property", ExpectedValue: "blue"}
				})

				Context("when the actual value is a matching color", func() {
					BeforeEach(func() {
						selection.CSSCall.ReturnValue = "rgba(0,0,255,1.0)"
					})

					It("should succeed", func() {
						Expect(matcher.Match(selection)).To(BeTrue())
					})

					Describe("#NegatedFailureMessage", func() {
						It("should return a negated failure message", func() {
							matcher.Match(selection)
							message := matcher.NegatedFailureMessage(selection)
							Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' not to have CSS matching\n    some-property: Color{R:0, G:0, B:255, A:1.00}"))
							Expect(message).To(ContainSubstring("but found\n    some-property: Color{R:0, G:0, B:255, A:1.00}"))
						})
					})
				})

				Context("when the actual value is a non-matching color", func() {
					BeforeEach(func() {
						selection.CSSCall.ReturnValue = "rgba(255,0,0,1.0)"
					})

					It("should fail", func() {
						Expect(matcher.Match(selection)).To(BeFalse())
					})

					Describe("#FailureMessage", func() {
						It("should return a failure message", func() {
							matcher.Match(selection)
							message := matcher.FailureMessage(selection)
							Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' to have CSS matching\n    some-property: Color{R:0, G:0, B:255, A:1.00}"))
							Expect(message).To(ContainSubstring("but found\n    some-property: Color{R:255, G:0, B:0, A:1.00}"))
						})
					})
				})

				Context("when the actual value is not a color", func() {
					BeforeEach(func() {
						selection.CSSCall.ReturnValue = "not-a-color"
					})

					It("should error", func() {
						_, err := matcher.Match(selection)
						Expect(err).To(MatchError("The expected value:\n    blue\nis a color:\n    Color{R:0, G:0, B:255, A:1.00}\nBut the actual value:\n    not-a-color\nis not.\n"))
					})
				})
			})
		})

		Context("when the actual object is not a selection", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a selection")
				Expect(err).To(MatchError("HaveCSS matcher requires a *Selection.  Got:\n    <string>: not a selection"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message", func() {
			selection.CSSCall.ReturnValue = "some other value"
			matcher.Match(selection)
			message := matcher.FailureMessage(selection)
			Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' to have CSS matching\n    some-property: \"some value\""))
			Expect(message).To(ContainSubstring("but found\n    some-property: \"some other value\""))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			selection.CSSCall.ReturnValue = "some value"
			matcher.Match(selection)
			message := matcher.NegatedFailureMessage(selection)
			Expect(message).To(ContainSubstring("Expected selection 'CSS: #selector' not to have CSS matching\n    some-property: \"some value\""))
			Expect(message).To(ContainSubstring("but found\n    some-property: \"some value\""))
		})
	})
})
