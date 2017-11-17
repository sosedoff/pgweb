package api

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/data"
)

// Middleware function to check database connection status before running queries
func dbCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.Replace(c.Request.URL.Path, command.Opts.Prefix, "", -1)

		if allowedPaths[path] == true {
			c.Next()
			return
		}

		// We dont care about sessions unless they're enabled
		if !command.Opts.Sessions {
			if DbClient == nil {
				c.JSON(400, Error{"Not connected"})
				c.Abort()
				return
			}

			c.Next()
			return
		}

		sessionId := getSessionId(c.Request)
		if sessionId == "" {
			c.JSON(400, Error{"Session ID is required"})
			c.Abort()
			return
		}

		conn := DbSessions[sessionId]
		if conn == nil {
			c.JSON(400, Error{"Not connected"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Middleware function to print out request parameters and body for debugging
func requestInspectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseForm()
		log.Println("Request params:", err, c.Request.Form)
	}
}

func serveStaticAsset(path string, c *gin.Context) {
	data, err := data.Asset("static" + path)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.Data(200, assetContentType(path), data)
}

func serveResult(result interface{}, err error, c *gin.Context) {
	if err != nil {
		c.JSON(400, NewError(err))
		return
	}

	c.JSON(200, result)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Allow-Origin", command.Opts.CorsOrigin)
	}
}
