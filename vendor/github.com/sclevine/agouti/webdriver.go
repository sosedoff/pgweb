package agouti

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

// A WebDriver controls a WebDriver process. This struct embeds api.WebDriver,
// which provides Start and Stop methods for starting and stopping the process.
type WebDriver struct {
	*api.WebDriver
	defaultOptions *config
}

// NewWebDriver returns an instance of a WebDriver specified by
// a templated URL and command. The URL should be the location of the
// WebDriver Wire Protocol web service brought up by the command. The
// command should be provided as a list of arguments (each of which are
// templated).
//
// The Timeout Option specifies how many seconds to wait for the web service
// to become available. The default timeout is 5 seconds.
//
// The HTTPClient Option specifies a *http.Client to use for all WebDriver
// communications. The default client is http.DefaultClient.
//
// Any other provided Options are treated as default Options for new pages.
//
// Valid template parameters are:
//   {{.Host}} - local address to bind to (usually 127.0.0.1)
//   {{.Port}} - arbitrary free port on the local address
//   {{.Address}} - {{.Host}}:{{.Port}}
//
// Selenium JAR example:
//   command := []string{"java", "-jar", "selenium-server.jar", "-port", "{{.Port}}"}
//   agouti.NewWebDriver("http://{{.Address}}/wd/hub", command)
func NewWebDriver(url string, command []string, options ...Option) *WebDriver {
	apiWebDriver := api.NewWebDriver(url, command)
	defaultOptions := config{Timeout: apiWebDriver.Timeout}.Merge(options)
	apiWebDriver.Timeout = defaultOptions.Timeout
	apiWebDriver.Debug = defaultOptions.Debug
	apiWebDriver.HTTPClient = defaultOptions.HTTPClient
	return &WebDriver{apiWebDriver, defaultOptions}
}

// NewPage returns a *Page that corresponds to a new WebDriver session.
// Provided Options configure the page. For instance, to disable JavaScript:
//    capabilities := agouti.NewCapabilities().Without("javascriptEnabled")
//    driver.NewPage(agouti.Desired(capabilities))
// For Selenium, a Browser Option (or a Desired Option with Capabilities that
// specify a Browser) must be provided. For instance:
//    seleniumDriver.NewPage(agouti.Browser("safari"))
// Specific Options (such as Browser) have precedence over Capabilities
// specified by the Desired Option.
//
// The HTTPClient Option will be ignored if passed to this function. New pages
// will always use the *http.Client provided to their WebDriver, or
// http.DefaultClient if none was provided.
func (w *WebDriver) NewPage(options ...Option) (*Page, error) {
	newOptions := w.defaultOptions.Merge(options)
	session, err := w.Open(newOptions.Capabilities())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebDriver: %s", err)
	}

	return newPage(session), nil
}
