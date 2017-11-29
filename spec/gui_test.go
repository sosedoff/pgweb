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
		Eventually(page.Find(helpers.CurrentDbSelector), 10).Should(BeVisible())
	})

	FContext("Tabs", func() {
		BeforeEach(func() {
			page.Find(helpers.TableBookSelector).Click()
		})

		It("clicks on Rows tab", func() {
			page.Find(helpers.TabRowsSelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.All("#results tr td").At(1).Text()).To(Equal("The Shining"))
		})
	})
})
