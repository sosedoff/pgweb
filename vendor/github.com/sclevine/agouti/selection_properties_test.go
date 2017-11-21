package agouti_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	. "github.com/sclevine/agouti/internal/matchers"
	"github.com/sclevine/agouti/internal/mocks"
)

var _ = Describe("Selection Properties", func() {
	var (
		selection         *MultiSelection
		session           *mocks.Session
		elementRepository *mocks.ElementRepository
		firstElement      *mocks.Element
		secondElement     *mocks.Element
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		firstElement = &mocks.Element{}
		secondElement = &mocks.Element{}
		elementRepository = &mocks.ElementRepository{}
		selection = NewTestMultiSelection(session, elementRepository, "#selector")
	})

	Describe("#Text", func() {
		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = firstElement
		})

		It("should successfully return the text", func() {
			firstElement.GetTextCall.ReturnText = "some text"
			Expect(selection.Text()).To(Equal("some text"))
		})

		Context("when the element repository fails to return exactly one element", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := selection.Text()
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the session fails to retrieve the element text", func() {
			It("should return an error", func() {
				firstElement.GetTextCall.Err = errors.New("some error")
				_, err := selection.Text()
				Expect(err).To(MatchError("failed to retrieve text for selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#Active", func() {
		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = firstElement
		})

		It("should successfully compare the active and selected elements", func() {
			activeElement := &api.Element{}
			session.GetActiveElementCall.ReturnElement = activeElement
			_, err := selection.Active()
			Expect(err).NotTo(HaveOccurred())
			Expect(firstElement.IsEqualToCall.Element).To(ExactlyEqual(activeElement))
		})

		Context("when the active element equals the selected element", func() {
			It("should successfully return true", func() {
				firstElement.IsEqualToCall.ReturnEquals = true
				Expect(selection.Active()).To(BeTrue())
			})
		})

		Context("when the active element does not equal the selected element", func() {
			It("should successfully return false", func() {
				firstElement.IsEqualToCall.ReturnEquals = false
				Expect(selection.Active()).To(BeFalse())
			})
		})

		Context("when the element repository fails to return exactly one element", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := selection.Active()
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the session fails to retrieve the active element", func() {
			It("should return an error", func() {
				session.GetActiveElementCall.Err = errors.New("some error")
				_, err := selection.Active()
				Expect(err).To(MatchError("failed to retrieve active element: some error"))
			})
		})

		Context("when the session fails to compare active element to the selected element", func() {
			It("should return an error", func() {
				firstElement.IsEqualToCall.Err = errors.New("some error")
				_, err := selection.Active()
				Expect(err).To(MatchError("failed to compare selection to active element: some error"))
			})
		})
	})

	Describe("#Attribute", func() {
		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = firstElement
		})

		It("should request the attribute value using the attribute name", func() {
			_, err := selection.Attribute("some-attribute")
			Expect(err).NotTo(HaveOccurred())
			Expect(firstElement.GetAttributeCall.Attribute).To(Equal("some-attribute"))
		})

		It("should successfully return the attribute value", func() {
			firstElement.GetAttributeCall.ReturnValue = "some value"
			Expect(selection.Attribute("some-attribute")).To(Equal("some value"))
		})

		Context("when the element repository fails to return exactly one element", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := selection.Attribute("some-attribute")
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the session fails to retrieve the requested element attribute", func() {
			It("should return an error", func() {
				firstElement.GetAttributeCall.Err = errors.New("some error")
				_, err := selection.Attribute("some-attribute")
				Expect(err).To(MatchError("failed to retrieve attribute value for selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#CSS", func() {
		BeforeEach(func() {
			elementRepository.GetExactlyOneCall.ReturnElement = firstElement
		})

		It("should request the CSS property value using the property name", func() {
			_, err := selection.CSS("some-property")
			Expect(err).NotTo(HaveOccurred())
			Expect(firstElement.GetCSSCall.Property).To(Equal("some-property"))
		})

		It("should successfully return the property value", func() {
			firstElement.GetCSSCall.ReturnValue = "some value"
			Expect(selection.CSS("some-property")).To(Equal("some value"))
		})

		Context("when the element repository fails to return exactly one element", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := selection.CSS("some-property")
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the the session fails to retrieve the requested element CSS property", func() {
			It("should return an error", func() {
				firstElement.GetCSSCall.Err = errors.New("some error")
				_, err := selection.CSS("some-property")
				Expect(err).To(MatchError("failed to retrieve CSS property value for selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#Selected", func() {
		BeforeEach(func() {
			elementRepository.GetAtLeastOneCall.ReturnElements = []element.Element{firstElement, secondElement}
		})

		It("should return true when all elements are selected", func() {
			firstElement.IsSelectedCall.ReturnSelected = true
			secondElement.IsSelectedCall.ReturnSelected = true
			Expect(selection.Selected()).To(BeTrue())
		})

		It("should return false when any elements are not selected", func() {
			firstElement.IsSelectedCall.ReturnSelected = true
			secondElement.IsSelectedCall.ReturnSelected = false
			Expect(selection.Selected()).To(BeFalse())
		})

		Context("when the element repository fails to return at least one element", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				_, err := selection.Selected()
				Expect(err).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the the session fails to retrieve any elements' selected status", func() {
			It("should return an error", func() {
				firstElement.IsSelectedCall.ReturnSelected = true
				secondElement.IsSelectedCall.Err = errors.New("some error")
				_, err := selection.Selected()
				Expect(err).To(MatchError("failed to determine whether selection 'CSS: #selector' is selected: some error"))
			})
		})
	})

	Describe("#Visible", func() {
		BeforeEach(func() {
			elementRepository.GetAtLeastOneCall.ReturnElements = []element.Element{firstElement, secondElement}
		})

		It("should return true when all elements are visible", func() {
			firstElement.IsDisplayedCall.ReturnDisplayed = true
			secondElement.IsDisplayedCall.ReturnDisplayed = true
			Expect(selection.Visible()).To(BeTrue())
		})

		It("should return false when any elements are not visible", func() {
			firstElement.IsDisplayedCall.ReturnDisplayed = true
			secondElement.IsDisplayedCall.ReturnDisplayed = false
			Expect(selection.Visible()).To(BeFalse())
		})

		Context("when the element repository fails to return at least one element", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				_, err := selection.Visible()
				Expect(err).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the the session fails to retrieve any elements' visible status", func() {
			It("should return an error", func() {
				firstElement.IsDisplayedCall.ReturnDisplayed = true
				secondElement.IsDisplayedCall.Err = errors.New("some error")
				_, err := selection.Visible()
				Expect(err).To(MatchError("failed to determine whether selection 'CSS: #selector' is visible: some error"))
			})
		})
	})

	Describe("#Enabled", func() {
		BeforeEach(func() {
			elementRepository.GetAtLeastOneCall.ReturnElements = []element.Element{firstElement, secondElement}
		})

		It("should return true when all elements are enabled", func() {
			firstElement.IsEnabledCall.ReturnEnabled = true
			secondElement.IsEnabledCall.ReturnEnabled = true
			Expect(selection.Enabled()).To(BeTrue())
		})

		It("should return false when any elements are not enabled", func() {
			firstElement.IsEnabledCall.ReturnEnabled = true
			secondElement.IsEnabledCall.ReturnEnabled = false
			Expect(selection.Enabled()).To(BeFalse())
		})

		Context("when the element repository fails to return at least one element", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				_, err := selection.Enabled()
				Expect(err).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the the session fails to retrieve any element's enabled status", func() {
			It("should return an error", func() {
				firstElement.IsEnabledCall.ReturnEnabled = true
				secondElement.IsEnabledCall.Err = errors.New("some error")
				_, err := selection.Enabled()
				Expect(err).To(MatchError("failed to determine whether selection 'CSS: #selector' is enabled: some error"))
			})
		})
	})
})
