package spec

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)


var _ = Describe("DbConnection", func() {
	var (
		selErrorBlock = "#connection_error"
		selCurrentDb = "#current_database"
	)

	var	txtConnectBtn = "Connect"

	var errorMsg = "pq: password authentication failed for user \"postgres\""
	var dbName = "booktown"

	It("connects to DB by connection string tab", func() {
		var (
			correctConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
						serverUser, serverPassword, serverHost, serverPort, serverDatabase)
			wrongConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				serverUser, "wrongpassword", serverHost, serverPort, serverDatabase)
		)

		Expect(page.Find("#connection_scheme").Click()).To(Succeed())

		Context("using wrong password", func() {
			page.Find("#connection_url").Fill(wrongConnStr)
			Expect(page.FindByButton(txtConnectBtn).Click()).To(Succeed())


			// TODO: find out the problem of chaotic test failure
			//
			// if I remove screenshot statement. the spec will fail
			// Maybe it is related to timeouts when we dealing AJAX;
			// clicking an element that will trigger AJAX. which take
			// arbitrary long time (see Codeception waitFor functions)
			screenshot(page, "scheme_wrong_password_after_connect")
			Expect(page.Find(selErrorBlock)).To(HaveText(errorMsg))
		})


		Context("using correct password", func() {
			page.Find("#connection_url").Fill(correctConnStr)
			Expect(page.FindByButton(txtConnectBtn).Click()).To(Succeed())
			Expect(page.Find(selCurrentDb)).To(BeVisible())
			Expect(page.Find(selCurrentDb)).Should(HaveText(dbName))
		})

	})

	It("connects to DB by Standard tab", func() {
		// Filling the form
		data := map[string]string {
			"#pg_user": serverUser,
			"#pg_password": "wrongpassword",
			"#pg_host": serverHost,
			"#pg_db": serverDatabase,
		}


		Context("using wrong password", func() {
			for selector, value := range data {
				page.Find(selector).Fill(value)
			}

			page.Find("#connection_ssl").Select("disable")

			Expect(page.FindByButton(txtConnectBtn).Click()).To(Succeed())
			Expect(page.Find(selErrorBlock)).To(HaveText(errorMsg))
		})


		Context("using correct password", func() {
			page.Find("#pg_password").Fill(serverPassword)

			Expect(page.FindByButton(txtConnectBtn).Click()).To(Succeed())
			screenshot(page, "standard_correct_password_after_connect")
			Expect(page.Find(selCurrentDb)).To(HaveText(dbName))
		})

	})
})
