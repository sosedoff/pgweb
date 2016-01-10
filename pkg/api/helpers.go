package api

import (
	"fmt"
	"mime"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

var extraMimeTypes = map[string]string{
	".icon": "image-x-icon",
	".ttf":  "application/x-font-ttf",
	".woff": "application/x-font-woff",
	".eot":  "application/vnd.ms-fontobject",
	".svg":  "image/svg+xml",
	".html": "text/html; charset-utf-8",
}

// Paths that dont require database connection
var allowedPaths = map[string]bool{
	"/api/sessions":  true,
	"/api/info":      true,
	"/api/connect":   true,
	"/api/bookmarks": true,
	"/api/history":   true,
}

type Error struct {
	Message string `json:"error"`
}

func getSessionId(c *gin.Context) string {
	id := c.Request.Header.Get("x-session-id")
	if id == "" {
		id = c.Request.URL.Query().Get("_session_id")
	}
	return id
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
