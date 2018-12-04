package api

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/command"
)

// Middleware to check database connection status before running queries
func dbCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.Replace(c.Request.URL.Path, command.Opts.Prefix, "", -1)

		// White-list provider paths
		if strings.HasPrefix(path, "/api/discovery") {
			c.Next()
			return
		}

		// Allow whitelisted paths
		if allowedPaths[path] {
			c.Next()
			return
		}

		// Check if session exists in single-session mode
		if !command.Opts.Sessions {
			if DbClient == nil {
				badRequest(c, "Not connected")
				return
			}

			c.Next()
			return
		}

		// Determine session ID from the client request
		sessionId := getSessionId(c.Request)
		if sessionId == "" {
			badRequest(c, "Session ID is required")
			return
		}

		// Determine the database connection handle for the session
		conn := DbSessions[sessionId]
		if conn == nil {
			badRequest(c, "Not connected")
			return
		}

		c.Next()
	}
}

// Middleware to print out request parameters and body for debugging
func requestInspectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseForm()
		log.Println("Request params:", err, c.Request.Form)
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
