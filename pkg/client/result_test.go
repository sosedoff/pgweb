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
		},
	}

	result.PrepareBigints()

	assert.Equal(t, 1234, result.Rows[0][0])
	assert.Equal(t, "9223372036854775807", result.Rows[1][0])
	assert.Equal(t, "-9223372036854775808", result.Rows[2][0])
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

	output := result.JSON()
	obj := []map[string]interface{}{}
	err := json.Unmarshal(output, &obj)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(obj))

	for i, row := range obj {
		for j, col := range result.Columns {
			assert.Equal(t, result.Rows[i][j], row[col])
		}
	}
}
