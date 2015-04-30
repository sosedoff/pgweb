package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		"foo.ttf":  "application/x-font-ttf",
		"foo.woff": "application/x-font-woff",
		"foo.foo":  "text/plain; charset=utf-8",
		"foo":      "text/plain; charset=utf-8",
	}

	for name, expected := range samples {
		assert.Equal(t, expected, assetContentType(name))
	}
}
