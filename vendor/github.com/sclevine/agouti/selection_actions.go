package agouti

import (
	"fmt"
	"path/filepath"

	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/target"
)

type actionsFunc func(element.Element) error

func (s *Selection) forEachElement(actions actionsFunc) error {
	elements, err := s.elements.GetAtLeastOne()
	if err != nil {
		return fmt.Errorf("failed to select elements from %s: %s", s, err)
	}

	for _, element := range elements {
		if err := actions(element); err != nil {
			return err
		}
	}
	return nil
}

// Click clicks on all of the elements that the selection refers to.
func (s *Selection) Click() error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := selectedElement.Click(); err != nil {
			return fmt.Errorf("failed to click on %s: %s", s, err)
		}
		return nil
	})
}

// DoubleClick double-clicks on all of the elements that the selection refers to.
func (s *Selection) DoubleClick() error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := s.session.MoveTo(selectedElement.(*api.Element), nil); err != nil {
			return fmt.Errorf("failed to move mouse to %s: %s", s, err)
		}
		if err := s.session.DoubleClick(); err != nil {
			return fmt.Errorf("failed to double-click on %s: %s", s, err)
		}
		return nil
	})
}

// Clear clears all fields the selection refers to.
func (s *Selection) Clear() error {
        return s.forEachElement(func(selectedElement element.Element) error {
                if err := selectedElement.Clear(); err != nil {
                        return fmt.Errorf("failed to clear %s: %s", s, err)
                }
                return nil
        })
}

// Fill fills all of the fields the selection refers to with the provided text.
func (s *Selection) Fill(text string) error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := selectedElement.Clear(); err != nil {
			return fmt.Errorf("failed to clear %s: %s", s, err)
		}
		if err := selectedElement.Value(text); err != nil {
			return fmt.Errorf("failed to enter text into %s: %s", s, err)
		}
		return nil
	})
}

// UploadFile uploads the provided file to all selected <input type="file" />.
// The provided filename may be a relative or absolute path.
// Returns an error if elements of any other type are in the selection.
func (s *Selection) UploadFile(filename string) error {
	absFilePath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to find absolute path for filename: %s", err)
	}
	return s.forEachElement(func(selectedElement element.Element) error {
		tagName, err := selectedElement.GetName()
		if err != nil {
			return fmt.Errorf("failed to determine tag name of %s: %s", s, err)
		}
		if tagName != "input" {
			return fmt.Errorf("element for %s is not an input element", s)
		}
		inputType, err := selectedElement.GetAttribute("type")
		if err != nil {
			return fmt.Errorf("failed to determine type attribute of %s: %s", s, err)
		}
		if inputType != "file" {
			return fmt.Errorf("element for %s is not a file uploader", s)
		}
		if err := selectedElement.Value(absFilePath); err != nil {
			return fmt.Errorf("failed to enter text into %s: %s", s, err)
		}
		return nil
	})
}

// Check checks all of the unchecked checkboxes that the selection refers to.
func (s *Selection) Check() error {
	return s.setChecked(true)
}

// Uncheck unchecks all of the checked checkboxes that the selection refers to.
func (s *Selection) Uncheck() error {
	return s.setChecked(false)
}

func (s *Selection) setChecked(checked bool) error {
	return s.forEachElement(func(selectedElement element.Element) error {
		elementType, err := selectedElement.GetAttribute("type")
		if err != nil {
			return fmt.Errorf("failed to retrieve type attribute of %s: %s", s, err)
		}

		if elementType != "checkbox" {
			return fmt.Errorf("%s does not refer to a checkbox", s)
		}

		elementChecked, err := selectedElement.IsSelected()
		if err != nil {
			return fmt.Errorf("failed to retrieve state of %s: %s", s, err)
		}

		if elementChecked != checked {
			if err := selectedElement.Click(); err != nil {
				return fmt.Errorf("failed to click on %s: %s", s, err)
			}
		}
		return nil
	})
}

