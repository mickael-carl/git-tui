package status

import (
	"fmt"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (s *statusPage) reset() {
	if err := s.state.Repository.ResetPreviousCommit(s.state.Branch); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-reset-err", fmt.Errorf("Failed to reset to previous commit: %v", err))
		return
	}

	s.update()
}
