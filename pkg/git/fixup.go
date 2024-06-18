package git

import (
	"fmt"

	git "github.com/libgit2/git2go/v36"
)

// This is useful for 2 reasons:
//  1. We can have a stable iteration in the right order over the list of
//     commits.
//  2. We can grab the first commit from the list and use that to set the text
//     of the message view on first draw, since tview is misbehaving on that
//     front.
type FixupCandidates struct {
	Hashes  []string
	Commits map[string]Commit
}

func (r *Repository) FixupList(branch string) (FixupCandidates, error) {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return FixupCandidates{}, err
	}
	defer repo.Free()

	head, err := repo.Head()
	if err != nil {
		return FixupCandidates{}, err
	}
	defer head.Free()

	commit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return FixupCandidates{}, err
	}
	defer commit.Free()

	mainBranch, err := repo.LookupBranch(r.Config.MainBranch, git.BranchLocal)
	if err != nil {
		return FixupCandidates{}, err
	}
	defer mainBranch.Free()

	mainBranchCommit, err := repo.LookupCommit(mainBranch.Reference.Target())
	if err != nil {
		return FixupCandidates{}, err
	}
	defer mainBranchCommit.Free()

	revWalk, err := repo.Walk()
	if err != nil {
		return FixupCandidates{}, err
	}
	defer revWalk.Free()

	if err = revWalk.PushRange(fmt.Sprintf("%s..%s", mainBranchCommit.Id(), commit.Id())); err != nil {
		return FixupCandidates{}, err
	}

	commits := FixupCandidates{
		Hashes:  []string{},
		Commits: map[string]Commit{},
	}

	err = revWalk.Iterate(func(commit *git.Commit) bool {
		commits.Hashes = append(commits.Hashes, commit.Id().String())
		commits.Commits[commit.Id().String()] = toTUICommit(commit)
		return true
	})

	return commits, err
}

func (r *Repository) Fixup(branch, commit string) error {
	return r.Commit(branch, fmt.Sprintf("fixup! %s", commit), commit)
}
