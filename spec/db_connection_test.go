package spec

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("DbConnection", func() {
	var page *agouti.Page

	BeforeEach(func() {
		var err error
		page, err = agoutiDriver.NewPage()
		Expect(page.Navigate("http://localhost:8081")).To(Succeed())
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})

	It("connects to DB by connection string tab", func() {
		var (
			correctConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
						serverUser, serverPassword, serverHost, serverPort, serverDatabase)
			wrongConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				serverUser, "wrongpassword", serverHost, serverPort, serverDatabase)
		)

		Expect(page.Find("#connection_scheme").Click()).To(Succeed())

		page.Find("#connection_url").Fill(wrongConnStr)
		Expect(page.FindByButton("Connect").Click()).To(Succeed())
		Expect(page.Find("#connection_error")).To(
								HaveText("pq: password authentication failed for user \"postgres\""))


		page.Find("#connection_url").Fill(correctConnStr)
		Expect(page.FindByButton("Connect").Click()).To(Succeed())
		Expect(page.Find("#current_database")).To(BeVisible())
		Expect(page.Find("#current_database")).Should(HaveText("booktown"))
	})
})
