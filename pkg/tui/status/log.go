package status

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/muesli/gamut"
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (s *statusPage) log() {
	url, err := s.state.Repository.RemoteHTTPURL()
	if err != nil {
		util.NewWarningWindow(s.state.Pages, "status-log-err", fmt.Errorf("Failed to get repository remote HTTP URL: %v", err))
	}

	textView := tview.NewTextView().SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetBorderPadding(1, 1, 2, 2)

	textViewElement := util.FlexElement{
		Primitive:    textView,
		RelativeSize: 1,
		Focus:        true,
	}

	logWindow := util.NewStatusFloatingWindow(s.state.Pages, "log", true, textViewElement)
	logWindow.Show()

	commits, err := s.state.Repository.Log(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-log-err", fmt.Errorf("Failed to get Git log: %v", err))
		return
	}

	rgbaColors, err := gamut.Generate(4, gamut.PastelGenerator{})
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-log-err", fmt.Errorf("Failed to generate color palette: %v", err))
		return
	}
	colors := []string{}
	for _, color := range rgbaColors {
		colors = append(colors, gamut.ToHex(color))
	}

	builder := strings.Builder{}
	regex := regexp.MustCompile(`#(\d+)`)
	for _, commit := range commits {
		if _, err := builder.WriteString(fmt.Sprintf(
			"[%s::bu]Commit\t%s[-::-]\n",
			colors[0],
			commit.Hash,
		)); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-log-err", fmt.Errorf("Failed to write to string while building Git log: %v", err))
			return
		}

		if _, err := builder.WriteString(fmt.Sprintf(
			"[::u]Author:[::-]\t%s\n",
			commit.Author,
		)); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-log-err", fmt.Errorf("Failed to write to string while building Git log: %v", err))
			return
		}

		if _, err := builder.WriteString(fmt.Sprintf(
			"[::u]Date:[::-]\t%s\n\n",
			commit.Date.Format(time.RFC3339),
		)); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-log-err", fmt.Errorf("Failed to write to string while building Git log: %v", err))
			return
		}

		if _, err := builder.WriteString(fmt.Sprintf(
			"%s\n\n",
			// As far as Github is concerned, `issues/X` redirects to `pull/X`
			// if X is a PR and not an actual issue.
			regex.ReplaceAllString(commit.Message, fmt.Sprintf("[:::%s/issues/$1]$0[:::-]", url)),
		)); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-log-err", fmt.Errorf("Failed to write to string while building Git log: %v", err))
			return
		}
	}

	textView.SetText(builder.String())
}
