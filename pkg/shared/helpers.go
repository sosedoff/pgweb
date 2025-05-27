package shared

import "regexp"

// Obfuscate credentials in URLs of the form //user:password@
func SanitizeConnectionString(str string) string {
	// This regex matches //user:password@ and replaces password with ***
	re := regexp.MustCompile(`(//[^:/@]+:)[^@]+(@)`) // matches //user:password@
	return re.ReplaceAllString(str, `$1***$2`)
}
