package branches

import (
	"fmt"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (b *branchesPage) fetch() {
	fetchWindow := util.NewTextNotificationWindow(b.state.Pages, "fetch", "Fetching...")
	fetchWindow.Show()
	defer fetchWindow.Hide()

	// TODO: this seems necessary, otherwise the page is just not displayed.
	b.state.App.ForceDraw()

	if err := b.state.Repository.Fetch(); err != nil {
		util.NewErrorWindow(b.state.Pages, "branches-fetch-err", fmt.Errorf("Failed to fetch remotes: %v", err))
		return
	}

	b.update()
}
