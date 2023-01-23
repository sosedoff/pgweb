package queries

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reMetaPrefix  = regexp.MustCompile(`^\s*--\s*pgweb:\s*(.+)`)
	reMetaContent = regexp.MustCompile(`([\w]+)\s*=\s*([^;]+)`)
	reMatchAll    = regexp.MustCompile(`^(.+)$`)
	reExpression  = regexp.MustCompile(`[\[\]\(\)\+\*]+`)

	allowedModes = map[string]bool{"readonly": true, "*": true}
)

type metadata struct {
	host     matcher
	user     matcher
	database matcher
	mode     matcher
}

func (m metadata) match(host, database, user, mode string) bool {
	// All input values should be set before we can match
	if host == "" || database == "" || user == "" || mode == "" {
		return false
	}

	return m.host.match(host) &&
		m.user.match(user) &&
		m.database.match(database) &&
		m.mode.match(mode)
}

func parseMetadata(input string) (*metadata, error) {
	matches := reMetaPrefix.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, nil
	}

	raw := strings.TrimSpace(matches[1])
	if len(raw) == 0 {
		return nil, nil
	}

	content := reMetaContent.FindAllStringSubmatch(raw, -1)
	if len(content) == 0 {
		return nil, nil
	}

	meta := &metadata{
		host:     stringMatcher{src: "localhost"},
		database: stringMatcher{src: "*", re: reMatchAll},
		user:     stringMatcher{src: "*", re: reMatchAll},
		mode:     valuesMatcher{src: "*", allowed: allowedModes},
	}

	foundKeys := map[string]bool{}

	for _, item := range content {
		var err error

		key := strings.TrimSpace(item[1])
		value := strings.TrimSpace(item[2])

		if foundKeys[key] {
			return nil, fmt.Errorf("duplicate key: %q", key)
		}
		foundKeys[key] = true

		switch key {
		case "host":
			meta.host, err = newStringMatcher(value)
		case "database":
			meta.database, err = newStringMatcher(value)
		case "user":
			meta.user, err = newStringMatcher(value)
		case "mode":
			if !allowedModes[value] {
				err = fmt.Errorf("invalid value for %q attribute: %q", "mode", value)
			}
			meta.mode = newValuesMatcher(value, allowedModes)
		default:
			err = fmt.Errorf("invalid meta attribute: %q", key)
		}

		if err != nil {
			return nil, err
		}
	}

	return meta, nil
}
