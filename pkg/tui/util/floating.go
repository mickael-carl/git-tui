package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FloatingWindow struct {
	name      string
	outerFlex *tview.Flex
	pages     *tview.Pages
}

type FlexElement struct {
	Primitive    tview.Primitive
	FixedSize    int
	RelativeSize int
	Focus        bool
}

func NewFloatingWindow(
	upperRow, lowerRow, leftColumn, rightColumn FlexElement,
	columnFixedSize, columnRelativeSize int,
	pages *tview.Pages,
	name string,
	focus bool,
	elements ...FlexElement,
) *FloatingWindow {
	innerFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, upperRow.FixedSize, upperRow.RelativeSize, false)
	for _, e := range elements {
		innerFlex.AddItem(e.Primitive, e.FixedSize, e.RelativeSize, e.Focus)
	}
	innerFlex.AddItem(nil, lowerRow.FixedSize, lowerRow.RelativeSize, false)

	outerFlex := tview.NewFlex().
		AddItem(nil, leftColumn.FixedSize, leftColumn.RelativeSize, false).
		AddItem(innerFlex, columnFixedSize, columnRelativeSize, focus).
		AddItem(nil, rightColumn.FixedSize, rightColumn.RelativeSize, false)

	window := &FloatingWindow{
		name:      name,
		outerFlex: outerFlex,
		pages:     pages,
	}

	outerFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Name() == "Esc" {
			window.Hide()
		}
		return event
	})
	return window
}

func NewTextNotificationWindow(pages *tview.Pages, name, text string) *FloatingWindow {
	textView := tview.NewTextView().SetText(text).SetTextAlign(tview.AlignCenter)
	textView.SetBorder(true)

	textElement := FlexElement{
		Primitive: textView,
		// The text itself takes one line, and the borders 2.
		FixedSize: 3,
	}
	upperRow := FlexElement{
		RelativeSize: 1,
	}
	lowerRow := FlexElement{
		RelativeSize: 1,
	}
	leftColumn := FlexElement{
		RelativeSize: 1,
	}
	rightColumn := FlexElement{
		RelativeSize: 1,
	}
	return NewFloatingWindow(
		upperRow, lowerRow, leftColumn, rightColumn,
		0, 1,
		pages, name, true, textElement,
	)
}

func NewMiddleFloatingWindow(pages *tview.Pages, name string, focus bool, elements ...FlexElement) *FloatingWindow {
	upperRow := FlexElement{
		FixedSize: 2,
	}
	lowerRow := FlexElement{
		FixedSize: 2,
	}
	leftColumn := FlexElement{
		RelativeSize: 1,
	}
	rightColumn := FlexElement{
		RelativeSize: 1,
	}

	return NewFloatingWindow(
		upperRow, lowerRow, leftColumn, rightColumn,
		0, 3,
		pages,
		name,
		focus,
		elements...,
	)
}

// The status page has a header for each column, and so the upper padding needs
// to be a bit larger otherwise the window will render a bit oddly.
func NewStatusFloatingWindow(pages *tview.Pages, name string, focus bool, elements ...FlexElement) *FloatingWindow {
	upperRow := FlexElement{
		FixedSize: 4,
	}
	lowerRow := FlexElement{
		FixedSize: 2,
	}
	leftColumn := FlexElement{
		RelativeSize: 1,
	}
	rightColumn := FlexElement{
		RelativeSize: 1,
	}

	return NewFloatingWindow(
		upperRow, lowerRow, leftColumn, rightColumn,
		0, 2,
		pages,
		name,
		focus,
		elements...,
	)
}

func NewErrorWindow(pages *tview.Pages, name string, err error) {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage(name)
		}).
		SetBackgroundColor(tcell.ColorDarkRed)
	modal.SetTitle("Error")
	pages.AddPage(name, modal, true, true)
}

func NewWarningWindow(pages *tview.Pages, name string, err error) {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage(name)
		}).
		SetBackgroundColor(tcell.ColorDarkOrange)
	modal.SetTitle("Warning")
	pages.AddPage(name, modal, true, true)
}

func (f *FloatingWindow) Show() {
	f.pages.AddPage(
		f.name,
		f.outerFlex,
		true,
		true,
	)
}

func (f *FloatingWindow) Hide() {
	f.pages.RemovePage(f.name)
}
