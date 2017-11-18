package target_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/api"
	. "github.com/sclevine/agouti/internal/target"
)

var _ = Describe("Selector", func() {
	Describe("#String", func() {
		It("should return valid string suffixes for the Selector", func() {
			Expect(Selector{Type: CSS, Value: "value"}.String()).To(Equal("CSS: value"))
			Expect(Selector{Type: CSS, Value: "value", Single: true}.String()).To(Equal("CSS: value [single]"))
			Expect(Selector{Type: CSS, Value: "value", Indexed: true, Index: 4}.String()).To(Equal("CSS: value [4]"))
		})

		It("should return valid string formatting for the Selector", func() {
			Expect(Selector{Type: CSS, Value: "value"}.String()).To(Equal("CSS: value"))
			Expect(Selector{Type: XPath, Value: "value"}.String()).To(Equal("XPath: value"))
			Expect(Selector{Type: Link, Value: "value"}.String()).To(Equal(`Link: "value"`))
			Expect(Selector{Type: Label, Value: "value"}.String()).To(Equal(`Label: "value"`))
			Expect(Selector{Type: Button, Value: "value"}.String()).To(Equal(`Button: "value"`))
			Expect(Selector{Type: Name, Value: "value"}.String()).To(Equal(`Name: "value"`))

		})
	})

	Describe("#API", func() {
		It("should return an API-consumable version of the Selector", func() {
			Expect(Selector{Type: CSS, Value: "value"}.API()).To(Equal(api.Selector{Using: "css selector", Value: "value"}))
			Expect(Selector{Type: XPath, Value: "value"}.API()).To(Equal(api.Selector{Using: "xpath", Value: "value"}))
			Expect(Selector{Type: Link, Value: "value"}.API()).To(Equal(api.Selector{Using: "link text", Value: "value"}))
			Expect(Selector{Type: Label, Value: "value"}.API()).To(Equal(api.Selector{Using: "xpath", Value: `//input[@id=(//label[normalize-space()="value"]/@for)] | //label[normalize-space()="value"]/input`}))
			Expect(Selector{Type: Button, Value: "value"}.API()).To(Equal(api.Selector{Using: "xpath", Value: `//input[@type="submit" or @type="button"][normalize-space(@value)="value"] | //button[normalize-space()="value"]`}))
			Expect(Selector{Type: Name, Value: "value"}.API()).To(Equal(api.Selector{Using: "name", Value: "value"}))
		})
	})
})
