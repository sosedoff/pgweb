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

func TestDetectDumpVersion(t *testing.T) {
	examples := []struct {
		input   string
		match   bool
		version string
	}{
		{"", false, ""},
		{"pg_dump (PostgreSQL) 9.6", true, "9.6"},
		{"pg_dump 10", true, "10"},
		{"pg_dump (PostgreSQL) 14.5 (Homebrew)", true, "14.5"},
	}

	for _, ex := range examples {
		t.Run("input:"+ex.input, func(t *testing.T) {
			match, version := detectDumpVersion(ex.input)

			assert.Equal(t, ex.match, match)
			assert.Equal(t, ex.version, version)
		})
	}
}

func TestGetMajorMinorVersion(t *testing.T) {
	examples := []struct {
		input string
		major int
		minor int
	}{
		{"", 0, 0},
		{"   ", 0, 0},
		{"0", 0, 0},
		{"9.6", 9, 6},
		{"9.6.1.1", 9, 6},
		{"10", 10, 0},
		{"10.1 ", 10, 1},
	}

	for _, ex := range examples {
		t.Run(ex.input, func(t *testing.T) {
			major, minor := getMajorMinorVersion(ex.input)
			assert.Equal(t, ex.major, major)
			assert.Equal(t, ex.minor, minor)
		})
	}
}

func TestCheckVersionRequirement(t *testing.T) {
	examples := []struct {
		client string
		server string
		result bool
	}{
		{"", "", true},
		{"0", "0", true},
		{"9.6", "9.7", false},
		{"9.6.10", "9.6.25", true},
		{"10.0", "10.1", true},
		{"10.5", "10.1", true},
		{"14.5", "15.1", false},
	}

	for _, ex := range examples {
		assert.Equal(t, ex.result, checkVersionRequirement(ex.client, ex.server))
	}
}
