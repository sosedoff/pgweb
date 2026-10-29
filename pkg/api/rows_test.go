package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
)

func withOptions(opts command.Options, fn func()) {
	original := command.Opts
	command.Opts = opts
	defer func() { command.Opts = original }()

	fn()
}

func Test_rowEditingEnabled(t *testing.T) {
	examples := []struct {
		name     string
		opts     command.Options
		expected bool
	}{
		{"default", command.Options{}, true},
		{"disabled", command.Options{DisableRowEditing: true}, false},
		{"read only", command.Options{ReadOnly: true}, false},
		{"both", command.Options{DisableRowEditing: true, ReadOnly: true}, false},
	}

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			withOptions(ex.opts, func() {
				assert.Equal(t, ex.expected, rowEditingEnabled())
			})
		})
	}
}

func Test_requireRowEditing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	examples := []struct {
		name         string
		opts         command.Options
		expectedCode int
		expectedBody string
	}{
		{"enabled", command.Options{}, http.StatusOK, "ok"},
		{"disabled", command.Options{DisableRowEditing: true}, http.StatusForbidden, `{"error":"Row editing is disabled","status":403}`},
		{"read only", command.Options{ReadOnly: true}, http.StatusForbidden, `{"error":"Row editing is disabled","status":403}`},
	}

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			withOptions(ex.opts, func() {
				router := gin.New()
				router.POST("/rows", requireRowEditing(), func(c *gin.Context) {
					c.String(http.StatusOK, "ok")
				})

				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, httptest.NewRequest("POST", "/rows", nil))

				assert.Equal(t, ex.expectedCode, resp.Code)
				assert.Equal(t, ex.expectedBody, resp.Body.String())
			})
		})
	}
}

func Test_parseRowRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid payload", func(t *testing.T) {
		body := `{"key":{"id":1},"values":{"name":"foo","note":null}}`

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("PUT", "/rows", bytes.NewBufferString(body))

		req, err := parseRowRequest(c)
		require.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"id": float64(1)}, req.Key)
		assert.Equal(t, map[string]interface{}{"name": "foo", "note": nil}, req.Values)
	})

	t.Run("invalid payload", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("PUT", "/rows", bytes.NewBufferString("not json"))

		_, err := parseRowRequest(c)
		assert.Equal(t, errInvalidRowPayload, err)
	})

	t.Run("empty payload", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("DELETE", "/rows", bytes.NewBufferString(""))

		_, err := parseRowRequest(c)
		assert.Equal(t, errInvalidRowPayload, err)
	})
}

func Test_disableWrites(t *testing.T) {
	meta := &client.TableMetadata{
		Schema:     "public",
		Table:      "books",
		Insertable: true,
		Editable:   true,
	}

	result := disableWrites(meta)

	assert.False(t, result.Insertable)
	assert.False(t, result.Editable)

	// The original metadata is cached and must not be mutated
	assert.True(t, meta.Insertable)
	assert.True(t, meta.Editable)
}

func Test_rowResponse(t *testing.T) {
	t.Run("with rows", func(t *testing.T) {
		result := &client.Result{
			Columns: []string{"id", "name"},
			Rows:    []client.Row{{int64(1), "foo"}},
		}

		response := rowResponse(result)

		assert.Equal(t, []string{"id", "name"}, response["columns"])
		assert.Equal(t, map[string]interface{}{"id": int64(1), "name": "foo"}, response["row"])
	})

	t.Run("without rows", func(t *testing.T) {
		result := &client.Result{Columns: []string{"id"}, Rows: []client.Row{}}
		response := rowResponse(result)

		assert.Equal(t, []string{"id"}, response["columns"])
		assert.NotContains(t, response, "row")
	})
}

func Test_parseFilters(t *testing.T) {
	examples := []struct {
		name     string
		input    string
		expected []client.Filter
		err      bool
	}{
		{"empty", "", nil, false},
		{
			"single filter",
			`[{"column":"status","op":"equal","value":"ACTIVE"}]`,
			[]client.Filter{{Column: "status", Op: "equal", Value: "ACTIVE"}},
			false,
		},
		{
			"multiple filters",
			`[{"column":"status","op":"equal","value":"ACTIVE"},{"column":"note","op":"not_null"}]`,
			[]client.Filter{
				{Column: "status", Op: "equal", Value: "ACTIVE"},
				{Column: "note", Op: "not_null"},
			},
			false,
		},
		{"invalid json", `[{"column":}]`, nil, true},
		{"not an array", `{"column":"status"}`, nil, true},
	}

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			filters, err := parseFilters(ex.input)

			if ex.err {
				assert.Equal(t, errInvalidFilters, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, ex.expected, filters)
		})
	}
}
