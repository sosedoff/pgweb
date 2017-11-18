package api

import (
	"errors"
	"path"
	"strings"
)

type Element struct {
	ID      string
	Session *Session
}

func (e *Element) Send(method, endpoint string, body, result interface{}) error {
	return e.Session.Send(method, path.Join("element", e.ID, endpoint), body, result)
}

func (e *Element) GetID() string {
	return e.ID
}

func (e *Element) GetElement(selector Selector) (*Element, error) {
	var result struct{ Element string }

	if err := e.Send("POST", "element", selector, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, e.Session}, nil
}

func (e *Element) GetElements(selector Selector) ([]*Element, error) {
	var results []struct{ Element string }

	if err := e.Send("POST", "elements", selector, &results); err != nil {
		return nil, err
	}

	elements := []*Element{}
	for _, result := range results {
		elements = append(elements, &Element{result.Element, e.Session})
	}

	return elements, nil
}

func (e *Element) GetText() (string, error) {
	var text string
	if err := e.Send("GET", "text", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (e *Element) GetName() (string, error) {
	var name string
	if err := e.Send("GET", "name", nil, &name); err != nil {
		return "", err
	}
	return name, nil
}

func (e *Element) GetAttribute(attribute string) (string, error) {
	var value string
	if err := e.Send("GET", path.Join("attribute", attribute), nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (e *Element) GetCSS(property string) (string, error) {
	var value string
	if err := e.Send("GET", path.Join("css", property), nil, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (e *Element) Click() error {
	return e.Send("POST", "click", nil, nil)
}

func (e *Element) Clear() error {
	return e.Send("POST", "clear", nil, nil)
}

func (e *Element) Value(text string) error {
	splitText := strings.Split(text, "")
	request := struct {
		Value []string `json:"value"`
	}{splitText}
	return e.Send("POST", "value", request, nil)
}

func (e *Element) IsSelected() (bool, error) {
	var selected bool
	if err := e.Send("GET", "selected", nil, &selected); err != nil {
		return false, err
	}
	return selected, nil
}

func (e *Element) IsDisplayed() (bool, error) {
	var displayed bool
	if err := e.Send("GET", "displayed", nil, &displayed); err != nil {
		return false, err
	}
	return displayed, nil
}

func (e *Element) IsEnabled() (bool, error) {
	var enabled bool
	if err := e.Send("GET", "enabled", nil, &enabled); err != nil {
		return false, err
	}
	return enabled, nil
}

func (e *Element) Submit() error {
	return e.Send("POST", "submit", nil, nil)
}

func (e *Element) IsEqualTo(other *Element) (bool, error) {
	if other == nil {
		return false, errors.New("nil element is invalid")
	}
	var equal bool
	if err := e.Send("GET", path.Join("equals", other.ID), nil, &equal); err != nil {
		return false, err
	}
	return equal, nil
}

func (e *Element) GetLocation() (x, y int, err error) {
	var location struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	if err := e.Send("GET", "location", nil, &location); err != nil {
		return 0, 0, err
	}
	return round(location.X), round(location.Y), nil
}

func (e *Element) GetSize() (width, height int, err error) {
	var size struct {
		Width float64 `json:"width"`
		Height float64 `json:"height"`
	}
	if err := e.Send("GET", "size", nil, &size); err != nil {
		return 0, 0, err
	}
	return round(size.Width), round(size.Height), nil
}

func round(number float64) int {
	return int(number + 0.5)
}
