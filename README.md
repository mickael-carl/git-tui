# git-tui

An opinionated Git TUI, that makes use of
[worktrees](https://git-scm.com/docs/git-worktree) rather than plain branches
and provides keybindings for common actions.

Built using Go and [git2go](https://github.com/libgit2/git2go) (Go bindings for
libgit2) (mostly), on top
of [tview](https://github.com/rivo/tview).

Currently it clones and manages repository in `~/dev`. For each repository, a
directory is created, with a subdirectory per branch/worktree (including the
main branch).

## Dependencies

* For interactive patching, Git itself.
* On macOS, git2go/libgit2 depends on pkg-config (`brew install pkg-config`).

## Building & Installing

```
make install
```

## Why

Existing Git TUIs (e.g. [LazyGit](https://github.com/jesseduffield/lazygit)) are
very unopinionated, don't necessarily support worktrees and don't adapt too well
to my personal workflow. Also: I wanted to see if I could write a TUI in Go.

