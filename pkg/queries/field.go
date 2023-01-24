package queries

import (
	"fmt"
	"regexp"
)

type field struct {
	value string
	re    *regexp.Regexp
}

func (f field) String() string {
	return f.value
}

func (f field) matches(input string) bool {
	if f.re != nil {
		return f.re.MatchString(input)
	}
	return f.value == input
}

func newField(value string) (field, error) {
	f := field{value: value}

	if value == "*" { // match everything
		f.re = reMatchAll
	} else if reExpression.MatchString(value) { // match by given expression
		re, err := regexp.Compile(fmt.Sprintf("^%s$", value))
		if err != nil {
			return f, err
		}
		f.re = re
	}

	return f, nil
}
