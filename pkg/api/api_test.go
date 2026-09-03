package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/sosedoff/pgweb/pkg/command"
)

func Test_assetContentType(t *testing.T) {
	samples := map[string]string{
		"foo.html": "text/html; charset=utf-8",
		"foo.css":  "text/css; charset=utf-8",
		"foo.js":   "application/javascript",
		"foo.icon": "image-x-icon",
		"foo.png":  "image/png",
		"foo.jpg":  "image/jpeg",
		"foo.gif":  "image/gif",
		"foo.eot":  "application/vnd.ms-fontobject",
		"foo.svg":  "image/svg+xml",
		"foo.foo":  "text/plain; charset=utf-8",
		"foo":      "text/plain; charset=utf-8",
	}

	alternatives := map[string]string{
		"foo.js": "text/javascript; charset=utf-8",
	}

	for name, expected := range samples {
		if alternatives[name] == "" {
			assert.Equal(t, expected, assetContentType(name))
			continue
		}

		actual := assetContentType(name)

		if actual != expected && actual != alternatives[name] {
			t.Errorf("expected %v but got %v (alternative value failed)", expected, actual)
		}
	}
}

func Test_Connect_RejectsDirectUrlWhenConnectBackendConfigured(t *testing.T) {
	originalBackend := command.Opts.ConnectBackend
	command.Opts.ConnectBackend = "http://backend.example"
	defer func() { command.Opts.ConnectBackend = originalBackend }()

	server := gin.Default()
	server.POST("/api/connect", Connect)

	form := url.Values{}
	form.Set("url", "postgres://attacker:pw@internal-host:5432/secret_db?sslmode=disable")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/connect", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	server.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"Not permitted","status":400}`, w.Body.String())
}
