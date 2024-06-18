package util

import (
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/git"
)

type State struct {
	App        *tview.Application
	Pages      *tview.Pages
	Repository *git.Repository
	Branch     string
}
