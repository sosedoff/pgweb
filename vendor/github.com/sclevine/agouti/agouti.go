// Package agouti is a universal WebDriver client for Go.
// It extends the agouti/api package to provide a feature-rich interface for
// controlling a web browser.
package agouti

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// PhantomJS returns an instance of a PhantomJS WebDriver.
//
// Provided Options will apply as default arguments for new pages.
// New pages will accept invalid SSL certificates by default. This
// may be disabled using the RejectInvalidSSL Option.
//
// The RejectInvalidSSL Option must be provided to the PhantomJS function
// (and not the NewPage method) for this Option to take effect on any
// PhantomJS page.
func PhantomJS(options ...Option) *WebDriver {
	command := []string{"phantomjs", "--webdriver={{.Address}}"}
	defaultOptions := config{}.Merge(options)
	if !defaultOptions.RejectInvalidSSL {
		command = append(command, "--ignore-ssl-errors=true")
	}
	return NewWebDriver("http://{{.Address}}", command, options...)
}

// ChromeDriver returns an instance of a ChromeDriver WebDriver.
//
// Provided Options will apply as default arguments for new pages.
// New pages will accept invalid SSL certificates by default. This
// may be disabled using the RejectInvalidSSL Option.
func ChromeDriver(options ...Option) *WebDriver {
	var binaryName string
	if runtime.GOOS == "windows" {
		binaryName = "chromedriver.exe"
	} else {
		binaryName = "chromedriver"
	}
	command := []string{binaryName, "--port={{.Port}}"}
	return NewWebDriver("http://{{.Address}}", command, options...)
}

// EdgeDriver returns an instance of a EdgeDriver WebDriver.
//
// Provided Options will apply as default arguments for new pages.
// New pages will accept invalid SSL certificates by default. This
// may be disabled using the RejectInvalidSSL Option.
func EdgeDriver(options ...Option) *WebDriver {
	var binaryName string
	if runtime.GOOS == "windows" {
		binaryName = "MicrosoftWebDriver.exe"
	} else {
		return nil
	}
	command := []string{binaryName, "--port={{.Port}}"}
	// Using {{.Address}} means using 127.0.0.1
	// But MicrosoftWebDriver only supports localhost, not 127.0.0.1
	return NewWebDriver("http://localhost:{{.Port}}", command, options...)
}

// Selenium returns an instance of a Selenium WebDriver.
//
// Provided Options will apply as default arguments for new pages.
// New pages will accept invalid SSL certificates by default. This
// may be disabled using the RejectInvalidSSL Option.
func Selenium(options ...Option) *WebDriver {
	command := []string{"selenium-server", "-port", "{{.Port}}"}
	return NewWebDriver("http://{{.Address}}/wd/hub", command, options...)
}

// Selendroid returns an instance of a Selendroid WebDriver.
//
// Provided Options will apply as default arguments for new pages.
// New pages will accept invalid SSL certificates by default. This
// may be disabled using the RejectInvalidSSL Option.
//
// The jarFile is a relative or absolute path to Selendroid JAR file.
// Selendroid will return nil if an invalid path is provided.
func Selendroid(jarFile string, options ...Option) *WebDriver {
	absJARPath, err := filepath.Abs(jarFile)
	if err != nil {
		return nil
	}

	command := []string{
		"java",
		"-jar", absJARPath,
		"-port", "{{.Port}}",
	}
	options = append([]Option{Timeout(90), Browser("android")}, options...)
	return NewWebDriver("http://{{.Address}}/wd/hub", command, options...)
}

// SauceLabs opens a Sauce Labs session and returns a *Page. Does not support Sauce Connect.
//
// This method takes the same Options as NewPage. Passing the Desired Option will
// completely override the provided name, platform, browser, and version.
func SauceLabs(name, platform, browser, version, username, accessKey string, options ...Option) (*Page, error) {
	url := fmt.Sprintf("http://%s:%s@ondemand.saucelabs.com/wd/hub", username, accessKey)
	capabilities := NewCapabilities().Browser(browser).Platform(platform).Version(version)
	capabilities["name"] = name
	return NewPage(url, append([]Option{Desired(capabilities)}, options...)...)
}
