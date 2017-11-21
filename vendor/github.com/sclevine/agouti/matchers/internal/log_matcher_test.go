package internal_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers/internal"
	"github.com/sclevine/agouti/matchers/internal/mocks"
)

var _ = Describe("LogMatcher", func() {
	var (
		matcher *LogMatcher
		page    *mocks.Page
	)

	BeforeEach(func() {
		page = &mocks.Page{}
		matcher = &LogMatcher{
			ExpectedMessages: []string{"some log", "some other log"},
			Levels:           []string{"FIRST", "SECOND", "THIRD"},
			Name:             "name",
			Type:             "some type",
		}
	})

	Describe("#Match", func() {
		Context("when the actual object is a page", func() {
			It("should request all of the logs with the provided type", func() {
				matcher.Match(page)
				Expect(page.ReadAllLogsCall.LogType).To(Equal("some type"))
			})

			Context("when all expected logs have been logged with an expected level", func() {
				It("should successfully return true", func() {
					page.ReadAllLogsCall.ReturnLogs = []agouti.Log{
						{"some log", "", "SECOND", time.Time{}},
						{"some other log", "", "FIRST", time.Time{}},
						{"another log", "", "OTHER", time.Time{}},
					}
					Expect(matcher.Match(page)).To(BeTrue())
				})
			})

			Context("when not all expected logs have been logged", func() {
				It("should successfully return false", func() {
					page.ReadAllLogsCall.ReturnLogs = []agouti.Log{
						{"some log", "", "FIRST", time.Time{}},
						{"another log", "", "SECOND", time.Time{}},
						{"yet another log", "", "THIRD", time.Time{}},
					}
					Expect(matcher.Match(page)).To(BeFalse())
				})
			})

			Context("when not all expected logs have been logged with expected levels", func() {
				It("should successfully return false", func() {
					page.ReadAllLogsCall.ReturnLogs = []agouti.Log{
						{"some log", "", "FIRST", time.Time{}},
						{"some other log", "", "OTHER", time.Time{}},
						{"another log", "", "SECOND", time.Time{}},
					}
					Expect(matcher.Match(page)).To(BeFalse())
				})
			})

			Context("when any log is expected", func() {
				BeforeEach(func() {
					matcher.ExpectedMessages = []string{}
				})

				Context("when any log of an expected type is logged", func() {
					It("should successfully return true", func() {
						page.ReadAllLogsCall.ReturnLogs = []agouti.Log{
							{"first log", "", "OTHER", time.Time{}},
							{"second log", "", "SECOND", time.Time{}},
						}
						Expect(matcher.Match(page)).To(BeTrue())
					})
				})

				Context("when no logs of an expected type are logged", func() {
					It("should successfully return false", func() {
						page.ReadAllLogsCall.ReturnLogs = []agouti.Log{
							{"first log", "", "OTHER", time.Time{}},
							{"second log", "", "ANOTHER", time.Time{}},
						}
						Expect(matcher.Match(page)).To(BeFalse())
					})
				})
			})

			Context("when retrieving the logs fails", func() {
				It("should return an error", func() {
					page.ReadAllLogsCall.Err = errors.New("some error")
					_, err := matcher.Match(page)
					Expect(err).To(MatchError("some error"))
				})
			})
		})

		Context("when the actual object is not a page", func() {
			It("should return an error", func() {
				_, err := matcher.Match("not a page")
				Expect(err).To(MatchError("HaveLoggedName matcher requires a Page.  Got:\n    <string>: not a page"))
			})
		})
	})

	Describe("#FailureMessage", func() {
		It("should return a failure message with logs if logs are present", func() {
			message := matcher.FailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page to have name logs matching\n    some log\n    some other log"))
		})

		It("should return a failure message without logs if logs are not present", func() {
			matcher.ExpectedMessages = []string{}
			message := matcher.FailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page to have logged name logs"))
		})
	})

	Describe("#NegatedFailureMessage", func() {
		It("should return a negated failure message", func() {
			message := matcher.NegatedFailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page not to have name logs matching\n    some log\n    some other log"))
		})

		It("should return a negated failure message without logs if logs are present", func() {
			matcher.ExpectedMessages = []string{}
			message := matcher.NegatedFailureMessage(page)
			Expect(message).To(ContainSubstring("Expected page not to have logged name logs"))
		})
	})
})
