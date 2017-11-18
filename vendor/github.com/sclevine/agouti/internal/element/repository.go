package element

import (
	"errors"
	"fmt"

	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/target"
)

type Repository struct {
	Client    Client
	Selectors target.Selectors
}

type Client interface {
	GetElement(selector api.Selector) (*api.Element, error)
	GetElements(selector api.Selector) ([]*api.Element, error)
}

type Element interface {
	Client
	GetID() string
	GetText() (string, error)
	GetName() (string, error)
	GetAttribute(attribute string) (string, error)
	GetCSS(property string) (string, error)
	IsSelected() (bool, error)
	IsDisplayed() (bool, error)
	IsEnabled() (bool, error)
	IsEqualTo(other *api.Element) (bool, error)
	Click() error
	Clear() error
	Value(text string) error
	Submit() error
	GetLocation() (x, y int, err error)
}

func (e *Repository) GetAtLeastOne() ([]Element, error) {
	elements, err := e.Get()
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, errors.New("no elements found")
	}

	return elements, nil
}

func (e *Repository) GetExactlyOne() (Element, error) {
	elements, err := e.GetAtLeastOne()
	if err != nil {
		return nil, err
	}

	if len(elements) > 1 {
		return nil, fmt.Errorf("method does not support multiple elements (%d)", len(elements))
	}

	return elements[0], nil
}

func (e *Repository) Get() ([]Element, error) {
	if len(e.Selectors) == 0 {
		return nil, errors.New("empty selection")
	}

	lastElements, err := retrieveElements(e.Client, e.Selectors[0])
	if err != nil {
		return nil, err
	}

	for _, selector := range e.Selectors[1:] {
		elements := []Element{}
		for _, element := range lastElements {
			subElements, err := retrieveElements(element, selector)
			if err != nil {
				return nil, err
			}

			elements = append(elements, subElements...)
		}
		lastElements = elements
	}
	return lastElements, nil
}

func retrieveElements(client Client, selector target.Selector) ([]Element, error) {
	if selector.Single {
		elements, err := client.GetElements(selector.API())
		if err != nil {
			return nil, err
		}

		if len(elements) == 0 {
			return nil, errors.New("element not found")
		} else if len(elements) > 1 {
			return nil, errors.New("ambiguous find")
		}

		return []Element{Element(elements[0])}, nil
	}

	if selector.Indexed && selector.Index > 0 {
		elements, err := client.GetElements(selector.API())
		if err != nil {
			return nil, err
		}

		if selector.Index >= len(elements) {
			return nil, errors.New("element index out of range")
		}

		return []Element{Element(elements[selector.Index])}, nil
	}

	if selector.Indexed && selector.Index == 0 {
		element, err := client.GetElement(selector.API())
		if err != nil {
			return nil, err
		}
		return []Element{Element(element)}, nil
	}

	elements, err := client.GetElements(selector.API())
	if err != nil {
		return nil, err
	}

	newElements := []Element{}
	for _, element := range elements {
		newElements = append(newElements, element)
	}

	return newElements, nil
}
