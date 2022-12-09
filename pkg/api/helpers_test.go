package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_desanitize64(t *testing.T) {
	examples := map[string]string{
		"test":        "test",
		"test+test+":  "test-test-",
		"test/test/":  "test_test_",
		"test=test==": "test.test..",
	}

	for expected, example := range examples {
		assert.Equal(t, expected, desanitize64(example))
	}
}

func Test_cleanQuery(t *testing.T) {
	assert.Equal(t, "a\nb\nc", cleanQuery("a\nb\nc"))
	assert.Equal(t, "", cleanQuery("--something"))
	assert.Equal(t, "test", cleanQuery("--test\ntest\n   -- test\n"))
}

func Test_sanitizeFilename(t *testing.T) {
	examples := map[string]string{
		"foo":              "foo",
		"fooBar":           "fooBar",
		"foo.bar":          "foo_bar",
		`"foo"."bar"`:      "foo_bar",
		"!@#$foo.&&*(&bar": "foo_bar",
	}

	for given, expected := range examples {
		t.Run(given, func(t *testing.T) {
			assert.Equal(t, expected, sanitizeFilename(given))
		})
	}
}

func Test_getSessionId(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	req.Header.Add("x-session-id", "token")
	assert.Equal(t, "token", getSessionId(req))

	req = &http.Request{}
	req.URL, _ = url.Parse("http://foobar/?_session_id=token")
	assert.Equal(t, "token", getSessionId(req))
}

func Test_serveResult(t *testing.T) {
	server := gin.Default()
	server.GET("/good", func(c *gin.Context) {
		serveResult(c, gin.H{"foo": "bar"}, nil)
	})
	server.GET("/bad", func(c *gin.Context) {
		serveResult(c, nil, errors.New("message"))
	})
	server.GET("/nodata", func(c *gin.Context) {
		serveResult(c, nil, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/good", nil)
	server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"foo":"bar"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/bad", nil)
	server.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"message","status":400}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/nodata", nil)
	server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `null`, w.Body.String())
}
