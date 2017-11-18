package agouti

import "github.com/sclevine/agouti/internal/target"

func NewTestSelection(session apiSession, elements elementRepository, firstSelector string) *Selection {
	selector := target.Selector{Type: target.CSS, Value: firstSelector, Single: true}
	return &Selection{selectable{session, target.Selectors{selector}}, elements}
}

func NewTestMultiSelection(session apiSession, elements elementRepository, firstSelector string) *MultiSelection {
	selector := target.Selector{Type: target.CSS, Value: firstSelector}
	selection := Selection{selectable{session, target.Selectors{selector}}, elements}
	return &MultiSelection{selection}
}

func NewTestPage(session apiSession) *Page {
	return &Page{selectable{session, nil}, nil}
}

func NewTestConfig() *config {
	return &config{}
}
