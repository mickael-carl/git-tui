package git

import (
	git "github.com/libgit2/git2go/v36"
)

// TODO: make most methods actually apply to a Branch struct, that has a
// reference to the repository it belongs to. OR add the current branch name to
// the Repository struct?
func (r *Repository) Add(branch string, paths []string) error {
	// TODO: Make this a method of Repository that yields the open repo and a Free as
	// well so we don't need to set `branchDir` everytime.
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return err
	}
	defer repo.Free()

	index, err := repo.Index()
	if err != nil {
		return err
	}
	defer index.Free()

	if err := index.AddAll(paths, git.IndexAddDefault, nil); err != nil {
		return err
	}

	return index.Write()
}

func (r *Repository) AddAll(branch string) error {
	return r.Add(branch, []string{"."})
}
