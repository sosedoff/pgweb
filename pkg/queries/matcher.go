package queries

import (
	"fmt"
	"regexp"
)

type matcher interface {
	match(string) bool
	input() string
}

type stringMatcher struct {
	src string
	re  *regexp.Regexp
}

func (m stringMatcher) match(input string) bool {
	if m.re == nil {
		return m.src == input
	}
	return m.re.MatchString(input)
}

func (m stringMatcher) input() string {
	return m.src
}

func newStringMatcher(input string) (matcher, error) {
	var (
		re  *regexp.Regexp
		err error
	)

	if input == "*" { // just match everything
		re = reMatchAll
	} else if reExpression.MatchString(input) { // check if input is regex on its own
		re, err = regexp.Compile(fmt.Sprintf("^%s$", input))
		if err != nil {
			return nil, err
		}
	}

	return stringMatcher{src: input, re: re}, nil
}

type valuesMatcher struct {
	src     string
	allowed map[string]bool
}

func (m valuesMatcher) match(input string) bool {
	return m.allowed[input]
}

func (m valuesMatcher) input() string {
	return m.src
}

func newValuesMatcher(src string, allowed map[string]bool) matcher {
	return valuesMatcher{allowed: allowed, src: src}
}
