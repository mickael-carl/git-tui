package status

import (
	"fmt"
	"os"
	"time"

	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/git"
	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (s *statusPage) commit() {
	tempfile, err := os.CreateTemp("", "")
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-commit-err", fmt.Errorf("Failed to create temporary file for commit message: %v", err))
		return
	}
	defer os.Remove(tempfile.Name())

	message, err := util.ReadCommitMessage(s.state.App, tempfile.Name())
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-commit-err", fmt.Errorf("Failed to read commit message: %v", err))
		return
	}

	if err := s.state.Repository.Commit(s.state.Branch, message, ""); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-commit-err", fmt.Errorf("Failed to create commit: %v", err))
		return
	}

	s.update()
}

func (s *statusPage) amend() {
	tempfile, err := os.CreateTemp("", "")
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-amend-err", fmt.Errorf("Failed to create temporary file for commit message while amending: %v", err))
		return
	}
	defer os.Remove(tempfile.Name())

	oldMessage, err := s.state.Repository.HeadCommitMessage(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-amend-err", fmt.Errorf("Failed to get previous commit message while amending: %v", err))
		return
	}

	if _, err := tempfile.WriteString(oldMessage); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-amend-err", fmt.Errorf("Failed to write commit message to tempfile while amending: %v", err))
		return
	}

	message, err := util.ReadCommitMessage(s.state.App, tempfile.Name())
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-amend-err", fmt.Errorf("Failed to read commit message while amending: %v", err))
		return
	}

	if err := s.state.Repository.Amend(s.state.Branch, message); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-amend-err", fmt.Errorf("Failed to amend commit: %v", err))
		return
	}

	s.update()
}

func formatForFixup(commit git.Commit) string {
	return fmt.Sprintf(
		"%s\n\n%s",
		commit.Date.Format(time.RFC3339),
		commit.Message,
	)
}

func (s *statusPage) fixup() {
	commits, err := s.state.Repository.FixupList(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-fixup-err", fmt.Errorf("Failed to list commits for fixup: %v", err))
		return
	}

	if len(commits.Hashes) == 0 {
		return
	}

	list := tview.NewList()
	list.SetBorder(true)
	for _, hash := range commits.Hashes {
		list.AddItem(hash, "", 0, nil)
	}

	messageView := tview.NewTextView()
	messageView.SetBorder(true)
	// On first render, tview does not call `SetChangedFunc`, despite the docs
	// saying it does, which leaves the message blank, so we force set it to
	// the first item on the list.
	messageView.SetText(formatForFixup(commits.Commits[commits.Hashes[0]]))

	list.SetChangedFunc(func(_ int, mainText, _ string, _ rune) {
		messageView.SetText(formatForFixup(commits.Commits[mainText]))
	})
	list.SetBorderPadding(1, 1, 2, 2)

	flex := tview.NewFlex().
		AddItem(list, 0, 1, true).
		AddItem(messageView, 0, 3, false)

	flexElement := util.FlexElement{
		Primitive:    flex,
		RelativeSize: 1,
		Focus:        true,
	}

	fixupWindow := util.NewStatusFloatingWindow(s.state.Pages, "fixup", true, flexElement)

	list.SetSelectedFunc(func(_ int, mainText, _ string, _ rune) {
		if err := s.state.Repository.Fixup(s.state.Branch, mainText); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-fixup-err", fmt.Errorf("Failed to fixup commit: %v", err))
			return
		}
		fixupWindow.Hide()
	})

	fixupWindow.Show()
}

func (s *statusPage) xxx() {
	if err := s.state.Repository.AddAll(s.state.Branch); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-xxx-err", fmt.Errorf("Failed to add all files: %v", err))
		return
	}

	if err := s.state.Repository.Commit(s.state.Branch, "XXX", ""); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-xxx-err", fmt.Errorf("Failed to create XXX commit: %v", err))
	}
}
