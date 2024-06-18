package repos

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/git"
	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

type reposPage struct {
	util.Page
	state       *util.State
	list        *tview.List
	flex        *tview.Flex
	form        *tview.Form
	formWindow  *util.FloatingWindow
	keybindings []util.Keybinding
}

func NewReposPage(state *util.State) *reposPage {
	rp := &reposPage{
		state: state,
	}

	list := tview.NewList().SetSelectedFunc(rp.setRepo)
	list.SetFocusFunc(rp.update)
	list.SetInputCapture(rp.handleInput)
	rp.list = list

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true)
	rp.flex = flex

	form := tview.NewForm().
		AddInputField("Name", "", 0, nil, nil).
		AddInputField("URL", "", 0, nil, nil).
		AddInputField("Main Branch Name", "", 0, nil, nil).
		AddInputField("Public Key Path", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub"), 0, nil, nil).
		AddInputField("Private Key Path", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"), 0, nil, nil).
		AddButton("Create", rp.createRepo)
	form.SetBorder(true)
	rp.form = form

	formElement := util.FlexElement{
		Primitive:    form,
		RelativeSize: 1,
		Focus:        true,
	}
	rp.formWindow = util.NewMiddleFloatingWindow(state.Pages, "newRepo", true, formElement)

	kbs := []util.Keybinding{
		{
			Key:         "r",
			Description: "refresh",
			Action:      rp.update,
		},
		{
			Key:         "n",
			Description: "new",
			Action:      rp.formWindow.Show,
		},
		{
			Key:         "d",
			Description: "delete",
			Action:      rp.deleteRepo,
		},
	}
	rp.keybindings = kbs

	page := util.NewPage(flex, kbs)
	rp.Page = page

	state.Pages.AddPage(util.ReposPageName, page.MainFlex, true, true)

	return rp
}

func (r *reposPage) update() {
	r.list.Clear()

	repos, err := git.GetRepositories()
	if err != nil {
		util.NewErrorWindow(r.state.Pages, "repos-update-err", err)
		return
	}

	i := '1'
	for _, repo := range repos {
		r.list.AddItem(repo.Name, repo.Config.Remote.URL, i, nil)
		i++
	}
}

func (r *reposPage) setRepo(index int, mainText, secondaryText string, shortcut rune) {
	repos, err := git.GetRepositories()
	if err != nil {
		util.NewErrorWindow(r.state.Pages, "repos-set-repo-err", err)
		return
	}

	for _, repo := range repos {
		if repo.Name == mainText {
			r.state.Repository = &repo
			break
		}
	}

	r.state.Pages.SwitchToPage(util.BranchesPageName)
}

func (r *reposPage) handleInput(event *tcell.EventKey) *tcell.EventKey {
	for _, kb := range r.keybindings {
		if kb.Key == event.Name() || kb.Key == string(event.Rune()) {
			kb.Action()
		}
	}
	return event
}
