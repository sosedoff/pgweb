package agouti

import (
	"fmt"

	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/target"
)

// Selection instances refer to a selection of elements.
// All Selection methods are also MultiSelection methods.
//
// Methods that take selectors apply their selectors to each element in the
// selection they are called on. If the selection they are called on refers to multiple
// elements, the resulting selection will refer to at least that many elements.
//
// Examples:
//
//    selection.Find("table").All("tr").At(2).First("td input[type=checkbox]").Check()
// Checks the first checkbox in the third row of the only table.
//    selection.Find("table").All("tr").Find("td").All("input[type=checkbox]").Check()
// Checks all checkboxes in the first-and-only cell of each row in the only table.
type Selection struct {
	selectable
	elements elementRepository
}

type elementRepository interface {
	Get() ([]element.Element, error)
	GetAtLeastOne() ([]element.Element, error)
	GetExactlyOne() (element.Element, error)
}

func newSelection(session apiSession, selectors target.Selectors) *Selection {
	return &Selection{
		selectable{session, selectors},
		&element.Repository{
			Client:    session,
			Selectors: selectors,
		},
	}
}

// String returns a string representation of the selection, ex.
//    selection 'CSS: .some-class | XPath: //table [3] | Link "click me" [single]'
func (s *Selection) String() string {
	return fmt.Sprintf("selection '%s'", s.selectors)
}

// Elements returns a []*api.Element that can be used to send direct commands
// to WebDriver elements. See: https://code.google.com/p/selenium/wiki/JsonWireProtocol
func (s *Selection) Elements() ([]*api.Element, error) {
	elements, err := s.elements.Get()
	if err != nil {
		return nil, err
	}
	apiElements := []*api.Element{}
	for _, selectedElement := range elements {
		apiElements = append(apiElements, selectedElement.(*api.Element))
	}
	return apiElements, nil
}

// Count returns the number of elements that the selection refers to.
func (s *Selection) Count() (int, error) {
	elements, err := s.elements.Get()
	if err != nil {
		return 0, fmt.Errorf("failed to select elements from %s: %s", s, err)
	}

	return len(elements), nil
}

// EqualsElement returns whether or not two selections of exactly
// one element refer to the same element.
func (s *Selection) EqualsElement(other interface{}) (bool, error) {
	otherSelection, ok := other.(*Selection)
	if !ok {
		multiSelection, ok := other.(*MultiSelection)
		if !ok {
			return false, fmt.Errorf("must be *Selection or *MultiSelection")
		}
		otherSelection = &multiSelection.Selection
	}

	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	otherElement, err := otherSelection.elements.GetExactlyOne()
	if err != nil {
		return false, fmt.Errorf("failed to select element from %s: %s", other, err)
	}

	equal, err := selectedElement.IsEqualTo(otherElement.(*api.Element))
	if err != nil {
		return false, fmt.Errorf("failed to compare %s to %s: %s", s, other, err)
	}

	return equal, nil
}

// MouseToElement moves the mouse over exactly one element in the selection.
func (s *Selection) MouseToElement() error {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	if err := s.session.MoveTo(selectedElement.(*api.Element), nil); err != nil {
		return fmt.Errorf("failed to move mouse to element for %s: %s", s, err)
	}

	return nil
}
