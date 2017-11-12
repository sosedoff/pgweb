package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("ConnectionOptions", func() {

	Context("Switching connections options tabs", func () {
		It("clicks on Standard tab", func () {
			Expect(page.Find("#connection_standard").Click()).To(Succeed())
			Expect(page.Find("#pg_user")).Should(BeVisible())
			Expect(page.Find("#pg_host")).Should(BeVisible())
			Expect(page.Find("#pg_db")).Should(BeVisible())

			Expect(page.Find("#connection_url")).ShouldNot(BeVisible())

			Expect(page.Find("#ssh_host")).ShouldNot(BeVisible())
		})

		It("clicks on Scheme tab", func() {
			Expect(page.Find("#connection_scheme").Click()).To(Succeed())
			Expect(page.Find("#connection_url")).Should(BeVisible())

			Expect(page.Find("#ssh_host")).ShouldNot(BeVisible())
		})


		It("clicks on SSH tab", func () {
			Expect(page.Find("#connection_ssh").Click()).To(Succeed())
			Expect(page.Find("#ssh_host")).Should(BeVisible())

			Expect(page.Find("#connection_url")).ShouldNot(BeVisible())
		})
	})
})
