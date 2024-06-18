package git

import (
	"path"

	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Fetch() error {
	mainBranchDir := path.Join(repoPath(r.Name), r.Config.MainBranch)

	repo, err := git.OpenRepository(mainBranchDir)
	if err != nil {
		return err
	}
	defer repo.Free()

	remotes, err := repo.Remotes.List()
	if err != nil {
		return err
	}

	for _, remoteName := range remotes {
		remote, err := repo.Remotes.Lookup(remoteName)
		if err != nil {
			return err
		}
		defer remote.Free()

		remoteCallbacks := git.RemoteCallbacks{
			CredentialsCallback: credentialsCallback,
		}

		options := &git.FetchOptions{
			RemoteCallbacks: remoteCallbacks,
		}

		if err := remote.Fetch([]string{}, options, ""); err != nil {
			return err
		}
	}

	return nil
}

func credentialsCallback(url, usernameFromURL string, allowedTypes git.CredentialType) (*git.Credential, error) {
	// TODO: don't hardcode this.
	return git.NewCredentialSSHKey(
		"git",
		"/Users/mickaelcarl/.ssh/github.pub",
		"/Users/mickaelcarl/.ssh/github",
		"",
	)
}
