package element_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/api"
	. "github.com/sclevine/agouti/internal/element"
	. "github.com/sclevine/agouti/internal/matchers"
	"github.com/sclevine/agouti/internal/mocks"
	"github.com/sclevine/agouti/internal/target"
)

var _ = Describe("ElementRepository", func() {
	var (
		client     *mocks.Session
		repository *Repository
	)

	BeforeEach(func() {
		client = &mocks.Session{}
		repository = &Repository{Client: client}
	})

	Describe("#GetAtLeastOne", func() {
		BeforeEach(func() {
			repository.Selectors = target.Selectors{target.Selector{}}
		})

		Context("when the client fails to retrieve any elements", func() {
			It("should fail with an error", func() {
				client.GetElementsCall.Err = errors.New("some error")
				_, err := repository.GetAtLeastOne()
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when the client retrieves zero elements", func() {
			It("should fail with an error", func() {
				client.GetElementsCall.ReturnElements = []*api.Element{}
				_, err := repository.GetAtLeastOne()
				Expect(err).To(MatchError("no elements found"))
			})
		})

		Context("when the client retrieves at least one element", func() {
			It("should successfully return those elements", func() {
				element := &api.Element{}
				client.GetElementsCall.ReturnElements = []*api.Element{element}
				elements, err := repository.GetAtLeastOne()
				Expect(err).NotTo(HaveOccurred())
				Expect(elements[0]).To(ExactlyEqual(element))
			})
		})
	})

	Describe("#GetExactlyOne", func() {
		BeforeEach(func() {
			repository.Selectors = target.Selectors{target.Selector{}}
		})

		Context("when the client retrieves zero elements", func() {
			It("should return an error", func() {
				client.GetElementsCall.ReturnElements = []*api.Element{}
				_, err := repository.GetExactlyOne()
				Expect(err).To(MatchError("no elements found"))
			})
		})

		Context("when the client retrieves more than one element", func() {
			It("should return an error", func() {
				client.GetElementsCall.ReturnElements = []*api.Element{{}, {}}
				_, err := repository.GetExactlyOne()
				Expect(err).To(MatchError("method does not support multiple elements (2)"))
			})
		})

		Context("when the client retrieves exactly one element", func() {
			It("should successfully return that element", func() {
				element := &api.Element{}
				client.GetElementsCall.ReturnElements = []*api.Element{element}
				Expect(repository.GetExactlyOne()).To(ExactlyEqual(element))
			})
		})
	})

	Describe("#Get", func() {
		var (
			firstParentBus     *mocks.Bus
			firstParent        *api.Element
			secondParentBus    *mocks.Bus
			secondParent       *api.Element
			children           []Element
			parentSelector     target.Selector
			parentSelectorJSON string
			childSelector      target.Selector
			childSelectorJSON  string
		)

		BeforeEach(func() {
			firstParentBus = &mocks.Bus{}
			firstParent = &api.Element{ID: "first parent", Session: &api.Session{Bus: firstParentBus}}
			secondParentBus = &mocks.Bus{}
			secondParent = &api.Element{ID: "second parent", Session: &api.Session{Bus: secondParentBus}}
			children = []Element{
				Element(&api.Element{ID: "first child", Session: &api.Session{Bus: firstParentBus}}),
				Element(&api.Element{ID: "second child", Session: &api.Session{Bus: firstParentBus}}),
				Element(&api.Element{ID: "third child", Session: &api.Session{Bus: secondParentBus}}),
				Element(&api.Element{ID: "fourth child", Session: &api.Session{Bus: secondParentBus}}),
			}
			firstParentBus.SendCall.Result = `[{"ELEMENT": "first child"}, {"ELEMENT": "second child"}]`
			secondParentBus.SendCall.Result = `[{"ELEMENT": "third child"}, {"ELEMENT": "fourth child"}]`
			client.GetElementsCall.ReturnElements = []*api.Element{firstParent, secondParent}
			parentSelector = target.Selector{Type: target.CSS, Value: "parents"}
			parentSelectorJSON = `{"using": "css selector", "value": "parents"}`
			childSelector = target.Selector{Type: target.XPath, Value: "children"}
			childSelectorJSON = `{"using": "xpath", "value": "children"}`
			repository.Selectors = target.Selectors{parentSelector, childSelector}
		})

		Context("when all unindexed elements are successfully retrieved", func() {
			It("should retrieve all parent and child elements and return the child elements", func() {
				Expect(repository.Get()).To(Equal(children))
				Expect(client.GetElementsCall.Selector).To(Equal(parentSelector.API()))
				Expect(firstParentBus.SendCall.BodyJSON).To(MatchJSON(childSelectorJSON))
				Expect(secondParentBus.SendCall.BodyJSON).To(MatchJSON(childSelectorJSON))
			})
		})

		Context("when a non-zero-indexed element is successfully retrieved", func() {
			BeforeEach(func() {
				parentSelector.Index = 1
				parentSelector.Indexed = true
				childSelector.Index = 1
				childSelector.Indexed = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
			})

			It("should retrieve the non-zero-indexed elements and return the child element", func() {
				Expect(repository.Get()).To(Equal([]Element{children[3]}))
				Expect(firstParentBus.SendCall.BodyJSON).To(BeEmpty())
				Expect(secondParentBus.SendCall.BodyJSON).To(MatchJSON(childSelectorJSON))
				Expect(client.GetElementsCall.Selector).To(Equal(parentSelector.API()))
			})
		})

		Context("when a zero-indexed element is successfully retrieved", func() {
			BeforeEach(func() {
				firstParentBus.SendCall.Result = `{"ELEMENT": "first child"}`
				client.GetElementCall.ReturnElement = firstParent
				parentSelector.Index = 0
				parentSelector.Indexed = true
				childSelector.Index = 0
				childSelector.Indexed = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
			})

			It("should retrieve the first parent and child elements and return the child element", func() {
				Expect(repository.Get()).To(Equal([]Element{children[0]}))
				Expect(firstParentBus.SendCall.BodyJSON).To(MatchJSON(childSelectorJSON))
				Expect(client.GetElementCall.Selector).To(Equal(parentSelector.API()))
			})
		})

		Context("when single-element-only elements are successfully retrieved", func() {
			BeforeEach(func() {
				firstParentBus.SendCall.Result = `[{"ELEMENT": "first child"}]`
				client.GetElementsCall.ReturnElements = []*api.Element{firstParent}
				parentSelector.Single = true
				childSelector.Single = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
			})

			It("should retrieve the parent and child elements and return the child element", func() {
				Expect(repository.Get()).To(Equal([]Element{children[0]}))
				Expect(firstParentBus.SendCall.BodyJSON).To(MatchJSON(childSelectorJSON))
				Expect(client.GetElementsCall.Selector).To(Equal(parentSelector.API()))
			})
		})

		Context("when there is no selection", func() {
			It("should return an error", func() {
				repository.Selectors = target.Selectors{}
				_, err := repository.Get()
				Expect(err).To(MatchError("empty selection"))
			})
		})

		Context("when a single-element-only parent selection refers to multiple parents", func() {
			It("should return an error", func() {
				parentSelector.Single = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("ambiguous find"))
			})
		})

		Context("when a single-element-only parent selection refers to no parents", func() {
			It("should return an error", func() {
				parentSelector.Single = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				client.GetElementsCall.ReturnElements = []*api.Element{}
				_, err := repository.Get()
				Expect(err).To(MatchError("element not found"))
			})
		})

		Context("when any single-element-only child selection refers to multiple child elements", func() {
			It("should return an error", func() {
				childSelector.Single = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				firstParentBus.SendCall.Result = `[{"ELEMENT": "first child"}]`
				_, err := repository.Get()
				Expect(err).To(MatchError("ambiguous find"))
			})
		})

		Context("when any single-element-only child selection refers to no child elements", func() {
			It("should return an error", func() {
				childSelector.Single = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				firstParentBus.SendCall.Result = `[]`
				_, err := repository.Get()
				Expect(err).To(MatchError("element not found"))
			})
		})

		Context("when the parent selection index is out of range", func() {
			It("should return an error", func() {
				parentSelector.Index = 2
				parentSelector.Indexed = true
				repository.Selectors = target.Selectors{parentSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("element index out of range"))
			})
		})

		Context("when child selection indices are out of range", func() {
			It("should return an error", func() {
				parentSelector.Index = 1
				parentSelector.Indexed = true
				childSelector.Index = 2
				childSelector.Indexed = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("element index out of range"))
			})
		})

		Context("when a zero-indexed parent selection element does not exist", func() {
			It("should return an error", func() {
				client.GetElementCall.Err = errors.New("some error")
				parentSelector.Index = 0
				parentSelector.Indexed = true
				repository.Selectors = target.Selectors{parentSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when a zero-indexed child selection element does not exist", func() {
			It("should return an error", func() {
				firstParentBus.SendCall.Err = errors.New("some error")
				client.GetElementCall.ReturnElement = firstParent
				parentSelector.Index = 0
				parentSelector.Indexed = true
				childSelector.Index = 0
				childSelector.Indexed = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when retrieving any non-zero-indexed parent selection element fails", func() {
			It("should return an error", func() {
				client.GetElementsCall.Err = errors.New("some error")
				parentSelector.Index = 1
				parentSelector.Indexed = true
				repository.Selectors = target.Selectors{parentSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when retrieving any non-zero-indexed child selection element fails", func() {
			It("should return an error", func() {
				firstParentBus.SendCall.Err = errors.New("some error")
				childSelector.Index = 1
				childSelector.Indexed = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when retrieving any single-element-only parent selection element fails", func() {
			It("should return an error", func() {
				client.GetElementsCall.Err = errors.New("some error")
				parentSelector.Single = true
				repository.Selectors = target.Selectors{parentSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when retrieving any single-element-only child selection element fails", func() {
			It("should return an error", func() {
				firstParentBus.SendCall.Err = errors.New("some error")
				childSelector.Single = true
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when retrieving any unindexed parent elements fails", func() {
			It("should return an error", func() {
				client.GetElementsCall.Err = errors.New("some error")
				repository.Selectors = target.Selectors{parentSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when retrieving any unindexed child elements fails", func() {
			It("should return an error", func() {
				secondParentBus.SendCall.Err = errors.New("some error")
				repository.Selectors = target.Selectors{parentSelector, childSelector}
				_, err := repository.Get()
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
