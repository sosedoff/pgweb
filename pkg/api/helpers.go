package api

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/shared"
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

// List of characters replaced by javascript code to make queries url-safe.
var base64subs = map[string]string{
	"-": "+",
	"_": "/",
	".": "=",
}

type Error struct {
	Message string `json:"error"`
}

func NewError(err error) Error {
	return Error{err.Error()}
}

// Returns a clean query without any comment statements
func cleanQuery(query string) string {
	lines := []string{}

	for _, line := range strings.Split(query, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "--") {
			continue
		}
		lines = append(lines, line)
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func desanitize64(query string) string {
	// Before feeding the string into decoded, we must "reconstruct" the base64 data.
	// Javascript replaces a few characters to be url-safe.
	for olds, news := range base64subs {
		query = strings.Replace(query, olds, news, -1)
	}

	return query
}

func getSessionId(req *http.Request) string {
	id := req.Header.Get("x-session-id")
	if id == "" {
		id = req.URL.Query().Get("_session_id")
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
		return defValue, fmt.Errorf("%s must be greater than 0", name)
	}

	return num, nil
}

func parseSshInfo(c *gin.Context) *shared.SSHInfo {
	info := shared.SSHInfo{
		Host:     c.Request.FormValue("ssh_host"),
		Port:     c.Request.FormValue("ssh_port"),
		User:     c.Request.FormValue("ssh_user"),
		Password: c.Request.FormValue("ssh_password"),
		Key:      c.Request.FormValue("ssh_key"),
	}

	if info.Port == "" {
		info.Port = "22"
	}

	return &info
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
