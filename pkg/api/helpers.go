package api

import (
	"fmt"
	"log"
	"mime"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sosedoff/pgweb/pkg/data"
)

var extraMimeTypes = map[string]string{
	".icon": "image-x-icon",
	".ttf":  "application/x-font-ttf",
	".woff": "application/x-font-woff",
	".eot":  "application/vnd.ms-fontobject",
	".svg":  "image/svg+xml",
	".html": "text/html; charset-utf-8",
}

type Error struct {
	Message string `json:"error"`
}

func getQueryParam(c *gin.Context, name string) string {
	result := ""
	q := c.Request.URL.Query()

	if len(q[name]) > 0 {
		result = q[name][0]
	}

	return result
}

func parseIntFormValue(c *gin.Context, name string, defValue int) (int, error) {
	val := c.Request.FormValue(name)

	if val == "" {
		return defValue, nil
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return defValue, fmt.Errorf("%s must be a number", name)
	}

	if num < 1 && defValue != 0 {
		return defValue, fmt.Errorf("%s must be greated than 0", name)
	}

	return num, nil
}

func assetContentType(name string) string {
	ext := filepath.Ext(name)
	result := mime.TypeByExtension(ext)

	if result == "" {
		result = extraMimeTypes[ext]
	}

	if result == "" {
		result = "text/plain; charset=utf-8"
	}

	return result
}

func NewError(err error) Error {
	return Error{err.Error()}
}

// Middleware function to check database connection status before running queries
func dbCheckMiddleware() gin.HandlerFunc {
	allowedPaths := []string{
		"/api/info",
		"/api/connect",
		"/api/bookmarks",
		"/api/history",
	}

	return func(c *gin.Context) {
		if DbClient != nil {
			c.Next()
			return
		}

		currentPath := c.Request.URL.Path
		allowed := false

		for _, path := range allowedPaths {
			if path == currentPath {
				allowed = true
				break
			}
		}

		if allowed {
			c.Next()
		} else {
			c.JSON(400, Error{"Not connected"})
			c.Abort()
		}

		return
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
