package client

import (
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
)

var (
	// List of keywords that are not allowed in read-only mode
	reRestrictedKeywords = regexp.MustCompile(`(?mi)\s?(CREATE|INSERT|DROP|DELETE|TRUNCATE|GRANT|OPEN|IMPORT|COPY)\s`)

	// Comment regular expressions
	reSlashComment = regexp.MustCompile(`(?m)/\*.+\*/`)
	reDashComment  = regexp.MustCompile(`(?m)--.+`)

	// Postgres version signature
	postgresSignature     = regexp.MustCompile(`(?i)postgresql ([\d\.]+)\s?`)
	postgresDumpSignature = regexp.MustCompile(`\s([\d\.]+)\s?`)
	postgresType          = "PostgreSQL"

	// Cockroach version signature
	cockroachSignature = regexp.MustCompile(`(?i)cockroachdb ccl v([\d\.]+)\s?`)
	cockroachType      = "CockroachDB"
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

func detectServerTypeAndVersion(version string) (bool, string, string) {
	version = strings.TrimSpace(version)

	// Detect postgresql
	matches := postgresSignature.FindAllStringSubmatch(version, 1)
	if len(matches) > 0 {
		return true, postgresType, matches[0][1]
	}

	// Detect cockroachdb
	matches = cockroachSignature.FindAllStringSubmatch(version, 1)
	if len(matches) > 0 {
		return true, cockroachType, matches[0][1]
	}

	return false, "", ""
}

// detectDumpVersion parses out version from `pg_dump -V` command.
func detectDumpVersion(version string) (bool, string) {
	matches := postgresDumpSignature.FindAllStringSubmatch(version, 1)
	if len(matches) > 0 {
		return true, matches[0][1]
	}
	return false, ""
}

func checkVersionRequirement(cur string, min string) (bool, error) {
	current, err := version.NewVersion(cur)
	if err != nil {
		return false, err
	}

	minimum, err := version.NewVersion(min)
	if err != nil {
		return false, err
	}

	return current.GreaterThanOrEqual(minimum), nil
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
