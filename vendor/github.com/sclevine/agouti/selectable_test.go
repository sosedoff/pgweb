package agouti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/mocks"
)

var _ = Describe("Selectable", func() {
	var (
		bus     *mocks.Bus
		session *api.Session
		page    *Page
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		session = &api.Session{Bus: bus}
		page = NewTestPage(session)
		bus.SendCall.Result = `[{"ELEMENT": ""}]`
	})

	Describe("#Find", func() {
		It("should apply a single CSS selector and return a selection with the same session", func() {
			Expect(page.Find("selector").String()).To(Equal("selection 'CSS: selector [single]'"))
			Expect(page.Find("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FindByXPath", func() {
		It("should apply a single XPath selector and return a selection with the same session", func() {
			Expect(page.FindByXPath("selector").String()).To(Equal("selection 'XPath: selector [single]'"))
			Expect(page.FindByXPath("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FindByLink", func() {
		It("should apply a single link selector and return a selection with the same session", func() {
			Expect(page.FindByLink("selector").String()).To(Equal(`selection 'Link: "selector" [single]'`))
			Expect(page.FindByLink("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FindByLabel", func() {
		It("should apply a single label selector and return a selection with the same session", func() {
			Expect(page.FindByLabel("selector").String()).To(Equal(`selection 'Label: "selector" [single]'`))
			Expect(page.FindByLabel("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FindByButton", func() {
		It("should apply a single button text selector and return a selection with the same session", func() {
			Expect(page.FindByButton("selector").String()).To(Equal(`selection 'Button: "selector" [single]'`))
			Expect(page.FindByButton("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FindByClass", func() {
		It("should apply a single class selector and return a selection with the same session", func() {
			Expect(page.FindByClass("selector").String()).To(Equal(`selection 'Class: selector [single]'`))
			Expect(page.FindByClass("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FindByID", func() {
		It("should apply a single ID selector and return a selection with the same session", func() {
			Expect(page.FindByID("selector").String()).To(Equal(`selection 'ID: selector [single]'`))
			Expect(page.FindByID("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#First", func() {
		It("should apply a zero-indexed CSS selector and return a selection with the same session", func() {
			Expect(page.First("selector").String()).To(Equal("selection 'CSS: selector [0]'"))
			Expect(page.First("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FirstByXPath", func() {
		It("should apply a zero-indexed XPath selector and return a selection with the same session", func() {
			Expect(page.FirstByXPath("selector").String()).To(Equal("selection 'XPath: selector [0]'"))
			Expect(page.FirstByXPath("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FirstByLink", func() {
		It("should apply a zero-indexed link selector and return a selection with the same session", func() {
			Expect(page.FirstByLink("selector").String()).To(Equal(`selection 'Link: "selector" [0]'`))
			Expect(page.FirstByLink("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FirstByLabel", func() {
		It("should apply a zero-indexed label selector and return a selection with the same session", func() {
			Expect(page.FirstByLabel("selector").String()).To(Equal(`selection 'Label: "selector" [0]'`))
			Expect(page.FirstByLabel("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FirstByButton", func() {
		It("should apply a zero-indexed button text selector and return a selection with the same session", func() {
			Expect(page.FirstByButton("selector").String()).To(Equal(`selection 'Button: "selector" [0]'`))
			Expect(page.FirstByButton("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#FirstByClass", func() {
		It("should apply a zero-indexed class selector and return a selection with the same session", func() {
			Expect(page.FirstByClass("selector").String()).To(Equal(`selection 'Class: selector [0]'`))
			Expect(page.FirstByClass("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#All", func() {
		It("should apply an un-indexed CSS selector and return a selection with the same session", func() {
			Expect(page.All("selector").String()).To(Equal("selection 'CSS: selector'"))
			Expect(page.All("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#AllByXPath", func() {
		It("should apply an un-indexed XPath selector and return a selection with the same session", func() {
			Expect(page.AllByXPath("selector").String()).To(Equal("selection 'XPath: selector'"))
			Expect(page.AllByXPath("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#AllByLink", func() {
		It("should apply an un-indexed link selector and return a selection with the same session", func() {
			Expect(page.AllByLink("selector").String()).To(Equal(`selection 'Link: "selector"'`))
			Expect(page.AllByLink("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#AllByLabel", func() {
		It("should apply an un-indexed label selector and return a selection with the same session", func() {
			Expect(page.AllByLabel("selector").String()).To(Equal(`selection 'Label: "selector"'`))
			Expect(page.AllByLabel("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#AllByButton", func() {
		It("should apply an un-indexed button text selector and return a selection with the same session", func() {
			Expect(page.AllByButton("selector").String()).To(Equal(`selection 'Button: "selector"'`))
			Expect(page.AllByButton("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#AllByClass", func() {
		It("should apply an un-indexed class selector and return a selection with the same session", func() {
			Expect(page.AllByClass("selector").String()).To(Equal(`selection 'Class: selector'`))
			Expect(page.AllByClass("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})

	Describe("#AllByID", func() {
		It("should apply an un-indexed id selector and return a selection with the same session", func() {
			Expect(page.AllByID("selector").String()).To(Equal(`selection 'ID: selector'`))
			Expect(page.AllByID("selector").Elements()).To(ContainElement(&api.Element{Session: session}))
		})
	})
})
