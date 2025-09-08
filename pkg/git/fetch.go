package git

import (
	"path"

	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Fetch(done func(err error)) {
	mainBranchDir := path.Join(repoPath(r.Name), r.Config.MainBranch)

	repo, err := git.OpenRepository(mainBranchDir)
	if err != nil {
		// TODO: defer done(err) doesn't seem to work, so we have to call it
		// before returning everytime.
		done(err)
		return
	}
	defer repo.Free()

	remotes, err := repo.Remotes.List()
	if err != nil {
		done(err)
		return
	}

	for _, remoteName := range remotes {
		remote, err := repo.Remotes.Lookup(remoteName)
		if err != nil {
			done(err)
			return
		}
		defer remote.Free()

		remoteCallbacks := git.RemoteCallbacks{
			CredentialsCallback: r.credentialsCallback,
		}

		options := &git.FetchOptions{
			RemoteCallbacks: remoteCallbacks,
		}

		err = remote.Fetch([]string{}, options, "")
		if err != nil {
			done(err)
			return
		}
	}

	done(nil)
	return
}

func (r *Repository) credentialsCallback(url, usernameFromURL string, allowedTypes git.CredentialType) (*git.Credential, error) {
	return git.NewCredentialSSHKey(
		"git",
		r.Config.PubKeyPath,
		r.Config.PrivateKeyPath,
		"",
	)
}
