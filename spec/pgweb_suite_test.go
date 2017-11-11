package spec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"

	"testing"
)

func TestPgweb(t *testing.T) {
	RegisterFailHandler(Fail)
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
