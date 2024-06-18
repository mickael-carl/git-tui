package git

import (
	"fmt"
	"os"
	"path"

	git "github.com/libgit2/git2go/v36"
)

func BranchPath(name, repo string) string {
	return path.Join(repoPath(repo), name)
}

func getBranches(name, mainBranch string) ([]Branch, error) {
	branches := []Branch{}

	mainBranchDir := BranchPath(mainBranch, name)

	repo, err := git.OpenRepository(mainBranchDir)
	if err != nil {
		return branches, err
	}
	defer repo.Free()

	branchesName, err := repo.Worktrees.List()
	if err != nil {
		return branches, err
	}
	branchesName = append(branchesName, mainBranch)

	for _, branch := range branchesName {
		head, err := repo.Head()
		if err != nil {
			return branches, err
		}
		defer head.Free()

		upstreamRef, err := head.Branch().Upstream()
		if err != nil {
			return branches, err
		}
		defer upstreamRef.Free()

		localOID := head.Target()
		upstreamOID := upstreamRef.Target()
		ahead, behind, err := repo.AheadBehind(localOID, upstreamOID)
		if err != nil {
			return branches, err
		}

		lastCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return branches, err
		}
		defer lastCommit.Free()

		branches = append(branches, Branch{
			Name:          branch,
			CommitsAhead:  ahead,
			CommitsBehind: behind,
			LastCommit:    lastCommit.Author().When,
		})
	}

	return branches, nil
}

func (r *Repository) CreateBranch(name string) error {
	// TODO: make a convenience method to open the main tree.
	mainBranchDir := BranchPath(r.Config.MainBranch, r.Name)

	repo, err := git.OpenRepository(mainBranchDir)
	if err != nil {
		return fmt.Errorf("couldn't open repository at %s: %v", mainBranchDir, err)
	}
	defer repo.Free()

	if _, err := repo.Worktrees.Add(
		name,
		BranchPath(name, r.Name),
		nil,
	); err != nil {
		return fmt.Errorf("couldn't create worktree: %v", err)
	}

	branches, err := getBranches(r.Name, r.Config.MainBranch)
	if err != nil {
		return fmt.Errorf("couldn't get branches: %v", err)
	}

	r.Branches = branches

	return nil
}

func (r *Repository) DeleteBranch(name string) error {
	if name == r.Config.MainBranch {
		return fmt.Errorf("cannot delete main branch %s for repository %s", name, r.Name)
	}

	branchDir := BranchPath(name, r.Name)
	// TODO: this needs some form of confirmation, i.e. force writing in some
	// input field the name of the branch to delete.
	if err := os.RemoveAll(branchDir); err != nil {
		return err
	}

	mainBranchDir := BranchPath(r.Config.MainBranch, r.Name)
	repo, err := git.OpenRepository(mainBranchDir)
	if err != nil {
		return err
	}
	defer repo.Free()

	branch, err := repo.LookupBranch(name, git.BranchLocal)
	if err != nil {
		return err
	}
	defer branch.Free()

	if err := branch.Delete(); err != nil {
		return err
	}

	worktree, err := repo.Worktrees.Lookup(name)
	if err != nil {
		return err
	}
	defer worktree.Free()

	if err := worktree.Prune(git.WorktreePruneWorkingTree); err != nil {
		return err
	}

	branches, err := getBranches(r.Name, r.Config.MainBranch)
	if err != nil {
		return err
	}

	r.Branches = branches

	return nil
}
