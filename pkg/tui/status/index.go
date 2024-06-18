package status

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"unicode"

	"github.com/creack/pty/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"

	"github.com/mickael-carl/git-tui/pkg/git"
	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (s *statusPage) add() {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textViewElement := util.FlexElement{
		Primitive:    textView,
		RelativeSize: 1,
	}

	// TODO: create a form, with checkboxes for every file?
	inputField := tview.NewInputField().SetPlaceholder("files:")
	inputFieldElement := util.FlexElement{
		Primitive: inputField,
		FixedSize: 1,
		Focus:     true,
	}

	addWindow := util.NewStatusFloatingWindow(s.state.Pages, "add", true, textViewElement, inputFieldElement)
	addWindow.Show()

	status, err := s.state.Repository.Status(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-add-err", fmt.Errorf("Failed to get status: %v", err))
		return
	}

	textView.SetText(strings.Join(status.ChangedFilesPaths(), "\n"))

	inputField.SetChangedFunc(func(text string) {
		changedFilesPaths := status.ChangedFilesPaths()
		if slices.Contains(changedFilesPaths, text) {
			textView.SetText(text)
			return
		}
		matched := fuzzy.Find(text, changedFilesPaths)
		textView.SetText(strings.Join(matched, "\n"))
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if err := s.state.Repository.Add(
				s.state.Branch,
				strings.Split(textView.GetText(true), "\n"),
			); err != nil {
				util.NewErrorWindow(s.state.Pages, "status-add-err", fmt.Errorf("Failed to add files to index: %v", err))
				return
			}
			addWindow.Hide()
		}
	})
}

func (s *statusPage) patch() {
	status, err := s.state.Repository.Status(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Failed to get repository status when adding patches: %v", err))
		return
	}

	if len(status.Unstaged) == 0 {
		util.NewErrorWindow(s.state.Pages, "status-patch-err", errors.New("Can't add patches: no changes left unstaged"))
		return
	}

	out := tview.NewTextView()
	out.SetBorder(true)
	out.SetChangedFunc(func() {
		out.ScrollToEnd()
		s.state.App.Draw()
	})
	out.SetDynamicColors(true)

	writer := tview.ANSIWriter(out)

	textViewElement := util.FlexElement{
		Primitive:    out,
		RelativeSize: 1,
		Focus:        true,
	}

	outputWindow := util.NewStatusFloatingWindow(s.state.Pages, "patch", true, textViewElement)
	outputWindow.Show()

	// Sadly libgit2 does not support adding hunks, only files. So in this one
	// case, we shell out.
	cmd := exec.Command("git", "add", "-p")
	cmd.Dir = git.BranchPath(s.state.Branch, s.state.Repository.Name)

	f, err := pty.Start(cmd)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Failed to start `git add -p` in PTY: %v", err))
		return
	}

	stop := make(chan struct{})

	go func() {
		r := bufio.NewReader(f)
		bs := make([]byte, 4096)
		for {
			select {
			case <-stop:
				return
			default:
				n, err := r.Read(bs)
				if err != nil && err != io.EOF {
					util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Failed to read `git add -p` output: %v", err))
					return
				}
				// In terminals, pressing Backspace does not result in a single
				// byte being written but three: `[Backspace] [Space]
				// [Backspace]`. See: https://unix.stackexchange.com/a/414246.
				filtered := strings.ReplaceAll(string(bs[:n]), "\x08 \x08", "")
				if _, err := writer.Write([]byte(filtered)); err != nil {
					util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Failed to write `git add -p` output to text view: %v", err))
					return
				}
			}
		}
	}()

	// We want to be able to backtrack on input, but not further than what
	// we've added, so we count each character that gets written.
	nInput := 0
	// We need to make sure we don't trigger backtracking multiple times at
	// once or at the same time as input.
	inputLock := sync.Mutex{}

	out.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		inputLock.Lock()
		defer inputLock.Unlock()
		switch event.Key() {
		case tcell.KeyEscape:
			if err := cmd.Process.Kill(); err != nil {
				util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Failed to run `git add -p`: %v", err))
				return event
			}
		case tcell.KeyEnter:
			// Once confirming, there is no more input to backspace over.
			nInput = 0
			if _, err := f.Write([]byte("\n")); err != nil {
				util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Failed to write to text view: %v", err))
				return event
			}

		// For some reason on my keyboard Backspace2 is emitted.
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if nInput > 0 {
				// Leaving the style tags in, since we'll rewrite the modified text
				// back into the view.
				text := out.GetText(false)
				if _, err := f.Write([]byte{byte(event.Key())}); err != nil {
					util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Failed to write Backspace to `git add -p` stdin: %v", err))
					return event
				}
				text = text[:len(text)-1]
				nInput--
				out.SetText(text)
			}
		default:
			// Ignore anything can't be displayed.
			if !unicode.IsPrint(event.Rune()) {
				break
			}
			if _, err := f.Write([]byte{byte(event.Rune())}); err != nil {
				util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Failed to write to `git add -p` stdin: %v", err))
				return event
			}
			nInput++
		}
		return event
	})

	go func() {
		if err := cmd.Wait(); err != nil {
			util.NewErrorWindow(s.state.Pages, "status-patch-err", fmt.Errorf("Running `git add -p` failed: %v", err))
		}
		stop <- struct{}{}
		outputWindow.Hide()
	}()
}

