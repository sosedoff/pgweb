package api

func NewTestWebDriver(service driverService) *WebDriver {
	return &WebDriver{service: service}
}
