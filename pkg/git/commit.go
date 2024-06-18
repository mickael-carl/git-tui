package git

import (
	"time"

	git "github.com/libgit2/git2go/v36"
)

func commitTree(repo *git.Repository, branch string) (*git.Tree, error) {
	index, err := repo.Index()
	if err != nil {
		return nil, err
	}
	defer index.Free()

	treeID, err := index.WriteTree()
	if err != nil {
		return nil, err
	}

	return repo.LookupTree(treeID)
}

func (r *Repository) Commit(branch, message, parent string) error {
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
	defer tree.Free()

	var parentCommit *git.Commit
	if parent == "" {
		head, err := repo.Head()
		if err != nil {
			return err
		}
		defer head.Free()

		parentCommit, err = repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}
	} else {
		oid, err := git.NewOid(parent)
		if err != nil {
			return err
		}

		parentCommit, err = repo.LookupCommit(oid)
		if err != nil {
			return err
		}
	}
	defer parentCommit.Free()

	// TODO: read from gitconfig.
	signature := &git.Signature{
		Name:  "Mickael Carl",
		Email: "contact@carlm.org",
		When:  time.Now(),
	}

	if _, err := repo.CreateCommit("HEAD", signature, signature, message, tree, parentCommit); err != nil {
		return err
	}

	return nil
}

func (r *Repository) HeadCommitMessage(branch string) (string, error) {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return "", err
	}
	defer repo.Free()

	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	defer head.Free()

	commit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return "", err
	}
	defer commit.Free()

	return commit.Message(), nil
}
