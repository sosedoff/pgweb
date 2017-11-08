package agouti

import "github.com/sclevine/agouti/internal/target"

// A MultiSelection is a Selection that may be indexed using the At() method.
// All Selection methods are available on a MultiSelection.
//
// A Selection returned by At() may still refer to multiple elements if any
// parent of the MultiSelection refers to multiple elements.
//
// Examples:
//    selection.All("section").All("form").At(1).Submit()
// Submits the second form in each section.
//    selection.All("div").Find("h1").Click()
// Clicks one h1 in each div, failing if any div does not contain exactly one h1.
type MultiSelection struct {
	Selection
}

func newMultiSelection(session apiSession, selectors target.Selectors) *MultiSelection {
	return &MultiSelection{*newSelection(session, selectors)}
}

// At finds an element at the provided index. It only applies to the immediate selection,
// meaning that the returned selection may still refer to multiple elements if any parent
// of the immediate selection is also a *MultiSelection.
func (s *MultiSelection) At(index int) *Selection {
	return newSelection(s.session, s.selectors.At(index))
}
