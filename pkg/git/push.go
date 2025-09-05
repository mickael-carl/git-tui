package git

import (
	"fmt"

	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Push(branch string) error {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return err
	}
	defer repo.Free()

	// TODO: for now this assumes that the remote is called `origin` which
	// won't always be true. That's because the Repository struct doesn't
	// really contain information about remotes.
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}
	defer remote.Free()

	remoteCallbacks := git.RemoteCallbacks{
		CredentialsCallback: r.credentialsCallback,
	}

	pushOptions := &git.PushOptions{
		RemoteCallbacks: remoteCallbacks,
	}

	return remote.Push(
		[]string{fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)},
		pushOptions,
	)
}
