package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/internal/mocks"
	. "github.com/sclevine/agouti/internal/matchers"
)

var _ = Describe("Element", func() {
	var (
		bus     *mocks.Bus
		session *Session
		element *Element
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		session = &Session{bus}
		element = &Element{"some-id", session}
	})

	Describe("#Send", func() {
		It("should successfully send a request to the provided endpoint", func() {
			Expect(element.Send("method", "endpoint", "body", nil)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("method"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/endpoint"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`"body"`))
		})

		It("should retrieve the result", func() {
			var result string
			bus.SendCall.Result = `"some result"`
			Expect(element.Send("method", "endpoint", "body", &result)).To(Succeed())
			Expect(result).To(Equal("some result"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				err := element.Send("method", "endpoint", "body", nil)
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElement", func() {
		It("should successfully send a POST request to the element endpoint", func() {
			_, err := element.GetElement(Selector{"css selector", "#selector"})
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/element"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		It("should return an element with an ID and session", func() {
			bus.SendCall.Result = `{"ELEMENT": "some-id"}`
			singleElement, err := element.GetElement(Selector{})
			Expect(err).NotTo(HaveOccurred())
			Expect(singleElement.ID).To(Equal("some-id"))
			Expect(singleElement.Session).To(ExactlyEqual(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetElement(Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		It("should successfully send a POST request to the elements endpoint", func() {
			_, err := element.GetElements(Selector{"css selector", "#selector"})
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/elements"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		It("should return a slice of elements with IDs and sessions", func() {
			bus.SendCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err := element.GetElements(Selector{"css selector", "#selector"})
			Expect(err).NotTo(HaveOccurred())
			Expect(elements[0].ID).To(Equal("some-id"))
			Expect(elements[0].Session).To(ExactlyEqual(session))
			Expect(elements[1].ID).To(Equal("some-other-id"))
			Expect(elements[1].Session).To(ExactlyEqual(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetElements(Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetText", func() {
		It("should successfully send a GET request to the text endpoint", func() {
			_, err := element.GetText()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/text"))
		})

		It("should return the visible text on the element", func() {
			bus.SendCall.Result = `"some text"`
			text, err := element.GetText()
			Expect(err).NotTo(HaveOccurred())
			Expect(text).To(Equal("some text"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to retrieve the text", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetText()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetName", func() {
		It("should successfully send a GET request to the name endpoint", func() {
			_, err := element.GetName()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/name"))
		})

		It("should return the tag name of the element", func() {
			bus.SendCall.Result = `"some-name"`
			text, err := element.GetName()
			Expect(err).NotTo(HaveOccurred())
			Expect(text).To(Equal("some-name"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to retrieve the tag name", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetName()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetAttribute", func() {
		It("should successfully send a GET request to the attribute/some-name endpoint", func() {
			_, err := element.GetAttribute("some-name")
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/attribute/some-name"))
		})

		It("should return the value of the attribute", func() {
			bus.SendCall.Result = `"some value"`
			value, err := element.GetAttribute("")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("some value"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetAttribute("some-name")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetCSS", func() {
		It("should successfully send a GET request to the css/some-property endpoint", func() {
			_, err := element.GetCSS("some-property")
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/css/some-property"))
		})

		It("should return the value of the CSS property", func() {
			bus.SendCall.Result = `"some value"`
			value, err := element.GetCSS("some-property")
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("some value"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.GetCSS("some-property")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Click", func() {
		It("should successfully send a POST request to the click endpoint", func() {
			Expect(element.Click()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/click"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(element.Click()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Clear", func() {
		It("should successfully send a POST request to the clear endpoint", func() {
			Expect(element.Clear()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/clear"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(element.Clear()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Value", func() {
		It("should successfully send a POST request to the value endpoint", func() {
			Expect(element.Value("text")).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/value"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"value": ["t", "e", "x", "t"]}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(element.Value("text")).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsSelected", func() {
		It("should successfully send a GET request to the selected endpoint", func() {
			_, err := element.IsSelected()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/selected"))
		})

		It("should return the selected status", func() {
			bus.SendCall.Result = "true"
			value, err := element.IsSelected()
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(BeTrue())
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.IsSelected()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsDisplayed", func() {
		It("should successfully send a GET request to the displayed endpoint", func() {
			_, err := element.IsDisplayed()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/displayed"))
		})

		It("should return the displayed status", func() {
			bus.SendCall.Result = "true"
			value, err := element.IsDisplayed()
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(BeTrue())
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.IsDisplayed()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsEnabled", func() {
		It("should successfully send a GET request to the enabled endpoint", func() {
			_, err := element.IsEnabled()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/enabled"))
		})

		It("should return the enabled status", func() {
			bus.SendCall.Result = "true"
			value, err := element.IsEnabled()
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(BeTrue())
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.IsEnabled()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Submit", func() {
		It("should successfully send a POST request to the submit endpoint", func() {
			Expect(element.Submit()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/submit"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(element.Submit()).To(MatchError("some error"))
			})
		})
	})

	Describe("#IsEqualTo", func() {
		It("should successfully send a GET request to the equals/other-id endpoint", func() {
			otherElement := &Element{"other-id", session}
			_, err := element.IsEqualTo(otherElement)
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/equals/other-id"))
		})

		It("should return whether the elements are equal", func() {
			bus.SendCall.Result = "true"
			equal, err := element.IsEqualTo(&Element{})
			Expect(err).NotTo(HaveOccurred())
			Expect(equal).To(BeTrue())
		})

		Context("when the other element is nil", func() {
			It("should return an error", func() {
				_, err := element.IsEqualTo(nil)
				Expect(err).To(MatchError("nil element is invalid"))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := element.IsEqualTo(&Element{})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetLocation", func() {
		It("should successfully send a GET request to the location endpoint", func() {
			_, _, err := element.GetLocation()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/location"))
		})

		It("should return the rounded location of the element", func() {
			bus.SendCall.Result = `{"x": 100.7, "y": 200}`
			x, y, err := element.GetLocation()
			Expect(err).NotTo(HaveOccurred())
			Expect(x).To(Equal(101))
			Expect(y).To(Equal(200))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to retrieve the location", func() {
				bus.SendCall.Err = errors.New("some error")
				_, _, err := element.GetLocation()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetSize", func() {
		It("should successfully send a GET request to the size endpoint", func() {
			_, _, err := element.GetSize()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/some-id/size"))
		})

		It("should return the rounded size of the element", func() {
			bus.SendCall.Result = `{"width": 100.7, "height": 200}`
			width, height, err := element.GetSize()
			Expect(err).NotTo(HaveOccurred())
			Expect(width).To(Equal(101))
			Expect(height).To(Equal(200))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to retrieve the location", func() {
				bus.SendCall.Err = errors.New("some error")
				_, _, err := element.GetSize()
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
