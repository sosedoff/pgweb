package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PrepareBigints(t *testing.T) {
	result := Result{
		Columns: []string{"value"},
		Rows: []Row{
			Row{int(1234)},
			Row{int64(9223372036854775807)},
			Row{int64(-9223372036854775808)},
			Row{float64(9223372036854775808.9223372036854775808)},
			Row{float64(999999999999999.9)},
		},
	}

	result.PrepareBigints()

	assert.Equal(t, 1234, result.Rows[0][0])
	assert.Equal(t, "9223372036854775807", result.Rows[1][0])
	assert.Equal(t, "-9223372036854775808", result.Rows[2][0])
	assert.Equal(t, "9.223372036854776e+18", result.Rows[3][0])
	assert.Equal(t, "9.999999999999999e+14", result.Rows[4][0])
}

func Test_CSV(t *testing.T) {
	result := Result{
		Columns: []string{"id", "name", "email"},
		Rows: []Row{
			Row{1, "John", "john@example.com"},
			Row{2, "Bob", "bob@example.com"},
		},
	}

	expected := "id,name,email\n1,John,john@example.com\n2,Bob,bob@example.com\n"
	output := string(result.CSV())

	assert.Equal(t, expected, output)
}

func Test_JSON(t *testing.T) {
	result := Result{
		Columns: []string{"id", "name", "email"},
		Rows: []Row{
			Row{1, "John", "john@example.com"},
			Row{2, "Bob", "bob@example.com"},
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
}
