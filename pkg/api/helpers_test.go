package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
