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

	Context("Tabs", func() {
		BeforeEach(func() {
			page.Size(1366, 900)
			page.Find(helpers.TableBookSelector).Click()
		})

		It("clicks on Rows tab", func() {
			page.Find(helpers.TabRowsSelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.All("#results tr td").At(1).Text()).To(Equal("The Shining"))
		})

		It("clicks on Structure editor tab", func() {
			page.Find(helpers.TabStructureSelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.All("#results tr").At(3).All("td").At(0).Text()).To(
				Equal("author_id"))
		})

		It("clicks on Indexes tab", func() {
			page.Find(helpers.TabIndexesSelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.All("#results tr td").At(0).Text()).To(Equal("books_id_pkey"))
		})

		It("clicks on Constraints tab", func() {
			page.Find(helpers.TabConstraintsSelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.All("#results tr td").At(0).Text()).To(Equal("PRIMARY KEY (id)"))
		})

		It("clicks on SQL query tab", func() {
			page.Find(helpers.TabQuerySelector).Click()
			Expect(page.Find(".actions #run")).Should(BeVisible())
		})

		It("clicks on History tab", func() {
			page.Find(helpers.TabHistorySelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.Find("#results tr td").Text()).To(Equal("No records found"))
		})

		It("clicks on Activity tab", func() {
			page.Find(helpers.TabActivitySelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.All("#results tr td").At(0).Text()).To(Equal("booktown"))
		})

		It("clicks on Connection tab", func() {
			page.Find(helpers.TabConnectionSelector).Click()
			Expect(page.Find("#results")).Should(BeVisible())
			Expect(page.All("#results tr td").At(0).Text()).To(Equal("current_database"))
			Expect(page.All("#results tr").At(1).All("td").At(1).Text()).To(
				Equal("booktown"))
		})

	})

	//Context("SQL editor", func() {
	//
	//})
})
