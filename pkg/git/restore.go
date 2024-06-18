package git

import (
	"os"

	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Restore(branch string, paths []string) error {
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

	tree, err := commit.Tree()
	if err != nil {
		return err
	}
	defer tree.Free()

	for _, path := range paths {
		entry, err := tree.EntryByPath(path)
		if err != nil {
			return err
		}

		blob, err := repo.LookupBlob(entry.Id)
		if err != nil {
			return err
		}
		defer blob.Free()

		if err := os.WriteFile(path, blob.Contents(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
