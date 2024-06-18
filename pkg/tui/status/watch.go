package status

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/mickael-carl/git-tui/pkg/git"
	"github.com/mickael-carl/git-tui/pkg/tui/util"
)

func (s *statusPage) watchAndUpdate() {
	// Call it at least once so we're actually filling the page even if no
	// changes occur.
	s.update()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		util.NewErrorWindow(s.state.Pages, "status-watch-err", fmt.Errorf("Failed to create watcher: %v", err))
		return
	}

	watchedDirs := map[string]struct{}{}

	dir := git.BranchPath(s.state.Branch, s.state.Repository.Name)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && !strings.Contains(path, ".git") {
			if err := watcher.Add(path); err != nil {
				return err
			}
			watchedDirs[path] = struct{}{}
		}

		return nil
	}); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-watch-err", fmt.Errorf("Failed to walk directory %s: %v", dir, err))
		return
	}

	mainBranchDir := git.BranchPath(s.state.Repository.Config.MainBranch, s.state.Repository.Name)
	if err := watcher.Add(path.Join(mainBranchDir, ".git")); err != nil {
		util.NewErrorWindow(s.state.Pages, "status-watch-err", fmt.Errorf("Failed to add `.git` to watcher: %v", err))
		return
	}

	go func() {
		defer watcher.Close()
		for s.grid.HasFocus() {
			event := <-watcher.Events
			if strings.Contains(event.Name, ".git") && !strings.Contains(event.Name, "index.lock") {
				continue
			}

			if event.Has(fsnotify.Create) {
				// We purposefully ignore errors here. Some editors will copy
				// files around and move them so fast that the `Stat` call will
				// fail sometimes (e.g. vim).
				if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
					_ = watcher.Add(event.Name)
					watchedDirs[event.Name] = struct{}{}
				}
			}
			// Renames will actually also cause a `fsnotify.Create` so for the
			// actual `fsnotify.Rename` we only need to care about removing
			// files from the watcher.
			if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				// Similarly here, we remove files from being watched as a best
				// effort.
				if _, ok := watchedDirs[event.Name]; ok {
					_ = watcher.Remove(event.Name)
					delete(watchedDirs, event.Name)
				}
			}

			s.state.App.QueueUpdateDraw(s.update)
		}
	}()
}