func (s *statusPage) restore() {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textViewElement := util.FlexElement{
		Primitive:    textView,
		RelativeSize: 1,
	}

	inputField := tview.NewInputField().SetPlaceholder("files:")
	inputFieldElement := util.FlexElement{
		Primitive: inputField,
		FixedSize: 1,
		Focus:     true,
	}

	restoreWindow := util.NewStatusFloatingWindow(s.state.Pages, "restore", true, textViewElement, inputFieldElement)
	restoreWindow.Show()

	status, err := s.state.Repository.Status(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-restore-err", fmt.Errorf("Failed to get status: %v", err))
		return
	}

	textView.SetText(strings.Join(status.ChangedFilesPaths(), "\n"))

	inputField.SetChangedFunc(func(text string) {
		changedFilesPaths := status.ChangedFilesPaths()
		if slices.Contains(changedFilesPaths, text) {
			textView.SetText(text)
			return
		}
		matched := fuzzy.Find(text, changedFilesPaths)
		textView.SetText(strings.Join(matched, "\n"))
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if err := s.state.Repository.Restore(
				s.state.Branch,
				strings.Split(textView.GetText(true), "\n"),
			); err != nil {
				util.NewErrorWindow(s.state.Pages, "status-restore-err", fmt.Errorf("Failed to restore files: %v", err))
				return
			}
			restoreWindow.Hide()
		}
	})
}

// TODO: this is almost identical to `add` above.
func (s *statusPage) unstage() {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textViewElement := util.FlexElement{
		Primitive:    textView,
		RelativeSize: 1,
	}

	// TODO: create a form, with checkboxes for every file?
	inputField := tview.NewInputField().SetPlaceholder("files:")
	inputFieldElement := util.FlexElement{
		Primitive: inputField,
		FixedSize: 1,
		Focus:     true,
	}

	unstageWindow := util.NewStatusFloatingWindow(s.state.Pages, "unstage", true, textViewElement, inputFieldElement)
	unstageWindow.Show()

	status, err := s.state.Repository.Status(s.state.Branch)
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-unstage-err", fmt.Errorf("Failed to get status: %v", err))
		return
	}

	textView.SetText(strings.Join(status.StagedFilesPaths(), "\n"))

	inputField.SetChangedFunc(func(text string) {
		stagedFilesPaths := status.StagedFilesPaths()
		if slices.Contains(stagedFilesPaths, text) {
			textView.SetText(text)
			return
		}
		matched := fuzzy.Find(text, stagedFilesPaths)
		textView.SetText(strings.Join(matched, "\n"))
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if err := s.state.Repository.Unstage(
				s.state.Branch,
				strings.Split(textView.GetText(true), "\n"),
			); err != nil {
				util.NewErrorWindow(s.state.Pages, "status-unstage-err", fmt.Errorf("Failed to unstage files: %v", err))
				return
			}
			unstageWindow.Hide()
		}
	})
}
