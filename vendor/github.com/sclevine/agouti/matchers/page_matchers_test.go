package matchers_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("Page Matchers", func() {
	var page *mocks.Page

	BeforeEach(func() {
		page = &mocks.Page{}
	})

	Describe("#HaveTitle", func() {
		It("should return a ValueMatcher with the 'Title' method", func() {
			page.TitleCall.ReturnTitle = "Some Title"
			Expect(page).To(HaveTitle("Some Title"))
			Expect(page).NotTo(HaveTitle("Some Other Title"))
		})

		It("should set the matcher property to 'title'", func() {
			Expect(HaveTitle("").FailureMessage(nil)).To(ContainSubstring("to have title"))
		})
	})

	Describe("#HaveURL", func() {
		It("should return a ValueMatcher with the 'URL' method", func() {
			page.URLCall.ReturnURL = "some/url"
			Expect(page).To(HaveURL("some/url"))
			Expect(page).NotTo(HaveURL("some/other/url"))
		})

		It("should set the matcher property to 'URL'", func() {
			Expect(HaveURL("").FailureMessage(nil)).To(ContainSubstring("to have URL"))
		})
	})

	Describe("#HavePopupText", func() {
		It("should return a ValueMatcher with the 'PopupText' method", func() {
			page.PopupTextCall.ReturnText = "some text"
			Expect(page).To(HavePopupText("some text"))
			Expect(page).NotTo(HavePopupText("some other text"))
		})

		It("should set the matcher property to 'popup text'", func() {
			Expect(HavePopupText("").FailureMessage(nil)).To(ContainSubstring("to have popup text"))
		})
	})

	Describe("#HaveWindowCount", func() {
		It("should return a ValueMatcher with the 'WindowCount' method", func() {
			page.WindowCountCall.ReturnCount = 1
			Expect(page).To(HaveWindowCount(1))
			Expect(page).NotTo(HaveWindowCount(2))
		})

		It("should set the matcher property to 'window count'", func() {
			Expect(HaveWindowCount(0).FailureMessage(nil)).To(ContainSubstring("to have window count"))
		})
	})

	Describe("#HaveLoggedError", func() {
		It("should return a LogMatcher matcher for SEVERE and WARNING browser logs", func() {
			page.ReadAllLogsCall.ReturnLogs = []agouti.Log{
				{"some log", "", "SEVERE", time.Time{}},
				{"some other log", "", "WARNING", time.Time{}},
				{"another log", "", "INFO", time.Time{}},
			}
			Expect(page).To(HaveLoggedError("some log"))
			Expect(page).To(HaveLoggedError("some other log"))
			Expect(page).NotTo(HaveLoggedError("another log"))
			Expect(page.ReadAllLogsCall.LogType).To(Equal("browser"))
		})

		It("should set the log name to 'error'", func() {
			Expect(HaveLoggedError().FailureMessage(nil)).To(ContainSubstring("to have logged error"))
		})
	})

	Describe("#HaveLoggedInfo", func() {
		It("should return a LogMatcher matcher for INFO browser logs", func() {
			page.ReadAllLogsCall.ReturnLogs = []agouti.Log{
				{"some log", "", "SEVERE", time.Time{}},
				{"some other log", "", "WARNING", time.Time{}},
				{"another log", "", "INFO", time.Time{}},
			}
			Expect(page).NotTo(HaveLoggedInfo("some log"))
			Expect(page).NotTo(HaveLoggedInfo("some other log"))
			Expect(page).To(HaveLoggedInfo("another log"))
			Expect(page.ReadAllLogsCall.LogType).To(Equal("browser"))
		})

		It("should set the log name to 'info'", func() {
			Expect(HaveLoggedInfo().FailureMessage(nil)).To(ContainSubstring("to have logged info"))
		})
	})
})
