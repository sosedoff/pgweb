package client

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess(t *testing.T) {
	t.Run("large numbers", func(t *testing.T) {
		result := Result{
			Columns: []string{"value"},
			Rows: []Row{
				{int(1234)},
				{int64(9223372036854775807)},
				{int64(-9223372036854775808)},
				{float64(9223372036854775808.9223372036854775808)},
				{float64(999999999999999.9)},
			},
		}

		result.PostProcess()

		assert.Equal(t, 1234, result.Rows[0][0])
		assert.Equal(t, "9223372036854775807", result.Rows[1][0])
		assert.Equal(t, "-9223372036854775808", result.Rows[2][0])
		assert.Equal(t, "9.223372036854776e+18", result.Rows[3][0])
		assert.Equal(t, "9.999999999999999e+14", result.Rows[4][0])
	})

	t.Run("binary encoding", func(t *testing.T) {
		result := Result{
			Columns: []string{"data"},
			Rows: []Row{
				{"text value"},
				{"text with symbols !@#$%"},
				{string([]byte{10, 11, 12, 13})},
			},
		}

		result.PostProcess()

		assert.Equal(t, "text value", result.Rows[0][0])
		assert.Equal(t, "text with symbols !@#$%", result.Rows[1][0])
		assert.Equal(t, "CgsMDQ==", result.Rows[2][0])
	})
}

func TestCSV(t *testing.T) {
	result := Result{
		Columns: []string{"id", "name", "email", "extra"},
		Rows: []Row{
			{1, "John", "john@example.com", "data"},
			{2, "Bob", "bob@example.com", nil},
		},
	}

	expected := strings.Join([]string{
		"id,name,email,extra",
		"1,John,john@example.com,data",
		"2,Bob,bob@example.com,",
	}, "\n") + "\n"

	assert.Equal(t, expected, string(result.CSV()))
}

func TestJSON(t *testing.T) {
	result := Result{
		Columns: []string{"id", "name", "email"},
		Rows: []Row{
			{1, "John", "john@example.com"},
			{2, "Bob", "bob@example.com"},
		},
	}

	obj := []struct {
		Id    int
		Name  string
		Email string
	}{}

	expected := []struct {
		Id    int
		Name  string
		Email string
	}{
		{1, "John", "john@example.com"},
		{2, "Bob", "bob@example.com"},
	}

	assert.NoError(t, json.Unmarshal(result.JSON(), &obj))
	assert.Equal(t, 2, len(obj))
	assert.Equal(t, expected, obj)

	t.Run("invalid time", func(t *testing.T) {
		loc, err := time.LoadLocation("UTC")
		if err != nil {
			panic(err)
		}

		command.Opts.DisablePrettyJSON = true
		defer func() {
			command.Opts.DisablePrettyJSON = false
		}()

		result := Result{
			Columns: []string{"value"},
			Rows: []Row{
				{time.Unix(1640995200, 0).In(loc)},
				{time.Unix(222539616000, 0).In(loc)},
				{time.Unix(254096611200, 0).In(loc)},
			},
		}

		result.PostProcess()
		assert.Equal(t, `[{"value":"2022-01-01T00:00:00Z"},{"value":"9022-01-01T00:00:00Z"},{"value":"ERR: INVALID_DATE"}]`, string(result.JSON()))
	})
}

func TestResultFormat(t *testing.T) {
	result := Result{
		Columns: []string{"col1", "col2", "col3", "col4"},
		Rows: []Row{
			{"1", "2", "3", nil},
			{"4", "5", "6", nil},
		},
	}

	expected := []map[string]interface{}{
		{"col1": "1", "col2": "2", "col3": "3", "col4": nil},
		{"col1": "4", "col2": "5", "col3": "6", "col4": nil},
	}

	assert.Equal(t, expected, result.Format())
}
