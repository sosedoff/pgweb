package agouti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/mocks"
)

var _ = Describe("MultiSelection", func() {
	var (
		bus       *mocks.Bus
		session   *api.Session
		selection *MultiSelection
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		session = &api.Session{Bus: bus}
		selection = NewTestMultiSelection(session, nil, "#selector")
	})

	Describe("#At", func() {
		It("should add an index to the current selection", func() {
			Expect(selection.At(4).String()).To(Equal("selection 'CSS: #selector [4]'"))
		})

		It("should provide the selectable's session to the element repository", func() {
			bus.SendCall.Result = `[{"ELEMENT": "some-id"}]`
			elements, _ := selection.At(0).Find("b").Elements()
			Expect(elements[0].ID).To(Equal("some-id"))
		})
	})
})
