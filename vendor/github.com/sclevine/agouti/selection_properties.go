package agouti

import (
	"fmt"

	"github.com/sclevine/agouti/internal/element"
)

// Text returns the entirety of the text content for exactly one element.
func (s *Selection) Text() (string, error) {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return "", fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	text, err := selectedElement.GetText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve text for %s: %s", s, err)
	}
	return text, nil
}

// Active returns true if the single element that the selection refers to is active.
func (s *Selection) Active() (bool, error) {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	activeElement, err := s.session.GetActiveElement()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve active element: %s", err)
	}

	equal, err := selectedElement.IsEqualTo(activeElement)
	if err != nil {
		return false, fmt.Errorf("failed to compare selection to active element: %s", err)
	}

	return equal, nil
}

type propertyMethod func(element element.Element, property string) (string, error)

func (s *Selection) hasProperty(method propertyMethod, property, name string) (string, error) {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return "", fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	value, err := method(selectedElement, property)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve %s value for %s: %s", name, s, err)
	}
	return value, nil
}

// Attribute returns an attribute value for exactly one element.
func (s *Selection) Attribute(attribute string) (string, error) {
	return s.hasProperty(element.Element.GetAttribute, attribute, "attribute")
}

// CSS returns a CSS style property value for exactly one element.
func (s *Selection) CSS(property string) (string, error) {
	return s.hasProperty(element.Element.GetCSS, property, "CSS property")
}

type stateMethod func(element element.Element) (bool, error)

func (s *Selection) hasState(method stateMethod, name string) (bool, error) {
	elements, err := s.elements.GetAtLeastOne()
	if err != nil {
		return false, fmt.Errorf("failed to select elements from %s: %s", s, err)
	}

	for _, selectedElement := range elements {
		pass, err := method(selectedElement)
		if err != nil {
			return false, fmt.Errorf("failed to determine whether %s is %s: %s", s, name, err)
		}
		if !pass {
			return false, nil
		}
	}

	return true, nil
}

// Selected returns true if all of the elements that the selection refers to are selected.
func (s *Selection) Selected() (bool, error) {
	return s.hasState(element.Element.IsSelected, "selected")
}

// Visible returns true if all of the elements that the selection refers to are visible.
func (s *Selection) Visible() (bool, error) {
	return s.hasState(element.Element.IsDisplayed, "visible")
}

// Enabled returns true if all of the elements that the selection refers to are enabled.
func (s *Selection) Enabled() (bool, error) {
	return s.hasState(element.Element.IsEnabled, "enabled")
}
