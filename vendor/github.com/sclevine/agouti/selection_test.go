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

var _ = Describe("Selection", func() {
	var (
		firstElement  *mocks.Element
		secondElement *api.Element
	)

	BeforeEach(func() {
		firstElement = &mocks.Element{}
		secondElement = &api.Element{}
	})

	Describe("#String", func() {
		It("should return a string representation of the selection", func() {
			selection := NewTestMultiSelection(nil, nil, "#selector")
			Expect(selection.AllByXPath("#subselector").String()).To(Equal("selection 'CSS: #selector | XPath: #subselector'"))
		})
	})

	Describe("#Elements", func() {
		var (
			selection         *Selection
			elementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			elementRepository = &mocks.ElementRepository{}
			selection = NewTestSelection(nil, elementRepository, "#selector")
		})

		It("should return a []*api.Elements retrieved from the element repository", func() {
			elements := []*api.Element{{ID: "first"}, {ID: "second"}}
			elementRepository.GetCall.ReturnElements = []element.Element{elements[0], elements[1]}
			Expect(selection.Elements()).To(Equal(elements))
		})

		Context("when retrieving the elements fails", func() {
			It("should return an error", func() {
				elementRepository.GetCall.Err = errors.New("some error")
				_, err := selection.Elements()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Count", func() {
		var (
			selection         *MultiSelection
			elementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			elementRepository = &mocks.ElementRepository{}
			selection = NewTestMultiSelection(nil, elementRepository, "#selector")
			elementRepository.GetCall.ReturnElements = []element.Element{firstElement, secondElement}
		})

		It("should successfully return the number of elements", func() {
			Expect(selection.Count()).To(Equal(2))
		})

		Context("when the the session fails to retrieve the elements", func() {
			It("should return an error", func() {
				elementRepository.GetCall.Err = errors.New("some error")
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#EqualsElement", func() {
		var (
			firstSelection          *Selection
			secondSelection         *Selection
			firstElementRepository  *mocks.ElementRepository
			secondElementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			firstElementRepository = &mocks.ElementRepository{}
			firstElementRepository.GetExactlyOneCall.ReturnElement = firstElement
			firstSelection = NewTestSelection(nil, firstElementRepository, "#first_selector")

			secondElementRepository = &mocks.ElementRepository{}
			secondElementRepository.GetExactlyOneCall.ReturnElement = secondElement
			secondSelection = NewTestSelection(nil, secondElementRepository, "#second_selector")
		})

		It("should compare the selection elements for equality", func() {
			firstSelection.EqualsElement(secondSelection)
			Expect(firstElement.IsEqualToCall.Element).To(ExactlyEqual(secondElement))
		})

		It("should successfully return true if they are equal", func() {
			firstElement.IsEqualToCall.ReturnEquals = true
			Expect(firstSelection.EqualsElement(secondSelection)).To(BeTrue())
		})

		It("should successfully return false if they are not equal", func() {
			firstElement.IsEqualToCall.ReturnEquals = false
			Expect(firstSelection.EqualsElement(secondSelection)).To(BeFalse())
		})

		Context("when the provided object is a *MultiSelection", func() {
			It("should not fail", func() {
				multiSelection := NewTestMultiSelection(nil, secondElementRepository, "#multi_selector")
				Expect(firstSelection.EqualsElement(multiSelection)).To(BeFalse())
				Expect(firstElement.IsEqualToCall.Element).To(ExactlyEqual(secondElement))
			})
		})

		Context("when the provided object is not a type of selection", func() {
			It("should return an error", func() {
				_, err := firstSelection.EqualsElement("not a selection")
				Expect(err).To(MatchError("must be *Selection or *MultiSelection"))
			})
		})

		Context("when there is an error retrieving elements from the selection", func() {
			It("should return an error", func() {
				firstElementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #first_selector [single]': some error"))
			})
		})

		Context("when there is an error retrieving elements from the other selection", func() {
			It("should return an error", func() {
				secondElementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #second_selector [single]': some error"))
			})
		})

		Context("when the session fails to compare the elements", func() {
			It("should return an error", func() {
				firstElement.IsEqualToCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to compare selection 'CSS: #first_selector [single]' to selection 'CSS: #second_selector [single]': some error"))
			})
		})
	})

	Describe("#MouseToElement", func() {
		var (
			selection         *Selection
			session           *mocks.Session
			elementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			elementRepository = &mocks.ElementRepository{}
			elementRepository.GetExactlyOneCall.ReturnElement = secondElement
			session = &mocks.Session{}
			selection = NewTestSelection(session, elementRepository, "#selector")
		})

		It("should successfully instruct the session to move the mouse over the selection", func() {
			Expect(selection.MouseToElement()).To(Succeed())
			Expect(session.MoveToCall.Element).To(Equal(secondElement))
			Expect(session.MoveToCall.Offset).To(BeNil())
		})

		Context("when the element repository fails to return exactly one element", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				err := selection.MouseToElement()
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #selector [single]': some error"))
			})
		})

		Context("when the session fails to move the mouse to the element", func() {
			It("should return an error", func() {
				session.MoveToCall.Err = errors.New("some error")
				err := selection.MouseToElement()
				Expect(err).To(MatchError("failed to move mouse to element for selection 'CSS: #selector [single]': some error"))
			})
		})
	})
})
