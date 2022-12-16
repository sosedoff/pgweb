package api

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_getRequestID(t *testing.T) {
	examples := []struct {
		headers map[string]string
		result  string
	}{
		{map[string]string{}, ""},
		{map[string]string{"X-Request-ID": "foo"}, "foo"},
		{map[string]string{"x-request-id": "foo"}, "foo"},
		{map[string]string{"x-request-id": "foo"}, "foo"},
		{map[string]string{"x-request-id": "foo", "x-amzn-trace-id": "amz"}, "foo"},
	}

	for _, ex := range examples {
		req := &http.Request{Header: http.Header{}}
		for k, v := range ex.headers {
			req.Header.Set(k, v)
		}

		assert.Equal(t, ex.result, getRequestID(&gin.Context{Request: req}))
	}
}
