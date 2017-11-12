package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("NavigateToApp", func() {

	It("should navigate to pgweb page", func () {
		Expect(page.Navigate("http://localhost:8081")).To(Succeed())
		Expect(page).To(HaveTitle("pgweb"))
		Expect(page.Find(".connection-settings h1")).To(HaveText("pgweb"))
	})
})



