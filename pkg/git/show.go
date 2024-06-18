package git

import (
	"errors"

	git "github.com/libgit2/git2go/v36"
)

func (r *Repository) Show(branch string) ([]Diff, error) {
	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return []Diff{}, err
	}
	defer repo.Free()

	head, err := repo.Head()
	if err != nil {
		return []Diff{}, err
	}
	defer head.Free()

	commit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return []Diff{}, err
	}
	defer commit.Free()

	if commit.ParentCount() == 0 {
		return []Diff{}, errors.New("no previous commit to show")
	}

	parent := commit.Parent(0)
	defer parent.Free()

	diffOpts, err := git.DefaultDiffOptions()
	if err != nil {
		return []Diff{}, err
	}

	parentTree, err := parent.Tree()
	if err != nil {
		return []Diff{}, err
	}
	defer parentTree.Free()

	commitTree, err := commit.Tree()
	if err != nil {
		return []Diff{}, err
	}
	defer commitTree.Free()

	diff, err := repo.DiffTreeToTree(parentTree, commitTree, &diffOpts)
	if err != nil {
		return []Diff{}, err
	}
	defer diff.Free()

	// TODO: the code below is identical to the one for Diff.
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
