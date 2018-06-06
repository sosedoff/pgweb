package client

import (
	"strings"
)

// Get short version from the string
// Example: 10.2.3.1 -> 10.2
func getMajorMinorVersion(str string) string {
	chunks := strings.Split(str, ".")
	if len(chunks) == 0 {
		return str
	}
	return strings.Join(chunks[0:2], ".")
}
