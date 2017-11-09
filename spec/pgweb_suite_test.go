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

	// https://github.com/onsi/ginkgo/issues/285#issuecomment-290575636

	// The Ginkgo test runner takes over os.Args and fills it with
	// its own flags.  This makes the cobra command arg parsing
	// fail because of unexpected options.  Work around this.

	// Save a copy of os.Args
	//origArgs := os.Args[:]

	// Trim os.Args to only the first arg
	//os.Args = os.Args[:1] // trim to only the first arg, which is the command itself

	// Run the command which parses os.Args
	//pwcli.Run()

	// Restore os.Args
	//os.Args = origArgs[:]

	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(agoutiDriver.Stop()).To(Succeed())
})