// Select may be called on a selection of any number of <select> elements to select
// any <option> elements under those <select> elements that match the provided text.
func (s *Selection) Select(text string) error {
	return s.forEachElement(func(selectedElement element.Element) error {
		optionXPath := fmt.Sprintf(`./option[normalize-space()="%s"]`, text)
		optionToSelect := target.Selector{Type: target.XPath, Value: optionXPath}
		options, err := selectedElement.GetElements(optionToSelect.API())
		if err != nil {
			return fmt.Errorf("failed to select specified option for %s: %s", s, err)
		}

		if len(options) == 0 {
			return fmt.Errorf(`no options with text "%s" found for %s`, text, s)
		}

		for _, option := range options {
			if err := option.Click(); err != nil {
				return fmt.Errorf(`failed to click on option with text "%s" for %s: %s`, text, s, err)
			}
		}
		return nil
	})
}

// Submit submits all selected forms. The selection may refer to a form itself
// or any input element contained within a form.
func (s *Selection) Submit() error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := selectedElement.Submit(); err != nil {
			return fmt.Errorf("failed to submit %s: %s", s, err)
		}
		return nil
	})
}

// Tap performs the provided Tap event on each element in the selection.
func (s *Selection) Tap(event Tap) error {
	var touchFunc func(*api.Element) error
	switch event {
	case SingleTap:
		touchFunc = s.session.TouchClick
	case DoubleTap:
		touchFunc = s.session.TouchDoubleClick
	case LongTap:
		touchFunc = s.session.TouchLongClick
	default:
		return fmt.Errorf("failed to %s on %s: invalid tap event", event, s)
	}

	return s.forEachElement(func(selectedElement element.Element) error {
		if err := touchFunc(selectedElement.(*api.Element)); err != nil {
			return fmt.Errorf("failed to %s on %s: %s", event, s, err)
		}
		return nil
	})
}

// Touch performs the provided Touch event at the location of each element in the
// selection.
func (s *Selection) Touch(event Touch) error {
	var touchFunc func(x, y int) error
	switch event {
	case HoldFinger:
		touchFunc = s.session.TouchDown
	case ReleaseFinger:
		touchFunc = s.session.TouchUp
	case MoveFinger:
		touchFunc = s.session.TouchMove
	default:
		return fmt.Errorf("failed to %s on %s: invalid touch event", event, s)
	}

	return s.forEachElement(func(selectedElement element.Element) error {
		x, y, err := selectedElement.GetLocation()
		if err != nil {
			return fmt.Errorf("failed to retrieve location of %s: %s", s, err)
		}
		if err := touchFunc(x, y); err != nil {
			return fmt.Errorf("failed to flick finger on %s: %s", s, err)
		}
		return nil
	})
}

// FlickFinger performs a flick touch action by the provided offset and at the
// provided speed on exactly one element.
func (s *Selection) FlickFinger(xOffset, yOffset int, speed uint) error {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	if err := s.session.TouchFlick(selectedElement.(*api.Element), api.XYOffset{X: xOffset, Y: yOffset}, api.ScalarSpeed(speed)); err != nil {
		return fmt.Errorf("failed to flick finger on %s: %s", s, err)
	}
	return nil
}

// ScrollFinger performs a scroll touch action by the provided offset on exactly
// one element.
func (s *Selection) ScrollFinger(xOffset, yOffset int) error {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	if err := s.session.TouchScroll(selectedElement.(*api.Element), api.XYOffset{X: xOffset, Y: yOffset}); err != nil {
		return fmt.Errorf("failed to scroll finger on %s: %s", s, err)
	}
	return nil
}

func (s *Selection) SendKeys(key string) error {
	return s.forEachElement(func(selectedElement element.Element) error {
		if err := selectedElement.Value(key); err != nil {
			return fmt.Errorf("failed to send key %s on %s: %s", key, s, err)
		}
		return nil
	})
}
