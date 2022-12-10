package client

import (
	"fmt"
	"regexp"
	"strings"
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

// Get major and minor version components
// Example: 10.2.3.1 -> 10.2
func getMajorMinorVersion(str string) (major int, minor int) {
	chunks := strings.Split(str, ".")
	fmt.Sscanf(chunks[0], "%d", &major)
	if len(chunks) > 1 {
		fmt.Sscanf(chunks[1], "%d", &minor)
	}
	return
}

// Get short version from the string
// Example: 10.2.3.1 -> 10.2
func getMajorMinorVersionString(str string) string {
	major, minor := getMajorMinorVersion(str)
	return fmt.Sprintf("%d.%d", major, minor)
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

func checkVersionRequirement(client, server string) bool {
	clientMajor, clientMinor := getMajorMinorVersion(client)
	serverMajor, serverMinor := getMajorMinorVersion(server)

	if serverMajor < 10 {
		return clientMajor >= serverMajor && clientMinor >= serverMinor
	}

	return clientMajor >= serverMajor
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
