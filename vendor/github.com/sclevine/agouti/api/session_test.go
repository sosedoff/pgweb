package api_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/api/internal/mocks"
	. "github.com/sclevine/agouti/internal/matchers"
)

var _ = Describe("Session", func() {
	var (
		bus     *mocks.Bus
		session *Session
	)

	BeforeEach(func() {
		bus = &mocks.Bus{}
		session = &Session{bus}
	})

	Describe("#Delete", func() {
		It("should successfully send a DELETE to the / endpoint", func() {
			Expect(session.Delete()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("DELETE"))
			Expect(bus.SendCall.Endpoint).To(Equal(""))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Delete()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElement", func() {
		It("should successfully send a POST to the element endpoint", func() {
			_, err := session.GetElement(Selector{"css selector", "#selector"})
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("element"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		It("should successfully return an element with an ID and session", func() {
			bus.SendCall.Result = `{"ELEMENT": "some-id"}`
			element, err := session.GetElement(Selector{})
			Expect(err).NotTo(HaveOccurred())
			Expect(element.ID).To(Equal("some-id"))
			Expect(element.Session).To(ExactlyEqual(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetElement(Selector{})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetElements", func() {
		It("should successfully send a POST to the elements endpoint", func() {
			_, err := session.GetElements(Selector{"css selector", "#selector"})
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("elements"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"using": "css selector", "value": "#selector"}`))
		})

		It("should return a slice of elements with IDs and sessions", func() {
			bus.SendCall.Result = `[{"ELEMENT": "some-id"}, {"ELEMENT": "some-other-id"}]`
			elements, err := session.GetElements(Selector{"css selector", "#selector"})
			Expect(err).NotTo(HaveOccurred())
			Expect(elements[0].ID).To(Equal("some-id"))
			Expect(elements[0].Session).To(ExactlyEqual(session))
			Expect(elements[1].ID).To(Equal("some-other-id"))
			Expect(elements[1].Session).To(ExactlyEqual(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetElements(Selector{"css selector", "#selector"})
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetActiveElement", func() {
		It("should successfully send a POST to the element/active endpoint", func() {
			_, err := session.GetActiveElement()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("element/active"))
		})

		It("should return the active element with an ID and session", func() {
			bus.SendCall.Result = `{"ELEMENT": "some-id"}`
			element, err := session.GetActiveElement()
			Expect(err).NotTo(HaveOccurred())
			Expect(element.ID).To(Equal("some-id"))
			Expect(element.Session).To(ExactlyEqual(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetActiveElement()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetWindow", func() {
		It("should successfully send a GET to the window_handle endpoint", func() {
			_, err := session.GetWindow()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("window_handle"))
		})

		It("should return the current window with the retrieved ID and session", func() {
			bus.SendCall.Result = `"some-id"`
			window, err := session.GetWindow()
			Expect(err).NotTo(HaveOccurred())
			Expect(window.ID).To(Equal("some-id"))
			Expect(window.Session).To(ExactlyEqual(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetWindow()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetWindows", func() {
		It("should successfully send a GET to the window_handles endpoint", func() {
			_, err := session.GetWindows()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("window_handles"))
		})

		It("should return all windows with their retrieved IDs and sessions", func() {
			bus.SendCall.Result = `["some-id", "some-other-id"]`
			windows, err := session.GetWindows()
			Expect(err).NotTo(HaveOccurred())
			Expect(windows[0].ID).To(Equal("some-id"))
			Expect(windows[0].Session).To(ExactlyEqual(session))
			Expect(windows[1].ID).To(Equal("some-other-id"))
			Expect(windows[1].Session).To(ExactlyEqual(session))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetWindows()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetWindow", func() {
		It("should successfully send a POST to the window endpoint", func() {
			window := &Window{ID: "some-id"}
			Expect(session.SetWindow(window)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("window"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"name": "some-id"}`))
		})

		Context("when the window is nil", func() {
			It("should return an error", func() {
				Expect(session.SetWindow(nil)).To(MatchError("nil window is invalid"))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetWindow(&Window{})).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetWindowByName", func() {
		It("should successfully send a POST to the window endpoint", func() {
			Expect(session.SetWindowByName("some name")).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("window"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"name": "some name"}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetWindowByName("")).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteWindow", func() {
		It("should successfully send a DELETE to the window endpoint", func() {
			Expect(session.DeleteWindow()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("DELETE"))
			Expect(bus.SendCall.Endpoint).To(Equal("window"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				err := session.DeleteWindow()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetCookies", func() {
		It("should successfully send a GET to the cookie endpoint", func() {
			_, err := session.GetCookies()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("cookie"))
		})

		It("should return the cookies", func() {
			bus.SendCall.Result = `[{"name": "some-cookie"}, {"name": "some-other-cookie"}]`
			cookies, err := session.GetCookies()
			Expect(err).NotTo(HaveOccurred())
			Expect(cookies).To(Equal([]*Cookie{
				{Name: "some-cookie"},
				{Name: "some-other-cookie"},
			}))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetCookies()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetCookie", func() {
		It("should successfully send a POST to the cookie endpoint", func() {
			Expect(session.SetCookie(&Cookie{Name: "some-cookie"})).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("cookie"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"cookie": {"name": "some-cookie", "value": ""}}`))
		})

		Context("when the cookie is nil", func() {
			It("should return an error", func() {
				Expect(session.SetCookie(nil)).To(MatchError("nil cookie is invalid"))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetCookie(&Cookie{})).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookie", func() {
		It("should successfully send a DELETE to the cookie/some-cookie endpoint", func() {
			Expect(session.DeleteCookie("some-cookie")).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("DELETE"))
			Expect(bus.SendCall.Endpoint).To(Equal("cookie/some-cookie"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DeleteCookie("")).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteCookies", func() {
		It("should successfully send a DELETE to the cookie endpoint", func() {
			Expect(session.DeleteCookies()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("DELETE"))
			Expect(bus.SendCall.Endpoint).To(Equal("cookie"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DeleteCookies()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetScreenshot", func() {
		It("should successfully send a GET to the screenshot endpoint", func() {
			_, err := session.GetScreenshot()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("screenshot"))
		})

		Context("when the image is valid base64", func() {
			It("should return the decoded image", func() {
				bus.SendCall.Result = `"c29tZS1wbmc="`
				image, err := session.GetScreenshot()
				Expect(err).NotTo(HaveOccurred())
				Expect(string(image)).To(Equal("some-png"))
			})
		})

		Context("when the image is not valid base64", func() {
			It("should return an error", func() {
				bus.SendCall.Result = `"..."`
				_, err := session.GetScreenshot()
				Expect(err).To(MatchError("illegal base64 data at input byte 0"))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetScreenshot()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetURL", func() {
		It("should successfully send a GET to the url endpoint", func() {
			_, err := session.GetURL()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("url"))
		})

		It("should return the page URL", func() {
			bus.SendCall.Result = `"http://example.com"`
			url, err := session.GetURL()
			Expect(err).NotTo(HaveOccurred())
			Expect(url).To(Equal("http://example.com"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetURL()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetURL", func() {
		It("should successfully send a POST to the url endpoint", func() {
			Expect(session.SetURL("http://example.com")).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("url"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"url": "http://example.com"}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetURL("")).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetTitle", func() {
		It("should successfully send a GET to the title endpoint", func() {
			_, err := session.GetTitle()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("title"))
		})

		It("should return the page title", func() {
			bus.SendCall.Result = `"Some Title"`
			title, err := session.GetTitle()
			Expect(err).NotTo(HaveOccurred())
			Expect(title).To(Equal("Some Title"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetTitle()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetSource", func() {
		It("should successfully send a GET to the source endpoint", func() {
			_, err := session.GetSource()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("source"))
		})

		It("should return the page source", func() {
			bus.SendCall.Result = `"some source"`
			source, err := session.GetSource()
			Expect(err).NotTo(HaveOccurred())
			Expect(source).To(Equal("some source"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetSource()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#MoveTo", func() {
		It("should successfully send a POST to the moveto endpoint", func() {
			Expect(session.MoveTo(nil, nil)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("moveto"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON("{}"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				err := session.MoveTo(nil, nil)
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when an element is provided", func() {
			It("should encode the element into the request JSON", func() {
				element := &Element{ID: "some-id"}
				session.MoveTo(element, nil)
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"element": "some-id"}`))
			})
		})

		Context("when a X offset is provided", func() {
			It("should encode the element into the request JSON", func() {
				session.MoveTo(nil, XOffset(100))
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"xoffset": 100}`))
			})
		})

		Context("when a Y offset is provided", func() {
			It("should encode the element into the request JSON", func() {
				session.MoveTo(nil, YOffset(200))
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"yoffset": 200}`))
			})
		})

		Context("when an XY offset is provided", func() {
			It("should encode the element into the request JSON", func() {
				session.MoveTo(nil, XYOffset{X: 300, Y: 400})
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"xoffset": 300, "yoffset": 400}`))
			})
		})
	})

	Describe("#Frame", func() {
		It("should successfully send a POST to the frame endpoint", func() {
			Expect(session.Frame(&Element{ID: "some-id"})).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("frame"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"id": {"ELEMENT": "some-id"}}`))
		})

		Context("When the provided frame in nil", func() {
			It("should successfully send a POST to the frame endpoint", func() {
				Expect(session.Frame(nil)).To(Succeed())
				Expect(bus.SendCall.Method).To(Equal("POST"))
				Expect(bus.SendCall.Endpoint).To(Equal("frame"))
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"id": null}`))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Frame(nil)).To(MatchError("some error"))
			})
		})
	})

	Describe("#FrameParent", func() {
		It("should successfully send a POST to the frame/parent endpoint", func() {
			Expect(session.FrameParent()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("frame/parent"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.FrameParent()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Execute", func() {
		It("should successfully send a POST to the execute endpoint", func() {
			Expect(session.Execute("some javascript code", []interface{}{1, "two"}, nil)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("execute"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"script": "some javascript code", "args": [1, "two"]}`))
		})

		It("should fill the provided results interface", func() {
			var result struct{ Some string }
			bus.SendCall.Result = `{"some": "result"}`
			err := session.Execute("some javascript code", []interface{}{1, "two"}, &result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Some).To(Equal("result"))
		})

		Context("when called with nil arguments", func() {
			It("should send an empty list for args", func() {
				session.Execute("some javascript code", nil, nil)
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"script": "some javascript code", "args": []}`))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Execute("", nil, nil)).To(MatchError("some error"))
			})
		})
	})

	Describe("#Forward", func() {
		It("should successfully send a POST to the forward endpoint", func() {
			Expect(session.Forward()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("forward"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Forward()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Back", func() {
		It("should successfully send a POST to the back endpoint", func() {
			Expect(session.Back()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("back"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to go back in history", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Back()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Refresh", func() {
		It("should successfully send a POST to the refresh endpoint", func() {
			Expect(session.Refresh()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("refresh"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error indicating the bus failed to refresh the page", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Refresh()).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetAlertText", func() {
		It("should successfully send a GET to the alert_text endpoint", func() {
			_, err := session.GetAlertText()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("alert_text"))
		})

		It("should return the current alert text", func() {
			bus.SendCall.Result = `"some text"`
			text, err := session.GetAlertText()
			Expect(err).NotTo(HaveOccurred())
			Expect(text).To(Equal("some text"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetAlertText()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#SetAlertText", func() {
		It("should successfully send a POST to the alert_text endpoint", func() {
			Expect(session.SetAlertText("some text")).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("alert_text"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"text": "some text"}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.SetAlertText("some text")).To(MatchError("some error"))
			})
		})
	})

	Describe("#AcceptAlert", func() {
		It("should successfully send a POST to the accept_alert endpoint", func() {
			Expect(session.AcceptAlert()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("accept_alert"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.AcceptAlert()).To(MatchError("some error"))
			})
		})
	})

	Describe("#DismissAlert", func() {
		It("should successfully send a POST to the dismiss_alert endpoint", func() {
			Expect(session.DismissAlert()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("dismiss_alert"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DismissAlert()).To(MatchError("some error"))
			})
		})
	})

	Describe("#NewLogs", func() {
		It("should successfully send a POST to the log endpoint", func() {
			_, err := session.NewLogs("browser")
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("log"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"type": "browser"}`))
		})

		It("should return all logs", func() {
			bus.SendCall.Result = `[
				{"message": "some message", "level": "INFO", "timestamp": 1417988844498},
				{"message": "some other message", "level": "WARNING", "timestamp": 1417982864598}
			]`
			logs, err := session.NewLogs("browser")
			Expect(err).NotTo(HaveOccurred())
			Expect(logs[0].Message).To(Equal("some message"))
			Expect(logs[0].Level).To(Equal("INFO"))
			Expect(logs[0].Timestamp).To(BeEquivalentTo(1417988844498))
			Expect(logs[1].Message).To(Equal("some other message"))
			Expect(logs[1].Level).To(Equal("WARNING"))
			Expect(logs[1].Timestamp).To(BeEquivalentTo(1417982864598))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.NewLogs("browser")
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#GetLogTypes", func() {
		It("should successfully send a GET to the log/types endpoint", func() {
			_, err := session.GetLogTypes()
			Expect(err).NotTo(HaveOccurred())
			Expect(bus.SendCall.Method).To(Equal("GET"))
			Expect(bus.SendCall.Endpoint).To(Equal("log/types"))
		})

		It("should return the current alert text", func() {
			bus.SendCall.Result = `["first type", "second type"]`
			types, err := session.GetLogTypes()
			Expect(err).NotTo(HaveOccurred())
			Expect(types).To(Equal([]string{"first type", "second type"}))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				_, err := session.GetLogTypes()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#DoubleClick", func() {
		It("should successfully send a POST to the doubleclick endpoint", func() {
			Expect(session.DoubleClick()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("doubleclick"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DoubleClick()).To(MatchError("some error"))
			})
		})
	})

	Describe("#Click", func() {
		It("should successfully send a POST to the click endpoint", func() {
			Expect(session.Click(RightButton)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("click"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"button": 2}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Click(RightButton)).To(MatchError("some error"))
			})
		})
	})

	Describe("#ButtonDown", func() {
		It("should successfully send a POST to the buttondown endpoint", func() {
			Expect(session.ButtonDown(RightButton)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("buttondown"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"button": 2}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.ButtonDown(RightButton)).To(MatchError("some error"))
			})
		})
	})

	Describe("#ButtonUp", func() {
		It("should successfully send a POST to the buttonup endpoint", func() {
			Expect(session.ButtonUp(RightButton)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("buttonup"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"button": 2}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.ButtonUp(RightButton)).To(MatchError("some error"))
			})
		})
	})

	Describe("#TouchDown", func() {
		It("should successfully send a POST to the touch/down endpoint", func() {
			Expect(session.TouchDown(100, 200)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("touch/down"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"x": 100, "y": 200}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.TouchDown(100, 200)).To(MatchError("some error"))
			})
		})
	})

	Describe("#TouchUp", func() {
		It("should successfully send a POST to the touch/up endpoint", func() {
			Expect(session.TouchUp(100, 200)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("touch/up"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"x": 100, "y": 200}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.TouchUp(100, 200)).To(MatchError("some error"))
			})
		})
	})

	Describe("#TouchMove", func() {
		It("should successfully send a POST to the touch/move endpoint", func() {
			Expect(session.TouchMove(100, 200)).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("touch/move"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"x": 100, "y": 200}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.TouchMove(100, 200)).To(MatchError("some error"))
			})
		})
	})

	Describe("#TouchClick", func() {
		It("should successfully send a POST to the touch/click endpoint", func() {
			Expect(session.TouchClick(&Element{ID: "some-element-id"})).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("touch/click"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"element": "some-element-id"}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.TouchClick(&Element{ID: "some-element-id"})).To(MatchError("some error"))
			})
		})

		Context("when the element is nil", func() {
			It("shoul return an error", func() {
				Expect(session.TouchClick(nil)).To(MatchError("nil element is invalid"))
			})
		})
	})

	Describe("#TouchDoubleClick", func() {
		It("should successfully send a POST to the touch/doubleclick endpoint", func() {
			Expect(session.TouchDoubleClick(&Element{ID: "some-element-id"})).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("touch/doubleclick"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"element": "some-element-id"}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.TouchDoubleClick(&Element{ID: "some-element-id"})).To(MatchError("some error"))
			})
		})

		Context("when the element is nil", func() {
			It("shoul return an error", func() {
				Expect(session.TouchDoubleClick(nil)).To(MatchError("nil element is invalid"))
			})
		})
	})

	Describe("#TouchLongClick", func() {
		It("should successfully send a POST to the touch/longclick endpoint", func() {
			Expect(session.TouchLongClick(&Element{ID: "some-element-id"})).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("touch/longclick"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"element": "some-element-id"}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.TouchLongClick(&Element{ID: "some-element-id"})).To(MatchError("some error"))
			})
		})

		Context("when the element is nil", func() {
			It("shoul return an error", func() {
				Expect(session.TouchLongClick(nil)).To(MatchError("nil element is invalid"))
			})
		})
	})

	Describe("#TouchFlick", func() {
		Context("when provided with an offset and element", func() {
			var (
				element *Element
				offset  Offset
			)

			BeforeEach(func() {
				element = &Element{ID: "some-element-id"}
				offset = XYOffset{X: 100, Y: 200}
			})

			Context("when provided with a scalar speed", func() {
				It("should successfully send a POST to the touch/flick endpoint", func() {
					Expect(session.TouchFlick(element, offset, ScalarSpeed(300))).To(Succeed())
					Expect(bus.SendCall.Method).To(Equal("POST"))
					Expect(bus.SendCall.Endpoint).To(Equal("touch/flick"))
					Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{
						"element": "some-element-id",
						"xoffset": 100,
						"yoffset": 200,
						"speed": 300
					}`))
				})
			})

			Context("when provided with a vector speed", func() {
				It("should successfully send a POST to the touch/flick endpoint", func() {
					Expect(session.TouchFlick(element, offset, VectorSpeed{X: 300, Y: 400})).To(Succeed())
					Expect(bus.SendCall.Method).To(Equal("POST"))
					Expect(bus.SendCall.Endpoint).To(Equal("touch/flick"))
					Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{
						"element": "some-element-id",
						"xoffset": 100,
						"yoffset": 200,
						"speed": 500
					}`))
				})
			})
		})

		Context("when provided with no offset or element", func() {
			Context("when provided with a scalar speed", func() {
				It("should successfully send a POST to the touch/flick endpoint", func() {
					Expect(session.TouchFlick(nil, nil, ScalarSpeed(5))).To(Succeed())
					Expect(bus.SendCall.Method).To(Equal("POST"))
					Expect(bus.SendCall.Endpoint).To(Equal("touch/flick"))
					Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{
						"xspeed": 3,
						"yspeed": 3
					}`))
				})
			})

			Context("when provided with a vector speed", func() {
				It("should successfully send a POST to the touch/flick endpoint", func() {
					Expect(session.TouchFlick(nil, nil, VectorSpeed{X: 100, Y: 200})).To(Succeed())
					Expect(bus.SendCall.Method).To(Equal("POST"))
					Expect(bus.SendCall.Endpoint).To(Equal("touch/flick"))
					Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{
						"xspeed": 100,
						"yspeed": 200
					}`))
				})
			})
		})

		Context("when provided with an element but no offset", func() {
			It("should return an error", func() {
				Expect(session.TouchFlick(&Element{}, nil, ScalarSpeed(0))).To(MatchError("element must be provided if offset is provided and vice versa"))
			})
		})

		Context("when provided with an offset but no element", func() {
			It("should return an error", func() {
				Expect(session.TouchFlick(nil, XYOffset{}, ScalarSpeed(0))).To(MatchError("element must be provided if offset is provided and vice versa"))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.TouchFlick(nil, nil, ScalarSpeed(0))).To(MatchError("some error"))
			})
		})
	})

	Describe("#TouchScroll", func() {
		Context("when provided with an offset and element", func() {
			It("should successfully send a POST to the touch/scroll endpoint", func() {
				element := &Element{ID: "some-element-id"}
				offset := XYOffset{X: 100, Y: 200}
				Expect(session.TouchScroll(element, offset)).To(Succeed())
				Expect(bus.SendCall.Method).To(Equal("POST"))
				Expect(bus.SendCall.Endpoint).To(Equal("touch/scroll"))
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{
					"element": "some-element-id",
					"xoffset": 100,
					"yoffset": 200
				}`))
			})
		})

		Context("when provided with only an offset", func() {
			It("should successfully send a POST to the touch/scroll endpoint", func() {
				offset := XYOffset{X: 100, Y: 200}
				Expect(session.TouchScroll(nil, offset)).To(Succeed())
				Expect(bus.SendCall.Method).To(Equal("POST"))
				Expect(bus.SendCall.Endpoint).To(Equal("touch/scroll"))
				Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{
					"xoffset": 100,
					"yoffset": 200
				}`))
			})
		})

		Context("when provided with no offset", func() {
			It("should return an error", func() {
				Expect(session.TouchScroll(&Element{}, nil)).To(MatchError("nil offset is invalid"))
			})
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				offset := XYOffset{X: 100, Y: 200}
				Expect(session.TouchScroll(nil, offset)).To(MatchError("some error"))
			})
		})
	})

	Describe("#Keys", func() {
		It("should successfully send a POST request to the keys endpoint", func() {
			Expect(session.Keys("text")).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("POST"))
			Expect(bus.SendCall.Endpoint).To(Equal("keys"))
			Expect(bus.SendCall.BodyJSON).To(MatchJSON(`{"value": ["t", "e", "x", "t"]}`))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.Keys("text")).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteLocalStorage", func() {
		It("should successfully send a POST to the delete local storage endpoint", func() {
			Expect(session.DeleteLocalStorage()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("DELETE"))
			Expect(bus.SendCall.Endpoint).To(Equal("local_storage"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DeleteLocalStorage()).To(MatchError("some error"))
			})
		})
	})

	Describe("#DeleteSessionStorage", func() {
		It("should successfully send a POST to the delete session storage endpoint", func() {
			Expect(session.DeleteSessionStorage()).To(Succeed())
			Expect(bus.SendCall.Method).To(Equal("DELETE"))
			Expect(bus.SendCall.Endpoint).To(Equal("session_storage"))
		})

		Context("when the bus indicates a failure", func() {
			It("should return an error", func() {
				bus.SendCall.Err = errors.New("some error")
				Expect(session.DeleteSessionStorage()).To(MatchError("some error"))
			})
		})
	})
})
