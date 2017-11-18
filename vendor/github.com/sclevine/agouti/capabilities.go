package agouti

import "encoding/json"

// A Capabilities instance defines the desired capabilities the WebDriver
// should use to configure a Page.
//
// For example, to open Firefox with JavaScript disabled:
//    capabilities := agouti.NewCapabilities().Browser("firefox").Without("javascriptEnabled")
//    driver.NewPage(agouti.Desired(capabilities))
// See: https://github.com/SeleniumHQ/selenium/wiki/DesiredCapabilities
//
// All methods called on this instance will modify the original instance.
type Capabilities map[string]interface{}

// NewCapabilities returns a Capabilities instance with any provided features enabled.
func NewCapabilities(features ...string) Capabilities {
	c := Capabilities{}
	for _, feature := range features {
		c.With(feature)
	}
	return c
}

// Browser sets the desired browser name.
// Possible values:
//    {android|chrome|firefox|htmlunit|internet explorer|iPhone|iPad|opera|safari}
func (c Capabilities) Browser(name string) Capabilities {
	c["browserName"] = name
	return c
}

// A ProxyConfig instance defines the desired proxy configuration the WebDriver
// should use to proxy a Page.
//
//    ProxyType is required and defines the type of proxy being used
//    Possible Values:
//        {direct|manual|pac|autodetect|system}
//
// See: https://github.com/SeleniumHQ/selenium/wiki/DesiredCapabilities#proxy-json-object
type ProxyConfig struct {
	ProxyType          string `json:"proxyType"`
	ProxyAutoconfigURL string `json:"proxyAutoconfigUrl,omitempty"`
	FTPProxy           string `json:"ftpProxy,omitempty"`
	HTTPProxy          string `json:"httpProxy,omitempty"`
	SSLProxy           string `json:"sslProxy,omitempty"`
	SOCKSProxy         string `json:"socksProxy,omitempty"`
	SOCKSUsername      string `json:"socksUsername,omitempty"`
	SOCKSPassword      string `json:"socksPassword,omitempty"`
	NoProxy            string `json:"noProxy,omitempty"`
}

// Proxy sets the desired proxy configuration.
func (c Capabilities) Proxy(p ProxyConfig) Capabilities {
	c["proxy"] = p
	return c
}

// Version sets the desired browser version (ex. "3.6").
func (c Capabilities) Version(version string) Capabilities {
	c["version"] = version
	return c
}

// Platform sets the desired browser platform.
// Possible values:
//    {WINDOWS|XP|VISTA|MAC|LINUX|UNIX|ANDROID|ANY}.
func (c Capabilities) Platform(platform string) Capabilities {
	c["platform"] = platform
	return c
}

// With enables the provided feature (ex. "trustAllSSLCertificates").
func (c Capabilities) With(feature string) Capabilities {
	c[feature] = true
	return c
}

// Without disables the provided feature (ex. "javascriptEnabled").
func (c Capabilities) Without(feature string) Capabilities {
	c[feature] = false
	return c
}

// JSON returns a JSON string representing the desired capabilities.
func (c Capabilities) JSON() (string, error) {
	capabilitiesJSON, err := json.Marshal(c)
	return string(capabilitiesJSON), err
}
