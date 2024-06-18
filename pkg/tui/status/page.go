package status

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

type statusPage struct {
	util.Page
	state             *util.State
	grid              *tview.Grid
	stagedTextView    *tview.TextView
	unstagedTextView  *tview.TextView
	untrackedTextView *tview.TextView
	keybindings       []util.Keybinding
}

func NewStatusPage(state *util.State) *statusPage {
	stagedTextView := tview.NewTextView()
	stagedTextView.SetChangedFunc(func() { state.App.Draw() })
	unstagedTextView := tview.NewTextView()
	unstagedTextView.SetChangedFunc(func() { state.App.Draw() })
	untrackedTextView := tview.NewTextView()
	untrackedTextView.SetChangedFunc(func() { state.App.Draw() })

	sp := &statusPage{
		state:             state,
		stagedTextView:    stagedTextView,
		unstagedTextView:  unstagedTextView,
		untrackedTextView: untrackedTextView,
	}

	grid := tview.NewGrid().
		// Set 3 rows (header, content and action) and 3 columns (staged,
		// unstaged and untracked). The 0 value allow for filling up
		// proportionally leftover space from fixed-size rows/columns.
		SetRows(1, 0).
		SetColumns(0, 0, 0).
		SetBorders(true).
		AddItem(tview.NewTextView().SetText("Staged").SetTextAlign(tview.AlignCenter), 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetText("Unstaged").SetTextAlign(tview.AlignCenter), 0, 1, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetText("Untracked").SetTextAlign(tview.AlignCenter), 0, 2, 1, 1, 0, 0, false).
		AddItem(stagedTextView, 1, 0, 1, 1, 0, 0, false).
		AddItem(unstagedTextView, 1, 1, 1, 1, 0, 0, false).
		AddItem(untrackedTextView, 1, 2, 1, 1, 0, 0, false)

	grid.SetFocusFunc(sp.watchAndUpdate)
	grid.SetInputCapture(sp.handleInput)
	sp.grid = grid

	kbs := []util.Keybinding{
		{
			Key:         "Esc",
			Description: "Back",
			Action:      sp.back,
		},
		{
			Key:         "a",
			Description: "Add",
			Action:      sp.add,
		},
		{
			Key:         "u",
			Description: "Unstage",
			Action:      sp.unstage,
		},
		{
			Key:         "c",
			Description: "Commit",
			Action:      sp.commit,
		},
		{
			Key:         "l",
			Description: "Log",
			Action:      sp.log,
		},
		{
			Key:         "d",
			Description: "Diff",
			Action:      sp.diff,
		},
		{
			Key:         "D",
			Description: "Diff Staged Changes",
			Action:      sp.diffStaged,
		},
		{
			Key:         "x",
			Description: "XXX Commit",
			Action:      sp.xxx,
		},
		{
			Key:         "s",
			Description: "Show Last Commit",
			Action:      sp.show,
		},
		{
			Key:         "C",
			Description: "Amend Last Commit",
			Action:      sp.amend,
		},
		{
			Key:         "R",
			Description: "Reset to Previous Commit",
			Action:      sp.reset,
		},
		{
			Key:         "f",
			Description: "Fixup",
			Action:      sp.fixup,
		},
		{
			Key:         "r",
			Description: "Restore",
			Action:      sp.restore,
		},
		{
			Key:         "p",
			Description: "Patch",
			Action:      sp.patch,
		},
	}

	sp.keybindings = kbs
	sp.Page = util.NewPage(grid, kbs)
	state.Pages.AddPage(util.StatusPageName, sp.Page.MainFlex, true, false)

	return sp
}

func (s *statusPage) update() {
	status, err := s.state.Repository.Status(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-update-err", fmt.Errorf("Failed to get status: %v", err))
		return
	}
	s.stagedTextView.SetText(status.FormatStaged())
	s.unstagedTextView.SetText(status.FormatUnstaged())
	s.untrackedTextView.SetText(status.FormatUntracked())
}

func (s *statusPage) back() {
	s.state.Pages.SwitchToPage(util.BranchesPageName)
}

func (s *statusPage) handleInput(event *tcell.EventKey) *tcell.EventKey {
	for _, kb := range s.keybindings {
		if kb.Key == event.Name() || kb.Key == string(event.Rune()) {
			kb.Action()
		}
	}
	return event
}
