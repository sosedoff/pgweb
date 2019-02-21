package client

import (
	"regexp"
	"strings"
)

// List of keywords that are not allowed in read-only mode
var restrictedKeywords = regexp.MustCompile(`(?mi)\s?(CREATE|INSERT|DROP|DELETE|TRUNCATE|GRANT|OPEN|IMPORT|COPY)\s`)

// Get short version from the string
// Example: 10.2.3.1 -> 10.2
func getMajorMinorVersion(str string) string {
	chunks := strings.Split(str, ".")
	if len(chunks) == 0 {
		return str
	}
	return strings.Join(chunks[0:2], ".")
}

// containsRestrictedKeywords returns true if given keyword is not allowed in read-only mode
func containsRestrictedKeywords(str string) bool {
	return restrictedKeywords.MatchString(str)
}
