package service_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api/internal/service"
)

var _ = Describe("Service", func() {
	var service *Service

	BeforeEach(func() {
		service = &Service{
			URLTemplate: "some-url",
			CmdTemplate: []string{"true"},
		}
	})

	Describe("#URL", func() {
		Context("when the server is not running", func() {
			It("should return an empty string", func() {
				Expect(service.URL()).To(BeEmpty())
				Expect(service.Start(false)).To(Succeed())
				Expect(service.Stop()).To(Succeed())
				Expect(service.URL()).To(BeEmpty())
			})
		})

		Context("when the server is running", func() {
			It("should successfully return the URL", func() {
				defer service.Stop()
				Expect(service.Start(false)).To(Succeed())
				Expect(service.URL()).To(MatchRegexp("some-url"))
			})
		})
	})

	Describe("#Start", func() {
		Context("when the service is started multiple times", func() {
			It("should return an error indicating that service is already running", func() {
				defer service.Stop()
				Expect(service.Start(false)).To(Succeed())
				Expect(service.Start(false)).To(MatchError("already running"))
			})
		})

		Context("when the binary is not available in PATH", func() {
			It("should return an error indicating the binary needs to be installed", func() {
				service.CmdTemplate = []string{"not-in-path"}
				Expect(service.Start(false)).To(MatchError("failed to run command: exec: \"not-in-path\": executable file not found in $PATH"))
			})
		})

		Context("when the service is started in debug mode", func() {
			It("should successfully start", func() {
				defer service.Stop()
				Expect(service.Start(true)).To(Succeed())
			})
		})

		Describe("the provided templated URL", func() {
			Context("when the template is invalid", func() {
				It("should return an error", func() {
					defer service.Stop()
					service.URLTemplate = "{{}}"
					Expect(service.Start(false)).To(MatchError("failed to parse URL: template: URL:1: missing value for command"))
				})
			})

			Context("when the template does not match the provided parameters", func() {
				It("should return an error", func() {
					defer service.Stop()
					service.URLTemplate = "{{.Bad}}"
					Expect(service.Start(false).Error()).To(MatchRegexp(`(failed to parse URL: template: URL:1:2: executing ){1}......(at <.Bad>: can't evaluate field Bad in type service.addressInfo){1}|(failed to parse URL: template: URL:1:2: executing ){1}......(at <.Bad>: Bad is not a field of struct type service.addressInfo){1}`))
				})
			})

			Context("when the template is valid", func() {
				It("should store a templated URL", func() {
					defer service.Stop()
					service.URLTemplate += "/status?test&{{.Address}}&{{.Host}}:{{.Port}}"
					service.Start(false)
					Expect(service.URL()).To(MatchRegexp(`test&127\.0\.0\.1:\d+&127\.0\.0\.1:\d+`))
				})
			})
		})

		Describe("the provided templated command", func() {
			Context("when the template is invalid", func() {
				It("should return an error", func() {
					defer service.Stop()
					service.CmdTemplate = []string{"correct", "{{}}"}
					Expect(service.Start(false)).To(MatchError("failed to parse command: template: command:1: missing value for command"))
				})
			})

			Context("when the template does not match the provided parameters", func() {
				It("should return an error", func() {
					defer service.Stop()
					service.CmdTemplate = []string{"correct", "{{.Bad}}"}
					Expect(service.Start(false).Error()).To(MatchRegexp(`(failed to parse command: template: command:1:2: executing ){1}..........(at <.Bad>: can't evaluate field Bad in type service.addressInfo){1}|(failed to parse command: template: command:1:2: executing ){1}..........(at <.Bad>: Bad is not a field of struct type service.addressInfo){1}`))
				})
			})

			Context("when the template is empty", func() {
				It("should return an error", func() {
					defer service.Stop()
					service.CmdTemplate = []string{}
					Expect(service.Start(false)).To(MatchError("failed to parse command: empty command"))
				})
			})

			Context("when the template is valid", func() {
				It("should not return an error", func() {
					defer service.Stop()
					service.CmdTemplate = []string{"true", "{{.Address}}{{.Host}}{{.Port}}"}
					Expect(service.Start(false)).To(Succeed())
				})
			})
		})
	})

	Describe("#Stop", func() {
		It("should stop a running server", func() {
			defer service.Stop()
			Expect(service.Start(false)).To(Succeed())
			Expect(service.Stop()).To(Succeed())
			Expect(service.Start(false)).To(Succeed())
		})

		Context("when the command is not started", func() {
			It("should return an error", func() {
				err := service.Stop()
				Expect(err).To(MatchError("already stopped"))
			})
		})
	})

	Describe("#WaitForBoot", func() {
		var (
			started bool
			server  *httptest.Server
		)

		BeforeEach(func() {
			started = false

			server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				if started && request.URL.Path == "/status" {
					response.WriteHeader(200)
				} else {
					response.WriteHeader(400)
				}
			}))

			service.URLTemplate = server.URL
		})

		AfterEach(func() {
			server.Close()
		})

		Context("when the service does not start before the provided timeout", func() {
			It("should return an error", func() {
				defer service.Stop()
				go func() {
					time.Sleep(3000 * time.Millisecond)
					started = true
				}()
				Expect(service.Start(false)).To(Succeed())
				Expect(service.WaitForBoot(1500 * time.Millisecond)).To(MatchError("failed to start before timeout"))
			})
		})

		Context("when the service starts before the provided timeout", func() {
			It("should not return an error", func() {
				defer service.Stop()
				go func() {
					time.Sleep(200 * time.Millisecond)
					started = true
				}()
				Expect(service.Start(false)).To(Succeed())
				Expect(service.WaitForBoot(1500 * time.Millisecond)).To(Succeed())
			})
		})
	})
})
