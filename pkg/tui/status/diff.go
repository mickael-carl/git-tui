package status

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (s *statusPage) diff() {
	diffs, err := s.state.Repository.Diff(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-diff-err", fmt.Errorf("Failed to get diff: %v", err))
		return
	}

	if len(diffs) == 0 {
		return
	}

	textView := tview.NewTextView().SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetBorderPadding(1, 1, 2, 2)

	textViewElement := util.FlexElement{
		Primitive:    textView,
		RelativeSize: 1,
		Focus:        true,
	}

	diffWindow := util.NewStatusFloatingWindow(s.state.Pages, "diff", true, textViewElement)
	diffWindow.Show()

	builder := strings.Builder{}
	for _, diff := range diffs {
		if _, err := builder.WriteString(
			diff.Patch,
		); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-diff-err", fmt.Errorf("Failed to write to string while building diff: %v", err))
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
		util.NewErrorWindow(s.state.Pages, "status-diff-err", fmt.Errorf("Failed to highlight diff: %v", err))
		return
	}

	// TODO: add an input field for filepath. On change, update the diff to
	// only include those files.
	textView.SetText(output.String())
}

// TODO: this is almost identical to `diff` above.
func (s *statusPage) diffStaged() {
	diffs, err := s.state.Repository.DiffStaged(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-diff-staged-err", fmt.Errorf("Failed to get diff: %v", err))
		return
	}

	if len(diffs) == 0 {
		return
	}

	textView := tview.NewTextView().SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetBorderPadding(1, 1, 2, 2)

	textViewElement := util.FlexElement{
		Primitive:    textView,
		RelativeSize: 1,
		Focus:        true,
	}

	diffWindow := util.NewStatusFloatingWindow(s.state.Pages, "diff", true, textViewElement)
	diffWindow.Show()

	builder := strings.Builder{}
	for _, diff := range diffs {
		if _, err := builder.WriteString(
			diff.Patch,
		); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-diff-staged-err", fmt.Errorf("Failed to write to string while building diff: %v", err))
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
		util.NewErrorWindow(s.state.Pages, "status-diff-staged-err", fmt.Errorf("Failed to highlight diff: %v", err))
		return
	}

	// TODO: add an input field for filepath. On change, update the diff to
	// only include those files.
	textView.SetText(output.String())
}
