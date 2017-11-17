package agouti_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/internal/matchers"
)

var _ = Describe("Options", func() {
	Describe("#Browser", func() {
		It("should return an Option that sets a browser name", func() {
			config := NewTestConfig()
			Browser("some browser")(config)
			Expect(config.BrowserName).To(Equal("some browser"))
		})
	})

	Describe("#Timeout", func() {
		It("should return an Option that sets a timeout as a time.Duration", func() {
			config := NewTestConfig()
			Timeout(5)(config)
			Expect(config.Timeout).To(Equal(5 * time.Second))
		})
	})

	Describe("#Desired", func() {
		It("should return an Option that sets desired capabilities", func() {
			config := NewTestConfig()
			capabilities := NewCapabilities("some feature")
			Desired(capabilities)(config)
			Expect(config.DesiredCapabilities).To(Equal(capabilities))
		})
	})

	Describe("#RejectInvalidSSL", func() {
		It("should return an Option that rejects invalid SSL certificates", func() {
			config := NewTestConfig()
			Expect(config.RejectInvalidSSL).To(BeFalse())
			RejectInvalidSSL(config)
			Expect(config.RejectInvalidSSL).To(BeTrue())
		})
	})

	Describe("#Debug", func() {
		It("should return an Option that debugs a WebDriver", func() {
			config := NewTestConfig()
			Expect(config.Debug).To(BeFalse())
			Debug(config)
			Expect(config.Debug).To(BeTrue())
		})
	})

	Describe("#HTTPClient", func() {
		It("should return an Option that sets a *http.Client", func() {
			config := NewTestConfig()
			client := &http.Client{}
			HTTPClient(client)(config)
			Expect(config.HTTPClient).To(ExactlyEqual(client))
		})
	})

	Describe("#ChromeOptions", func() {
		It("should return an Option with ChromeOptions set", func() {
			config := NewTestConfig()
			ChromeOptions("args", []string{"v1", "v2"})(config)
			Expect(config.ChromeOptions["args"]).To(Equal([]string{"v1", "v2"}))
			ChromeOptions("other", "value")(config)
			Expect(config.ChromeOptions["args"]).To(Equal([]string{"v1", "v2"}))
			Expect(config.ChromeOptions["other"]).To(Equal("value"))
		})
	})

	Describe("#Merge", func() {
		It("should apply any provided options to an existing config", func() {
			config := NewTestConfig()
			Browser("some browser")(config)
			newConfig := config.Merge([]Option{Timeout(5), Debug})
			Expect(newConfig.BrowserName).To(Equal("some browser"))
			Expect(newConfig.Timeout).To(Equal(5 * time.Second))
			Expect(newConfig.Debug).To(BeTrue())
		})
	})

	Describe("#Capabilities", func() {
		It("should return a merged copy of the desired capabilities", func() {
			config := NewTestConfig()
			capabilities := NewCapabilities().Browser("some browser")
			Desired(capabilities)(config)
			Expect(config.Capabilities()["browserName"]).To(Equal("some browser"))
			Expect(config.Capabilities()["acceptSslCerts"]).To(BeTrue())
			Browser("some other browser")(config)
			RejectInvalidSSL(config)
			Expect(config.Capabilities()["browserName"]).To(Equal("some other browser"))
			Expect(config.Capabilities()["acceptSslCerts"]).To(BeFalse())
			ChromeOptions("args", "someArg")(config)
			Expect(config.Capabilities()["chromeOptions"]).To(
				Equal(map[string]interface{}{"args": "someArg"}),
			)
		})
	})
})
