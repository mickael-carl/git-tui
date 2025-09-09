package git

import (
	"fmt"

	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Push(branch string, done func(err error)) {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		done(err)
		return
	}
	defer repo.Free()

	// TODO: for now this assumes that the remote is called `origin` which
	// won't always be true. That's because the Repository struct doesn't
	// really contain information about remotes.
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		done(err)
		return
	}
	defer remote.Free()

	remoteCallbacks := git.RemoteCallbacks{
		CredentialsCallback: r.credentialsCallback,
	}

	pushOptions := &git.PushOptions{
		RemoteCallbacks: remoteCallbacks,
	}

	err = remote.Push(
		[]string{fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)},
		pushOptions,
	)
	done(err)
	return
}
