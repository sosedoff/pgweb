package target

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

type Type string

const (
	CSS        Type = "CSS: %s"
	XPath      Type = "XPath: %s"
	Link       Type = `Link: "%s"`
	Label      Type = `Label: "%s"`
	Button     Type = `Button: "%s"`
	Name       Type = `Name: "%s"`
	A11yID     Type = "Accessibility ID: %s"
	AndroidAut Type = "Android UIAut.: %s"
	IOSAut     Type = "iOS UIAut.: %s"
	Class      Type = "Class: %s"
	ID         Type = "ID: %s"

	labelXPath  = `//input[@id=(//label[normalize-space()="%s"]/@for)] | //label[normalize-space()="%[1]s"]/input`
	buttonXPath = `//input[@type="submit" or @type="button"][normalize-space(@value)="%s"] | //button[normalize-space()="%[1]s"]`
)

func (t Type) format(value string) string {
	return fmt.Sprintf(string(t), value)
}

type Selector struct {
	Type    Type
	Value   string
	Index   int
	Indexed bool
	Single  bool
}

func (s Selector) String() string {
	var suffix string

	if s.Single {
		suffix = " [single]"
	} else if s.Indexed {
		suffix = fmt.Sprintf(" [%d]", s.Index)
	}

	return s.Type.format(s.Value) + suffix
}

func (s Selector) API() api.Selector {
	return api.Selector{Using: s.apiType(), Value: s.value()}
}

func (s Selector) apiType() string {
	switch s.Type {
	case CSS:
		return "css selector"
	case Class:
		return "class name"
	case ID:
		return "id"
	case Link:
		return "link text"
	case Name:
		return "name"
	case A11yID:
		return "accessibility id"
	case AndroidAut:
		return "-android uiautomator"
	case IOSAut:
		return "-ios uiautomation"
	}
	return "xpath"
}

func (s Selector) value() string {
	switch s.Type {
	case Label:
		return fmt.Sprintf(labelXPath, s.Value)
	case Button:
		return fmt.Sprintf(buttonXPath, s.Value)
	}
	return s.Value
}
