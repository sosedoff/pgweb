package agouti

import (
	"fmt"

	"github.com/sclevine/agouti/api"
)

// SwitchToFrame focuses on the frame specified by the selection. All new and
// existing selections will refer to the new frame. All further Page methods
// will apply to this frame as well.
func (s *Selection) SwitchToFrame() error {
	selectedElement, err := s.elements.GetExactlyOne()
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %s", s, err)
	}

	if err := s.session.Frame(selectedElement.(*api.Element)); err != nil {
		return fmt.Errorf("failed to switch to frame referred to by %s: %s", s, err)
	}
	return nil
}
