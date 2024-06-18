package status

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

// TODO: a lot of this is almost identical to the `diff` function.
func (s *statusPage) show() {
	textView := tview.NewTextView().SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetBorderPadding(1, 1, 2, 2)

	textViewElement := util.FlexElement{
		Primitive:    textView,
		RelativeSize: 1,
		Focus:        true,
	}

	showWindow := util.NewStatusFloatingWindow(s.state.Pages, "show", true, textViewElement)
	showWindow.Show()

	diffs, err := s.state.Repository.Show(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-show-err", fmt.Errorf("Failed to get diff for last commit: %v", err))
		return
	}

	builder := strings.Builder{}
	for _, diff := range diffs {
		if _, err := builder.WriteString(
			diff.Patch,
		); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-show-err", fmt.Errorf("Failed to write to string while building diff for last commit: %v", err))
			return
		}
	}

	output := strings.Builder{}
	if err := quick.Highlight(
		tview.ANSIWriter(&output),
		builder.String(),
		"diff",
		"terminal16m",
		"monokai",
	); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-show-err", fmt.Errorf("Failed to highlight diff for last commit: %v", err))
		return
	}

	// TODO: add an input field for filepath. On change, update the diff to
	// only include those files.
	textView.SetText(output.String())
}
