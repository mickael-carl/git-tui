package tui

import (
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/tui/branches"
	"github.com/mickael-carl/git-tui/pkg/tui/repos"
	"github.com/mickael-carl/git-tui/pkg/tui/status"
	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func Run() error {
	// The default tview behaviour is to put double lines around focused
	// primitives when borders are on. I don't particularly like that
	// behaviour, so this overrides the default border to all be identical when
	// focused.
	tview.Borders.HorizontalFocus = tview.BoxDrawingsLightHorizontal
	tview.Borders.VerticalFocus = tview.BoxDrawingsLightVertical
	tview.Borders.TopLeftFocus = tview.BoxDrawingsLightDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsLightDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsLightUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsLightUpAndLeft

	app := tview.NewApplication()

	// TODO: if run from a worktree, open directly to the right repo and
	// branch.
	pages := tview.NewPages().SwitchToPage(util.ReposPageName)

	app.SetRoot(pages, true)
	state := &util.State{
		App:   app,
		Pages: pages,
	}

	status.NewStatusPage(state)
	branches.NewBranchesPage(state)
	repos.NewReposPage(state)

	return app.Run()
}
