package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("Selection Matchers", func() {
	var selection *mocks.Selection

	BeforeEach(func() {
		selection = &mocks.Selection{}
	})

	Describe("#HaveText", func() {
		It("should return a ValueMatcher with the 'Text' method", func() {
			selection.TextCall.ReturnText = "some text"
			Expect(selection).To(HaveText("some text"))
			Expect(selection).NotTo(HaveText("some other text"))
		})

		It("should set the matcher property to 'text'", func() {
			Expect(HaveText("").FailureMessage(nil)).To(ContainSubstring("to have text"))
		})
	})

	Describe("#MatchText", func() {
		It("should return a MatchText matcher", func() {
			selection.TextCall.ReturnText = "some text"
			Expect(selection).To(MatchText("s[^t]+text"))
			Expect(selection).NotTo(MatchText("so*text"))
		})
	})

	Describe("#HaveCount", func() {
		It("should return a ValueMatcher with the 'Count' method", func() {
			selection.CountCall.ReturnCount = 1
			Expect(selection).To(HaveCount(1))
			Expect(selection).NotTo(HaveCount(2))
		})

		It("should set the matcher property to 'element count'", func() {
			Expect(HaveCount(0).FailureMessage(nil)).To(ContainSubstring("to have element count"))
		})
	})

	Describe("#HaveAttribute", func() {
		It("should return a HaveAttribute matcher", func() {
			selection.AttributeCall.ReturnValue = "some value"
			Expect(selection).To(HaveAttribute("some-attribute", "some value"))
			Expect(selection).NotTo(HaveAttribute("some-attribute", "some other value"))
		})
	})

	Describe("#HaveCSS", func() {
		It("should return a HaveCSS matcher", func() {
			selection.CSSCall.ReturnValue = "some value"
			Expect(selection).To(HaveCSS("some-property", "some value"))
			Expect(selection).NotTo(HaveCSS("some-property", "some other value"))
		})
	})

	Describe("#BeSelected", func() {
		It("should return a BooleanMatcher with the 'Selected' method", func() {
			selection.SelectedCall.ReturnSelected = true
			Expect(selection).To(BeSelected())
			selection.SelectedCall.ReturnSelected = false
			Expect(selection).NotTo(BeSelected())
		})

		It("should set the matcher property to 'selected'", func() {
			Expect(BeSelected().FailureMessage(nil)).To(HaveSuffix("to be selected"))
		})
	})

	Describe("#BeVisible", func() {
		It("should return a BooleanMatcher with the 'Visible' method", func() {
			selection.VisibleCall.ReturnVisible = true
			Expect(selection).To(BeVisible())
			selection.VisibleCall.ReturnVisible = false
			Expect(selection).NotTo(BeVisible())
		})

		It("should set the matcher property to 'visible'", func() {
			Expect(BeVisible().FailureMessage(nil)).To(HaveSuffix("to be visible"))
		})
	})

	Describe("#BeEnabled", func() {
		It("should return a BooleanMatcher with the 'Enabled' method", func() {
			selection.EnabledCall.ReturnEnabled = true
			Expect(selection).To(BeEnabled())
			selection.EnabledCall.ReturnEnabled = false
			Expect(selection).NotTo(BeEnabled())
		})

		It("should set the matcher property to 'enabled'", func() {
			Expect(BeEnabled().FailureMessage(nil)).To(HaveSuffix("to be enabled"))
		})
	})

	Describe("#BeActive", func() {
		It("should return a BooleanMatcher with the 'Active' method", func() {
			selection.ActiveCall.ReturnActive = true
			Expect(selection).To(BeActive())
			selection.ActiveCall.ReturnActive = false
			Expect(selection).NotTo(BeActive())
		})

		It("should set the matcher property to 'active'", func() {
			Expect(BeActive().FailureMessage(nil)).To(HaveSuffix("to be active"))
		})
	})

	Describe("#BeFound", func() {
		It("should return a BeFound matcher", func() {
			selection.CountCall.ReturnCount = 1
			Expect(selection).To(BeFound())
			selection.CountCall.ReturnCount = 0
			Expect(selection).NotTo(BeFound())
		})
	})

	Describe("#EqualElement", func() {
		It("should return a EqualElement matcher", func() {
			selection.EqualsElementCall.ReturnEquals = true
			Expect(selection).To(EqualElement(selection))
			selection.EqualsElementCall.ReturnEquals = false
			Expect(selection).NotTo(EqualElement(selection))
		})
	})
})
