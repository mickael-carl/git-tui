package branches

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (b *branchesPage) createBranch() {
	name := b.form.GetFormItem(0).(*tview.InputField).GetText()
	if err := b.state.Repository.CreateBranch(name); err != nil {
		util.NewErrorWindow(b.state.Pages, "branches-create-err", fmt.Errorf("Failed to create branch %s: %v", name, err))
	}

	b.formWindow.Hide()

	b.update()
}

func (b *branchesPage) deleteBranch() {
	name, _ := b.list.GetItemText(b.list.GetCurrentItem())
	if err := b.state.Repository.DeleteBranch(name); err != nil {
		util.NewErrorWindow(b.state.Pages, "branches-delete-err", fmt.Errorf("Failed to delete branch %s: %v", name, err))
		return
	}
	b.update()
}
