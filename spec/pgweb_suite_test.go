package spec

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	"github.com/sosedoff/pgweb/spec/helpers"
)

func TestPgweb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pgweb Suite")
}

var agoutiDriver *agouti.WebDriver

var _ = BeforeSuite(func() {
	helpers.CreateBooktownDB()

	agoutiDriver = agouti.ChromeDriver()
	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	helpers.DropBooktownDb()

	Expect(agoutiDriver.Stop()).To(Succeed())
})

var page *agouti.Page

var _ = BeforeEach(func() {
	var err error

	page, err = agoutiDriver.NewPage()
	page.Size(1366, 900)
	Expect(page.Navigate("http://localhost:8081")).To(Succeed())

	if visible, _ := page.Find("#close_connection").Visible(); visible {
		Expect(page.Find("#close_connection").Click()).To(Succeed())
		page.ConfirmPopup()
	}

	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	Expect(page.Destroy()).To(Succeed())
})
