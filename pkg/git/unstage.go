package git

import (
	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Unstage(branch string, paths []string) error {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return err
	}
	defer repo.Free()

	head, err := repo.Head()
	if err != nil {
		return err
	}
	defer head.Free()

	commit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return err
	}
	defer commit.Free()

	return repo.ResetDefaultToCommit(commit, paths)
}
