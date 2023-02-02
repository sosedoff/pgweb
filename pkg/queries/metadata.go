package queries

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reMetaPrefix  = regexp.MustCompile(`(?m)^\s*--\s*pgweb:\s*(.+)`)
	reMetaContent = regexp.MustCompile(`([\w]+)\s*=\s*"([^"]+)"`)
	reMatchAll    = regexp.MustCompile(`^(.+)$`)
	reExpression  = regexp.MustCompile(`[\[\]\(\)\+\*]+`)

	allowedKeys  = []string{"title", "description", "host", "user", "database", "mode", "timeout"}
	allowedModes = map[string]bool{"readonly": true, "*": true}
)

type Metadata struct {
	Title       string
	Description string
	Host        field
	User        field
	Database    field
	Mode        field
	Timeout     *time.Duration
}

func parseMetadata(input string) (*Metadata, error) {
	fields, err := parseFields(input)
	if err != nil {
		return nil, err
	}
	if fields == nil {
		return nil, nil
	}

	// Host must be set to limit queries availability
	if fields["host"] == "" {
		return nil, fmt.Errorf("host field must be set")
	}

	// Allow matching for any user, database and mode by default
	if fields["user"] == "" {
		fields["user"] = "*"
	}
	if fields["database"] == "" {
		fields["database"] = "*"
	}
	if fields["mode"] == "" {
		fields["mode"] = "*"
	}

	hostField, err := newField(fields["host"])
	if err != nil {
		return nil, fmt.Errorf(`error initializing "host" field: %w`, err)
	}

	userField, err := newField(fields["user"])
	if err != nil {
		return nil, fmt.Errorf(`error initializing "user" field: %w`, err)
	}

	dbField, err := newField(fields["database"])
	if err != nil {
		return nil, fmt.Errorf(`error initializing "database" field: %w`, err)
	}

	if !allowedModes[fields["mode"]] {
		return nil, fmt.Errorf(`invalid "mode" field value: %q`, fields["mode"])
	}
	modeField, err := newField(fields["mode"])
	if err != nil {
		return nil, fmt.Errorf(`error initializing "mode" field: %w`, err)
	}

	var timeout *time.Duration
	if fields["timeout"] != "" {
		timeoutSec, err := strconv.Atoi(fields["timeout"])
		if err != nil {
			return nil, fmt.Errorf(`error initializing "timeout" field: %w`, err)
		}
		timeoutVal := time.Duration(timeoutSec) * time.Second
		timeout = &timeoutVal
	}

	return &Metadata{
		Title:       fields["title"],
		Description: fields["description"],
		Host:        hostField,
		User:        userField,
		Database:    dbField,
		Mode:        modeField,
		Timeout:     timeout,
	}, nil
}

func parseFields(input string) (map[string]string, error) {
	result := map[string]string{}
	seenKeys := map[string]bool{}

	allowed := map[string]bool{}
	for _, key := range allowedKeys {
		allowed[key] = true
	}

	matches := reMetaPrefix.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	for _, match := range matches {
		content := reMetaContent.FindAllStringSubmatch(match[1], -1)
		if len(content) == 0 {
			continue
		}

		for _, field := range content {
			key := field[1]
			value := field[2]

			if !allowed[key] {
				return result, fmt.Errorf("unknown key: %q", key)
			}
			if seenKeys[key] {
				return result, fmt.Errorf("duplicate key: %q", key)
			}

			seenKeys[key] = true
			result[key] = value
		}
	}

	return result, nil
}

func sanitizeMetadata(input string) string {
	lines := []string{}
	for _, line := range strings.Split(input, "\n") {
		line = reMetaPrefix.ReplaceAllString(line, "")
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}
