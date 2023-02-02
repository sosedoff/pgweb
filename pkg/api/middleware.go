package api

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/command"
)

// Middleware to check database connection status before running queries
func dbCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.Replace(c.Request.URL.Path, command.Opts.Prefix, "", -1)

		// Allow whitelisted paths
		if allowedPaths[path] {
			c.Next()
			return
		}

		// Check if session exists in single-session mode
		if !command.Opts.Sessions {
			if DbClient == nil {
				badRequest(c, errNotConnected)
				return
			}

			c.Next()
			return
		}

		// Determine session ID from the client request
		sid := getSessionId(c.Request)
		if sid == "" {
			badRequest(c, errSessionRequired)
			return
		}

		// Determine the database connection handle for the session
		conn := DbSessions.Get(sid)
		if conn == nil {
			badRequest(c, errNotConnected)
			return
		}

		c.Next()
	}
}

// Middleware to inject CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Allow-Origin", command.Opts.CorsOrigin)
	}
}

func requireLocalQueries() gin.HandlerFunc {
	return func(c *gin.Context) {
		if QueryStore == nil {
			badRequest(c, "local queries are disabled")
			return
		}

		c.Next()
	}
}
