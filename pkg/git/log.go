package git

import (
	"fmt"
	"time"

	git "github.com/libgit2/git2go/v36"
)

type Commit struct {
	Author  string
	Date    time.Time
	Hash    string
	Message string
}

func (r *Repository) Log(branch string) ([]Commit, error) {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return []Commit{}, err
	}
	defer repo.Free()

	head, err := repo.Head()
	if err != nil {
		return []Commit{}, err
	}
	defer head.Free()

	commit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return []Commit{}, err
	}
	defer commit.Free()

	revWalk, err := repo.Walk()
	if err != nil {
		return []Commit{}, err
	}
	defer revWalk.Free()

	if err := revWalk.Push(commit.Object.Id()); err != nil {
		return []Commit{}, err
	}

	revWalk.Sorting(git.SortTopological)

	items := []Commit{}

	if err = revWalk.Iterate(func(commit *git.Commit) bool {
		items = append(items, toTUICommit(commit))
		return true
	}); err != nil {
		return []Commit{}, err
	}

	return items, nil
}

func toTUICommit(commit *git.Commit) Commit {
	return Commit{
		Author:  fmt.Sprintf("%s <%s>", commit.Author().Name, commit.Author().Email),
		Date:    commit.Author().When,
		Hash:    commit.Id().String(),
		Message: commit.Message(),
	}
}
