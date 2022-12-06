package api

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger

	reConnectToken = regexp.MustCompile("/connect/(.*)")
)

func init() {
	if logger == nil {
		logger = logrus.New()
	}
}

// TODO: Move this into server struct when it's ready
func SetLogger(l *logrus.Logger) {
	logger = l
}

func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	debug := logger.Level > logrus.InfoLevel

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		if !debug {
			// Skip static assets logging
			if strings.Contains(path, "/static/") {
				return
			}

			path = sanitizeLogPath(path)
		}

		status := c.Writer.Status()
		end := time.Now()
		latency := end.Sub(start)

		fields := logrus.Fields{
			"status":      status,
			"method":      c.Request.Method,
			"remote_addr": c.ClientIP(),
			"duration":    latency,
			"path":        path,
		}

		if err := c.Errors.Last(); err != nil {
			fields["error"] = err.Error()
		}

		// Additional fields for debugging
		if debug {
			fields["raw_query"] = c.Request.URL.RawQuery

			if c.Request.Method != http.MethodGet {
				fields["raw_form"] = c.Request.Form
			}
		}

		entry := logger.WithFields(fields)
		msg := "http_request"

		switch {
		case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
			entry.Warn(msg)
		case status >= http.StatusInternalServerError:
			entry.Error(msg)
		default:
			entry.Info(msg)
		}
	}
}

func sanitizeLogPath(str string) string {
	return reConnectToken.ReplaceAllString(str, "/connect/REDACTED")
}
