package git

import (
	"errors"

	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) ResetPreviousCommit(branch string) error {
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

	if commit.ParentCount() == 0 {
		return errors.New("no parent commit for HEAD")
	}

	return repo.ResetToCommit(commit.Parent(0), git.ResetMixed, nil)
}
