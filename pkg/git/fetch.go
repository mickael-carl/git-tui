package git

import (
	"path"

	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Fetch(done func(err error)) error {
	var err error
	defer done(err)

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
			CredentialsCallback: r.credentialsCallback,
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

func (r *Repository) credentialsCallback(url, usernameFromURL string, allowedTypes git.CredentialType) (*git.Credential, error) {
	return git.NewCredentialSSHKey(
		"git",
		r.Config.PubKeyPath,
		r.Config.PrivateKeyPath,
		"",
	)
}
