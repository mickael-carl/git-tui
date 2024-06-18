package repos

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/git"
	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (r *reposPage) createRepo() {
	name := r.form.GetFormItem(0).(*tview.InputField).GetText()
	url := r.form.GetFormItem(1).(*tview.InputField).GetText()
	mainBranch := r.form.GetFormItem(2).(*tview.InputField).GetText()
	pubKeyPath := r.form.GetFormItem(3).(*tview.InputField).GetText()
	privateKeyPath := r.form.GetFormItem(4).(*tview.InputField).GetText()

	r.formWindow.Hide()

	if err := git.NewRepository(name, url, mainBranch, pubKeyPath, privateKeyPath); err != nil {
		util.NewErrorWindow(r.state.Pages, "repos-create-err", fmt.Errorf("Failed to create repository %s with URL %s on main branch %s: %v", name, url, mainBranch, err))
		return
	}

	r.update()
}

func (r *reposPage) deleteRepo() {
	name, _ := r.list.GetItemText(r.list.GetCurrentItem())
	if err := git.DeleteRepository(name); err != nil {
		util.NewErrorWindow(r.state.Pages, "repos-delete-err", fmt.Errorf("Failed to delete repository %s: %v", name, err))
		return
	}

	r.update()
}
