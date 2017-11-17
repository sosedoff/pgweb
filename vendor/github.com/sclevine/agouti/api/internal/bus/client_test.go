package bus_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api/internal/bus"
)

var _ = Describe("Session", func() {
	var (
		requestPath        string
		requestMethod      string
		requestBody        string
		requestContentType string
		responseBody       string
		responseStatus     int
		server             *httptest.Server
	)

	BeforeEach(func() {
		responseBody, responseStatus = "", 200
		requestPath, requestMethod, requestBody, requestContentType = "", "", "", ""
		server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			requestPath = request.URL.Path
			requestMethod = request.Method
			requestBodyBytes, _ := ioutil.ReadAll(request.Body)
			requestBody = string(requestBodyBytes)
			requestContentType = request.Header.Get("Content-Type")
			response.WriteHeader(responseStatus)
			response.Write([]byte(responseBody))
		}))
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("#Send", func() {
		var client *Client

		BeforeEach(func() {
			client = &Client{
				SessionURL: server.URL + "/session/some-id",
				HTTPClient: http.DefaultClient,
			}
		})

		It("should make a request with the method and full session endpoint", func() {
			client.Send("GET", "some/endpoint", nil, nil)
			Expect(requestPath).To(Equal("/session/some-id/some/endpoint"))
			Expect(requestMethod).To(Equal("GET"))
		})

		It("should use the provided HTTP client", func() {
			var path string
			client.HTTPClient = &http.Client{Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
				path = request.URL.Path
				return nil, errors.New("some error")
			})}
			err := client.Send("GET", "some/endpoint", nil, nil)
			Expect(err).To(MatchError(ContainSubstring("some error")))
			Expect(path).To(Equal("/session/some-id/some/endpoint"))
		})

		Context("with a valid request body", func() {
			It("should make a request with the provided body and application/json content type", func() {
				body := struct{ SomeValue string }{"some request value"}
				Expect(client.Send("POST", "some/endpoint", body, nil)).To(Succeed())
				Expect(requestBody).To(Equal(`{"SomeValue":"some request value"}`))
				Expect(requestContentType).To(Equal("application/json"))
			})
		})

		Context("with an invalid request body", func() {
			It("should return an invalid request body error", func() {
				err := client.Send("POST", "some/endpoint", func() {}, nil)
				Expect(err).To(MatchError("invalid request body: json: unsupported type: func()"))
			})
		})

		Context("when the provided body is nil", func() {
			It("should make a request without a body", func() {
				Expect(client.Send("POST", "some/endpoint", nil, nil)).To(Succeed())
				Expect(requestBody).To(BeEmpty())
			})
		})

		Context("when the session endpoint is empty", func() {
			It("should make a request to the session itself", func() {
				Expect(client.Send("GET", "", nil, nil)).To(Succeed())
				Expect(requestPath).To(Equal("/session/some-id"))
			})
		})

		Context("with an invalid URL", func() {
			It("should return an invalid request error", func() {
				client.SessionURL = "%@#$%"
				err := client.Send("GET", "some/endpoint", nil, nil)
				Expect(err).To(MatchError(`invalid request: parse %@: invalid URL escape "%@"`))
			})
		})

		Context("when the request fails entirely", func() {
			It("should return an error indicating that the request failed", func() {
				server.Close()
				err := client.Send("GET", "some/endpoint", nil, nil)
				Expect(err.Error()).To(MatchRegexp("request failed: .+ connection refused"))
			})
		})

		Context("when the server responds with a non-2xx status code", func() {
			BeforeEach(func() {
				responseStatus = 400
			})

			Context("when the server has a valid error message", func() {
				It("should return an error from the server indicating that the request failed", func() {
					responseBody = `{"value": {"message": "{\"errorMessage\": \"some error\"}"}}`
					err := client.Send("GET", "some/endpoint", nil, nil)
					Expect(err).To(MatchError("request unsuccessful: some error"))
				})
			})

			Context("when the server does not have a valid message", func() {
				It("should return an error indicating that the request failed with no details", func() {
					responseBody = `$$$`
					err := client.Send("GET", "some/endpoint", nil, nil)
					Expect(err).To(MatchError("request unsuccessful: $$$"))
				})
			})

			Context("when the server does not have a valid JSON-encoded error message", func() {
				It("should return an error with the entire message output", func() {
					responseBody = `{"value": {"message": "$$$"}}`
					err := client.Send("GET", "some/endpoint", nil, nil)
					Expect(err).To(MatchError("request unsuccessful: $$$"))
				})
			})
		})

		Context("when the request succeeds", func() {
			var result struct{ Some string }

			BeforeEach(func() {
				responseBody = `{"value": {"some": "response value"}}`
			})

			Context("with a valid response body", func() {
				It("should successfully unmarshal the returned JSON into the result", func() {
					Expect(client.Send("GET", "some/endpoint", nil, &result)).To(Succeed())
					Expect(result.Some).To(Equal("response value"))
				})
			})

			Context("with a response body value that cannot be read", func() {
				It("should return a failed to extract value from response error", func() {
					responseBody = "some unexpected response"
					err := client.Send("GET", "some/endpoint", nil, &result)
					Expect(err).To(MatchError("unexpected response: some unexpected response"))
				})
			})
		})
	})
})
