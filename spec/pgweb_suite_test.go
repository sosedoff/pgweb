package spec

import (
	"os"
	"testing"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
)


var (
	serverHost     string
	serverPort     string
	serverUser     string
	serverPassword string
	serverDatabase string
)


func screenshot(page *agouti.Page, name string) {
	page.Screenshot(fmt.Sprintf("_output/%s.png", name))
}


func getVar(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}
	return val
}

func initVars() {
	serverHost = getVar("PGHOST", "postgres")
	serverPort = getVar("PGPORT", "5432")
	serverUser = getVar("PGUSER", "postgres")
	serverPassword = getVar("PGPASSWORD", "postgres")
	serverDatabase = getVar("PGDATABASE", "booktown")
}


func TestPgweb(t *testing.T) {
	RegisterFailHandler(Fail)
	initVars()
	RunSpecs(t, "Pgweb Suite")
}

var agoutiDriver *agouti.WebDriver

var _ = BeforeSuite(func() {
	agoutiDriver = agouti.ChromeDriver()
	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(agoutiDriver.Stop()).To(Succeed())
})

var page *agouti.Page

var _ = BeforeEach(func() {
	var err error

	page, err = agoutiDriver.NewPage()
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