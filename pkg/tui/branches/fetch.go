package branches

import (
	"fmt"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (b *branchesPage) fetch() {
	done := util.NewProgressWindow(b.state.Pages, "branches-fetch-progress", "Fetching...")

	go b.state.Repository.Fetch(func(err error) {
		done()
		b.update()

		// TODO: this seems necessary, otherwise the progress window is never
		// removed.
		b.state.App.Draw()

		if err != nil {
			util.NewErrorWindow(b.state.Pages, "branches-fetch-err", fmt.Errorf("Failed to fetch remotes: %v", err))
		}
	})
}
