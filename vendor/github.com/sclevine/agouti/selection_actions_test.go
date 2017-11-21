package agouti_test

import (
	"errors"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	. "github.com/sclevine/agouti/internal/matchers"
	"github.com/sclevine/agouti/internal/mocks"
)

var _ = Describe("Selection Actions", func() {
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
		elementRepository.GetAtLeastOneCall.ReturnElements = []element.Element{firstElement, secondElement}
	})

	Describe("#Click", func() {
		It("should successfully click on all selected elements", func() {
			Expect(selection.Click()).To(Succeed())
			Expect(firstElement.ClickCall.Called).To(BeTrue())
			Expect(secondElement.ClickCall.Called).To(BeTrue())
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Click()).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when any click fails", func() {
			It("should return an error", func() {
				secondElement.ClickCall.Err = errors.New("some error")
				Expect(selection.Click()).To(MatchError("failed to click on selection 'CSS: #selector': some error"))
			})
		})
	})

	// TODO: extend mock to test multiple calls
	Describe("#DoubleClick", func() {
		var apiElement *api.Element

		BeforeEach(func() {
			apiElement = &api.Element{}
			elementRepository.GetAtLeastOneCall.ReturnElements = []element.Element{&api.Element{}, apiElement}
		})

		It("should successfully move the mouse to the middle of each selected element", func() {
			Expect(selection.DoubleClick()).To(Succeed())
			Expect(session.MoveToCall.Element).To(ExactlyEqual(apiElement))
			Expect(session.MoveToCall.Offset).To(BeNil())
		})

		It("should successfully double-click on each element", func() {
			Expect(selection.DoubleClick()).To(Succeed())
			Expect(session.DoubleClickCall.Called).To(BeTrue())
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when moving over any element fails", func() {
			It("should retun an error", func() {
				session.MoveToCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to move mouse to selection 'CSS: #selector': some error"))
			})
		})

		Context("when the double-clicking any element fails", func() {
			It("should return an error", func() {
				session.DoubleClickCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to double-click on selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#Fill", func() {
		It("should successfully clear each element", func() {
			Expect(selection.Fill("some text")).To(Succeed())
			Expect(firstElement.ClearCall.Called).To(BeTrue())
			Expect(secondElement.ClearCall.Called).To(BeTrue())
		})

		It("should successfully fill each element with the provided text", func() {
			Expect(selection.Fill("some text")).To(Succeed())
			Expect(firstElement.ValueCall.Text).To(Equal("some text"))
			Expect(secondElement.ValueCall.Text).To(Equal("some text"))
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Fill("some text")).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when clearing any element fails", func() {
			It("should return an error", func() {
				secondElement.ClearCall.Err = errors.New("some error")
				Expect(selection.Fill("some text")).To(MatchError("failed to clear selection 'CSS: #selector': some error"))
			})
		})

		Context("when entering text into any element fails", func() {
			It("should return an error", func() {
				secondElement.ValueCall.Err = errors.New("some error")
				Expect(selection.Fill("some text")).To(MatchError("failed to enter text into selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#Clear", func() {
		It("should successfully clear each element", func() {
			Expect(selection.Clear()).To(Succeed())
			Expect(firstElement.ClearCall.Called).To(BeTrue())
			Expect(secondElement.ClearCall.Called).To(BeTrue())
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Clear()).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when clearing any element fails", func() {
			It("should return an error", func() {
				secondElement.ClearCall.Err = errors.New("some error")
				Expect(selection.Clear()).To(MatchError("failed to clear selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#UploadFile", func() {
		BeforeEach(func() {
			firstElement.GetAttributeCall.ReturnValue = "file"
			firstElement.GetNameCall.ReturnName = "input"
			secondElement.GetAttributeCall.ReturnValue = "file"
			secondElement.GetNameCall.ReturnName = "input"
		})

		It("should successfully enter the absolute file path into each element", func() {
			Expect(selection.UploadFile("some-file")).To(Succeed())
			Expect(firstElement.ValueCall.Text).To(HaveSuffix(filepath.Join("agouti", "some-file")))
			Expect(secondElement.ValueCall.Text).To(HaveSuffix(filepath.Join("agouti", "some-file")))
		})

		It("should request the 'type' attribute for each element", func() {
			Expect(selection.UploadFile("some-file")).To(Succeed())
			Expect(firstElement.GetAttributeCall.Attribute).To(Equal("type"))
			Expect(secondElement.GetAttributeCall.Attribute).To(Equal("type"))
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.UploadFile("/some/file")).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when any element has a tag name other than 'input'", func() {
			It("should return an error", func() {
				secondElement.GetNameCall.ReturnName = "notinput"
				err := selection.UploadFile("some-file")
				Expect(err).To(MatchError("element for selection 'CSS: #selector' is not an input element"))
			})
		})

		Context("when the tag name of any element is not retrievable", func() {
			It("should return an error", func() {
				secondElement.GetNameCall.Err = errors.New("some error")
				err := selection.UploadFile("some-file")
				Expect(err).To(MatchError("failed to determine tag name of selection 'CSS: #selector': some error"))
			})
		})

		Context("when any element has a type attribute other than 'file'", func() {
			It("should return an error", func() {
				secondElement.GetAttributeCall.ReturnValue = "notfile"
				err := selection.UploadFile("some-file")
				Expect(err).To(MatchError("element for selection 'CSS: #selector' is not a file uploader"))
			})
		})

		Context("when the type attribute of any element is not retrievable", func() {
			It("should return an error", func() {
				secondElement.GetAttributeCall.Err = errors.New("some error")
				err := selection.UploadFile("some-file")
				Expect(err).To(MatchError("failed to determine type attribute of selection 'CSS: #selector': some error"))
			})
		})

		Context("when entering text into any element fails", func() {
			It("should return an error", func() {
				secondElement.ValueCall.Err = errors.New("some error")
				Expect(selection.UploadFile("/some/file")).To(MatchError("failed to enter text into selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#Check", func() {
		It("should successfully check the type of each checkbox", func() {
			firstElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.GetAttributeCall.ReturnValue = "checkbox"
			Expect(selection.Check()).To(Succeed())
			Expect(firstElement.GetAttributeCall.Attribute).To(Equal("type"))
			Expect(secondElement.GetAttributeCall.Attribute).To(Equal("type"))
		})

		Context("when all elements are checkboxes", func() {
			BeforeEach(func() {
				firstElement.GetAttributeCall.ReturnValue = "checkbox"
				secondElement.GetAttributeCall.ReturnValue = "checkbox"
			})

			It("should not click on the checked checkbox successfully", func() {
				firstElement.IsSelectedCall.ReturnSelected = true
				Expect(selection.Check()).To(Succeed())
				Expect(firstElement.ClickCall.Called).To(BeFalse())
			})

			It("should click on the unchecked checkboxes successfully", func() {
				secondElement.IsSelectedCall.ReturnSelected = false
				Expect(selection.Check()).To(Succeed())
				Expect(secondElement.ClickCall.Called).To(BeTrue())
			})

			Context("when the determining the selected status of any element fails", func() {
				It("should return an error", func() {
					secondElement.IsSelectedCall.Err = errors.New("some error")
					Expect(selection.Check()).To(MatchError("failed to retrieve state of selection 'CSS: #selector': some error"))
				})
			})

			Context("when clicking on the checkbox fails", func() {
				It("should return an error", func() {
					secondElement.ClickCall.Err = errors.New("some error")
					Expect(selection.Check()).To(MatchError("failed to click on selection 'CSS: #selector': some error"))
				})
			})
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Check()).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when any element fails to retrieve the 'type' attribute", func() {
			It("should return an error", func() {
				firstElement.GetAttributeCall.ReturnValue = "checkbox"
				secondElement.GetAttributeCall.Err = errors.New("some error")
				Expect(selection.Check()).To(MatchError("failed to retrieve type attribute of selection 'CSS: #selector': some error"))
			})
		})

		Context("when any element is not a checkbox", func() {
			It("should return an error", func() {
				firstElement.GetAttributeCall.ReturnValue = "checkbox"
				secondElement.GetAttributeCall.ReturnValue = "banana"
				Expect(selection.Check()).To(MatchError("selection 'CSS: #selector' does not refer to a checkbox"))
			})
		})
	})

	Describe("#Uncheck", func() {
		It("should successfully click on a checked checkbox", func() {
			firstElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.IsSelectedCall.ReturnSelected = true
			Expect(selection.Uncheck()).To(Succeed())
			Expect(firstElement.ClickCall.Called).To(BeFalse())
			Expect(secondElement.ClickCall.Called).To(BeTrue())
		})
	})

	Describe("#Select", func() {
		var (
			firstOptionBuses  []*mocks.Bus
			secondOptionBuses []*mocks.Bus
			firstOptions      []*api.Element
			secondOptions     []*api.Element
		)

		BeforeEach(func() {
			firstOptionBuses = []*mocks.Bus{{}, {}}
			secondOptionBuses = []*mocks.Bus{{}, {}}
			firstOptions = []*api.Element{
				{ID: "one", Session: &api.Session{Bus: firstOptionBuses[0]}},
				{ID: "two", Session: &api.Session{Bus: firstOptionBuses[1]}},
			}
			secondOptions = []*api.Element{
				{ID: "three", Session: &api.Session{Bus: secondOptionBuses[0]}},
				{ID: "four", Session: &api.Session{Bus: secondOptionBuses[1]}},
			}
			firstElement.GetElementsCall.ReturnElements = []*api.Element{firstOptions[0], firstOptions[1]}
			secondElement.GetElementsCall.ReturnElements = []*api.Element{secondOptions[0], secondOptions[1]}
		})

		It("should successfully retrieve the options with matching text for each selected element", func() {
			Expect(selection.Select("some text")).To(Succeed())
			Expect(firstElement.GetElementsCall.Selector.Using).To(Equal("xpath"))
			Expect(firstElement.GetElementsCall.Selector.Value).To(Equal(`./option[normalize-space()="some text"]`))
			Expect(secondElement.GetElementsCall.Selector.Using).To(Equal("xpath"))
			Expect(secondElement.GetElementsCall.Selector.Value).To(Equal(`./option[normalize-space()="some text"]`))
		})

		It("should successfully click on all options with matching text", func() {
			Expect(selection.Select("some text")).To(Succeed())
			Expect(firstOptionBuses[0].SendCall.Endpoint).To(Equal("element/one/click"))
			Expect(firstOptionBuses[1].SendCall.Endpoint).To(Equal("element/two/click"))
			Expect(secondOptionBuses[0].SendCall.Endpoint).To(Equal("element/three/click"))
			Expect(secondOptionBuses[1].SendCall.Endpoint).To(Equal("element/four/click"))
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Select("some text")).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when we fail to retrieve any option", func() {
			It("should return an error", func() {
				secondElement.GetElementsCall.Err = errors.New("some error")
				Expect(selection.Select("some text")).To(MatchError("failed to select specified option for selection 'CSS: #selector': some error"))
			})
		})

		Context("when any of the elements has no options with matching text", func() {
			It("should return an error", func() {
				secondElement.GetElementsCall.ReturnElements = []*api.Element{}
				Expect(selection.Select("some text")).To(MatchError(`no options with text "some text" found for selection 'CSS: #selector'`))
			})
		})

		Context("when the click fails for any of the options", func() {
			It("should return an error", func() {
				secondOptionBuses[1].SendCall.Err = errors.New("some error")
				Expect(selection.Select("some text")).To(MatchError(`failed to click on option with text "some text" for selection 'CSS: #selector': some error`))
			})
		})
	})

	Describe("#Submit", func() {
		It("should successfully submit all selected elements", func() {
			Expect(selection.Submit()).To(Succeed())
			Expect(firstElement.SubmitCall.Called).To(BeTrue())
			Expect(secondElement.SubmitCall.Called).To(BeTrue())
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Submit()).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when any submit fails", func() {
			It("should return an error", func() {
				secondElement.SubmitCall.Err = errors.New("some error")
				Expect(selection.Submit()).To(MatchError("failed to submit selection 'CSS: #selector': some error"))
			})
		})
	})

	// TODO: implement call tracking in mocks
	Describe("#Tap", func() {
		var (
			firstElement  *api.Element
			secondElement *api.Element
		)
		BeforeEach(func() {
			firstElement = &api.Element{}
			secondElement = &api.Element{}
			elementRepository.GetAtLeastOneCall.ReturnElements = []element.Element{firstElement, secondElement}
		})

		It("should successfully tap on all selected elements for each event type", func() {
			Expect(selection.Tap(SingleTap)).To(Succeed())
			Expect(session.TouchClickCall.Element).To(ExactlyEqual(secondElement))
			Expect(selection.Tap(DoubleTap)).To(Succeed())
			Expect(session.TouchDoubleClickCall.Element).To(ExactlyEqual(secondElement))
			Expect(selection.Tap(LongTap)).To(Succeed())
			Expect(session.TouchLongClickCall.Element).To(ExactlyEqual(secondElement))
		})

		Context("when the tap event is invalid", func() {
			It("should return an error", func() {
				err := selection.Tap(-1)
				Expect(err).To(MatchError("failed to perform tap on selection 'CSS: #selector': invalid tap event"))
			})
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Tap(SingleTap)).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})

		Context("when any tap fails", func() {
			It("should return an error", func() {
				session.TouchClickCall.Err = errors.New("some error")
				Expect(selection.Tap(SingleTap)).To(MatchError("failed to tap on selection 'CSS: #selector': some error"))
			})
		})
	})

	// TODO: implement call tracking in mocks
	Describe("#Touch", func() {
		It("should successfully instruct the session to touch using the provided offset for each event type", func() {
			firstElement.GetLocationCall.ReturnX = 100
			firstElement.GetLocationCall.ReturnY = 200
			secondElement.GetLocationCall.ReturnX = 300
			secondElement.GetLocationCall.ReturnY = 400

			Expect(selection.Touch(HoldFinger)).To(Succeed())
			Expect(session.TouchDownCall.X).To(Equal(300))
			Expect(session.TouchDownCall.Y).To(Equal(400))

			Expect(selection.Touch(ReleaseFinger)).To(Succeed())
			Expect(session.TouchUpCall.X).To(Equal(300))
			Expect(session.TouchUpCall.Y).To(Equal(400))

			Expect(selection.Touch(MoveFinger)).To(Succeed())
			Expect(session.TouchMoveCall.X).To(Equal(300))
			Expect(session.TouchMoveCall.Y).To(Equal(400))
		})

		Context("when retrieving an element's location fails", func() {
			It("should return an error", func() {
				secondElement.GetLocationCall.Err = errors.New("some error")
				Expect(selection.Touch(HoldFinger)).To(MatchError("failed to retrieve location of selection 'CSS: #selector': some error"))
			})
		})

		Context("when the touch event fails", func() {
			It("should return an error of each event type", func() {
				session.TouchDownCall.Err = errors.New("some touch down error")
				Expect(selection.Touch(HoldFinger)).To(MatchError("failed to flick finger on selection 'CSS: #selector': some touch down error"))

				session.TouchUpCall.Err = errors.New("some touch up error")
				Expect(selection.Touch(ReleaseFinger)).To(MatchError("failed to flick finger on selection 'CSS: #selector': some touch up error"))

				session.TouchMoveCall.Err = errors.New("some touch move error")
				Expect(selection.Touch(MoveFinger)).To(MatchError("failed to flick finger on selection 'CSS: #selector': some touch move error"))
			})
		})

		Context("when the touch event is invalid", func() {
			It("should return an error", func() {
				err := selection.Touch(-1)
				Expect(err).To(MatchError("failed to perform touch on selection 'CSS: #selector': invalid touch event"))
			})
		})
	})

	Describe("#FlickFinger", func() {
		var firstElement *api.Element

		BeforeEach(func() {
			firstElement = &api.Element{}
			elementRepository.GetExactlyOneCall.ReturnElement = firstElement
		})

		It("should successfully flick on the selected element", func() {
			Expect(selection.FlickFinger(100, 200, 300)).To(Succeed())
			Expect(session.TouchFlickCall.Element).To(ExactlyEqual(firstElement))
			Expect(session.TouchFlickCall.Offset).To(Equal(api.XYOffset{X: 100, Y: 200}))
			Expect(session.TouchFlickCall.Speed).To(Equal(api.ScalarSpeed(300)))
		})

		Context("when exactly one element is not returned", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				Expect(selection.FlickFinger(100, 200, 300)).To(MatchError("failed to select element from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the flick fails", func() {
			It("should return an error", func() {
				session.TouchFlickCall.Err = errors.New("some error")
				Expect(selection.FlickFinger(100, 200, 300)).To(MatchError("failed to flick finger on selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#ScrollFinger", func() {
		var firstElement *api.Element

		BeforeEach(func() {
			firstElement = &api.Element{}
			elementRepository.GetExactlyOneCall.ReturnElement = firstElement
		})

		It("should successfully scroll on the selected element", func() {
			Expect(selection.ScrollFinger(100, 200)).To(Succeed())
			Expect(session.TouchScrollCall.Element).To(ExactlyEqual(firstElement))
			Expect(session.TouchScrollCall.Offset).To(Equal(api.XYOffset{X: 100, Y: 200}))
		})

		Context("when exactly one element is not returned", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				Expect(selection.ScrollFinger(100, 200)).To(MatchError("failed to select element from selection 'CSS: #selector': some error"))
			})
		})

		Context("when the scroll fails", func() {
			It("should return an error", func() {
				session.TouchScrollCall.Err = errors.New("some error")
				Expect(selection.ScrollFinger(100, 200)).To(MatchError("failed to scroll finger on selection 'CSS: #selector': some error"))
			})
		})
	})
})
