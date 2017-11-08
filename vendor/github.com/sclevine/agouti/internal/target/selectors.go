package target

import "strings"

type Selectors []Selector

func (s Selectors) canMergeType(selectorType Type) bool {
	if len(s) == 0 {
		return false
	}
	last := s[len(s)-1]
	bothCSS := selectorType == CSS && last.Type == CSS
	return bothCSS && !last.Indexed && !last.Single
}

func (s Selectors) Append(selectorType Type, value string) Selectors {
	selector := Selector{Type: selectorType, Value: value}

	if s.canMergeType(selectorType) {
		lastIndex := len(s) - 1
		selector.Value = s[lastIndex].Value + " " + selector.Value
		return s[:lastIndex].append(selector)
	}
	return s.append(selector)
}

func (s Selectors) Single() Selectors {
	lastIndex := len(s) - 1
	if lastIndex < 0 {
		return nil
	}

	selector := s[lastIndex]
	selector.Single = true
	selector.Indexed = false
	return s[:lastIndex].append(selector)
}

func (s Selectors) At(index int) Selectors {
	lastIndex := len(s) - 1
	if lastIndex < 0 {
		return nil
	}

	selector := s[lastIndex]
	selector.Single = false
	selector.Indexed = true
	selector.Index = index
	return s[:lastIndex].append(selector)
}

func (s Selectors) String() string {
	var tags []string

	for _, selector := range s {
		tags = append(tags, selector.String())
	}

	return strings.Join(tags, " | ")
}

func (s Selectors) append(selector Selector) Selectors {
	selectorsCopy := append(Selectors(nil), s...)
	return append(selectorsCopy, selector)
}
