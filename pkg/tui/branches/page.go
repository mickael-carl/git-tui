package branches

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

type branchesPage struct {
	util.Page
	state       *util.State
	list        *tview.List
	form        *tview.Form
	formWindow  *util.FloatingWindow
	keybindings []util.Keybinding
}

func NewBranchesPage(state *util.State) *branchesPage {
	bp := &branchesPage{
		state: state,
	}

	list := tview.NewList().SetSelectedFunc(bp.setBranch)
	list.SetFocusFunc(bp.update)
	list.SetInputCapture(bp.handleBranchInput)
	bp.list = list

	form := tview.NewForm().
		AddInputField("Name", "", 0, nil, nil).
		AddButton("Create", bp.createBranch)
	form.SetBorder(true)
	bp.form = form

	formElement := util.FlexElement{
		Primitive:    form,
		RelativeSize: 1,
		Focus:        true,
	}
	bp.formWindow = util.NewMiddleFloatingWindow(state.Pages, "newBranch", true, formElement)

	kbs := []util.Keybinding{
		{
			Key:         "f",
			Description: "fetch",
			Action:      bp.fetch,
		},
		{
			Key:         "r",
			Description: "refresh",
			Action:      bp.update,
		},
		{
			Key:         "d",
			Description: "delete",
			Action:      bp.deleteBranch,
		},
		{
			Key:         "Esc",
			Description: "back",
			Action:      bp.back,
		},
		{
			Key:         "n",
			Description: "new branch",
			Action:      bp.formWindow.Show,
		},
	}
	bp.keybindings = kbs

	page := util.NewPage(list, kbs)
	bp.Page = page

	state.Pages.AddPage(util.BranchesPageName, page.MainFlex, true, false)

	return bp
}

func (b *branchesPage) update() {
	b.list.Clear()

	i := '1'
	for _, branch := range b.state.Repository.Branches {
		b.list.AddItem(
			branch.Name,
			fmt.Sprintf(
				"↑ %d, ↓ %d, ⌚️ %s",
				branch.CommitsAhead,
				branch.CommitsBehind,
				branch.LastCommit.Format(time.DateTime),
			),
			i,
			nil,
		)
		i++
	}
}

func (b *branchesPage) back() {
	b.state.Pages.SwitchToPage(util.ReposPageName)
}

func (b *branchesPage) handleBranchInput(event *tcell.EventKey) *tcell.EventKey {
	for _, kb := range b.keybindings {
		if kb.Key == event.Name() || kb.Key == string(event.Rune()) {
			kb.Action()
		}
	}
	return event
}

func (b *branchesPage) setBranch(index int, mainText, secondaryText string, shortcut rune) {
	b.state.Branch = mainText
	b.state.Pages.SwitchToPage(util.StatusPageName)
}
