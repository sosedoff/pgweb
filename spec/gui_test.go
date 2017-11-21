package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"

	"github.com/sosedoff/pgweb/spec/helpers"
)

var _ = Describe("Gui", func() {

	BeforeEach(func() {
		helpers.ConnectByStandardTab(page)
	})

	Context("Tabs", func() {
		BeforeEach(func() {
			sel := page.Find(helpers.TableBookSelector)
			sel.Click()
		})

		FIt("clicks on Rows tab", func() {
			page.Find(helpers.TabRowsSelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.Find("#results > tr td:nth-child(2)").Text()).To(Equal("The Shining"))
		})
	})
})
