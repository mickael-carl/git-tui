package git

import (
	git "github.com/libgit2/git2go/v36"
)

type Diff struct {
	Path    string
	NewPath string
	Patch   string
}

func (r *Repository) Diff(branch string) ([]Diff, error) {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return []Diff{}, err
	}
	defer repo.Free()

	index, err := repo.Index()
	if err != nil {
		return []Diff{}, err
	}
	defer index.Free()

	diff, err := repo.DiffIndexToWorkdir(index, nil)
	if err != nil {
		return []Diff{}, err
	}
	defer diff.Free()

	deltas, err := diff.NumDeltas()
	if err != nil {
		return []Diff{}, err
	}

	// TODO: use make here.
	diffs := []Diff{}
	for i := range deltas {
		delta, err := diff.Delta(i)
		if err != nil {
			return []Diff{}, err
		}

		patch, err := diff.Patch(i)
		if err != nil {
			return []Diff{}, err
		}

		patchString, err := patch.String()
		if err != nil {
			return []Diff{}, err
		}

		diff := Diff{
			Patch: patchString,
			Path:  delta.OldFile.Path,
		}

		if delta.NewFile.Path != "" {
			diff.NewPath = delta.NewFile.Path
		}

		diffs = append(diffs, diff)
	}

	return diffs, nil
}

func (r *Repository) DiffStaged(branch string) ([]Diff, error) {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return []Diff{}, err
	}
	defer repo.Free()

	index, err := repo.Index()
	if err != nil {
		return []Diff{}, err
	}
	defer index.Free()

	head, err := repo.Head()
	if err != nil {
		return []Diff{}, err
	}
	defer head.Free()

	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return []Diff{}, err
	}
	defer headCommit.Free()

	tree, err := headCommit.Tree()
	if err != nil {
		return []Diff{}, err
	}

	diff, err := repo.DiffTreeToIndex(tree, index, nil)
	if err != nil {
		return []Diff{}, err
	}
	defer diff.Free()

	deltas, err := diff.NumDeltas()
	if err != nil {
		return []Diff{}, err
	}

	// TODO: use make here.
	diffs := []Diff{}
	for i := range deltas {
		delta, err := diff.Delta(i)
		if err != nil {
			return []Diff{}, err
		}

		patch, err := diff.Patch(i)
		if err != nil {
			return []Diff{}, err
		}

		patchString, err := patch.String()
		if err != nil {
			return []Diff{}, err
		}

		diff := Diff{
			Patch: patchString,
			Path:  delta.OldFile.Path,
		}

		if delta.NewFile.Path != "" {
			diff.NewPath = delta.NewFile.Path
		}

		diffs = append(diffs, diff)
	}

	return diffs, nil
}
