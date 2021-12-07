package client

import (
	"regexp"
	"strings"
)

var (
	// List of keywords that are not allowed in read-only mode
	reRestrictedKeywords = regexp.MustCompile(`(?mi)\s?(CREATE|INSERT|DROP|DELETE|TRUNCATE|GRANT|OPEN|IMPORT|COPY)\s`)

	// Comment regular expressions
	reSlashComment = regexp.MustCompile(`(?m)/\*.+\*/`)
	reDashComment  = regexp.MustCompile(`(?m)--.+`)
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

// containsRestrictedKeywords returns true if given keyword is not allowed in read-only mode
func containsRestrictedKeywords(str string) bool {
	str = reSlashComment.ReplaceAllString(str, "")
	str = reDashComment.ReplaceAllString(str, "")

	return reRestrictedKeywords.MatchString(str)
}

func hasBinary(data string, checkLen int) bool {
	for idx, chr := range data {
		if int(chr) < 32 || int(chr) > 126 {
			return true
		}
		if idx >= checkLen {
			break
		}
	}
	return false
}
