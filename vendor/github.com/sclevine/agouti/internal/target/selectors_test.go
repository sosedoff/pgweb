package target_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/internal/target"
)

var _ = Describe("Selectors", func() {
	var selectors Selectors

	BeforeEach(func() {
		selectors = Selectors{}
	})

	Describe("#Append", func() {
		Context("when the provided selector is not a CSS selector", func() {
			It("should append a new selector", func() {
				selectors := selectors.Append(XPath, "//selector")
				Expect(selectors.String()).To(Equal("XPath: //selector"))
				Expect(selectors.Append(Link, "some link").String()).To(Equal(`XPath: //selector | Link: "some link"`))
			})
		})

		Context("when the provided selector is a CSS selector", func() {
			Context("when the last selector is an unindexed CSS selector", func() {
				It("should modify the last CSS selector to include the new selector", func() {
					Expect(selectors.Append(CSS, "#selector").Append(CSS, "#subselector").String()).To(Equal("CSS: #selector #subselector"))
				})
			})

			Context("when there are no selectors", func() {
				It("should append a new selector", func() {
					Expect(selectors.Append(CSS, "#selector").String()).To(Equal("CSS: #selector"))
				})
			})

			Context("when the last selector is a non-CSS selector", func() {
				It("should append a new selector", func() {
					xpath := selectors.Append(XPath, "//selector")
					Expect(xpath.Append(CSS, "#subselector").String()).To(Equal("XPath: //selector | CSS: #subselector"))
				})
			})

			Context("when the last selector is an indexed selector", func() {
				It("should append a new selector", func() {
					Expect(selectors.Append(CSS, "#selector").At(0).Append(CSS, "#subselector").String()).To(Equal("CSS: #selector [0] | CSS: #subselector"))
				})
			})

			Context("when the last selector is a single-element-only selector", func() {
				It("should append a new selector", func() {
					Expect(selectors.Append(CSS, "#selector").Single().Append(CSS, "#subselector").String()).To(Equal("CSS: #selector [single] | CSS: #subselector"))
				})
			})
		})
	})

	Describe("#At", func() {
		Context("when called on a selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(selectors.At(1).String()).To(Equal(""))
			})
		})

		Context("when called on a selection with selectors", func() {
			It("should select an index of the current selection", func() {
				Expect(selectors.Append(CSS, "#selector").At(1).String()).To(Equal("CSS: #selector [1]"))
			})
		})
	})

	Describe("#Single", func() {
		Context("when called on a selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(selectors.Single().String()).To(Equal(""))
			})
		})

		Context("when called on a selection with selectors", func() {
			It("should select a single element of the current selection", func() {
				Expect(selectors.Append(CSS, "#selector").Single().String()).To(Equal("CSS: #selector [single]"))
			})
		})
	})

	Describe("selectors are always copied", func() {
		Context("when two CSS selections are created from the same XPath parent", func() {
			It("should not overwrite the first created child", func() {
				parent := selectors.Append(XPath, "//one").Append(XPath, "//two").Append(XPath, "//parent")
				firstChild := parent.Append(CSS, "#firstChild")
				parent.Append(CSS, "#secondChild")
				Expect(firstChild.String()).To(Equal("XPath: //one | XPath: //two | XPath: //parent | CSS: #firstChild"))
			})
		})
	})
})
