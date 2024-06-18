package git

import (
	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Amend(branch, message string) error {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return err
	}
	defer repo.Free()

	tree, err := commitTree(repo, branch)
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}
	defer head.Free()

	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return err
	}
	defer headCommit.Free()

	if _, err := headCommit.Amend("HEAD", nil, nil, message, tree); err != nil {
		return err
	}

	return nil
}
