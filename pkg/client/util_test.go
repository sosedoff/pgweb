package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectServerType(t *testing.T) {
	examples := []struct {
		input      string
		match      bool
		serverType string
		version    string
	}{
		{input: "",
			match:      false,
			serverType: "",
			version:    "",
		},
		{
			input:      " postgresql 15 ",
			match:      true,
			serverType: postgresType,
			version:    "15",
		},
		{
			input:      "PostgreSQL 14.5 (Homebrew) on aarch64-apple-darwin21.6.0",
			match:      true,
			serverType: postgresType,
			version:    "14.5",
		},
		{
			input:      "PostgreSQL 11.16, compiled by Visual C++ build 1800, 64-bit",
			match:      true,
			serverType: postgresType,
			version:    "11.16",
		},
	}

	for _, ex := range examples {
		t.Run("input:"+ex.input, func(t *testing.T) {
			match, stype, version := detectServerTypeAndVersion(ex.input)

			assert.Equal(t, ex.match, match)
			assert.Equal(t, ex.serverType, stype)
			assert.Equal(t, ex.version, version)
		})
	}
}
