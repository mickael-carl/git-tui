package util

import (
	"github.com/rivo/tview"
)

const (
	StatusPageName   = "status"
	BranchesPageName = "branches"
	ReposPageName    = "repos"
)

type Page struct {
	MainFlex *tview.Flex
}

func NewPage(primitive tview.Primitive, keybindings []Keybinding) Page {
	footer := newHelp(keybindings)
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(primitive, 0, 1, true).
		AddItem(footer, 4, 0, false)
	page := Page{
		MainFlex: mainFlex,
	}
	return page
}

func newHelp(keybindings []Keybinding) *tview.Table {
	helpTable := tview.NewTable()
	helpTable.SetBorderPadding(0, 0, 2, 2)
	row := 0
	column := 0
	for _, kb := range keybindings {
		if row > 3 {
			row = 0
			column += 4
		}
		helpTable.SetCell(
			row, column,
			tview.NewTableCell(kb.Key),
		)
		helpTable.SetCell(
			row, column+1,
			tview.NewTableCell(kb.Description).SetAlign(tview.AlignRight),
		)
		row++
	}
	return helpTable
}
