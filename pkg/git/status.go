package git

import (
	"fmt"

	git "github.com/libgit2/git2go/v36"
)

type BaseStatus struct {
	Change  git.Status
	Path    string
	NewPath string
}

type Status struct {
	Staged    []BaseStatus
	Unstaged  []BaseStatus
	Untracked []string
}

var gitStatusToSymbol map[git.Status]string = map[git.Status]string{
	git.StatusIndexNew:        "🆕",
	git.StatusIndexModified:   "📝",
	git.StatusIndexDeleted:    "💥",
	git.StatusIndexRenamed:    "➡️",
	git.StatusIndexTypeChange: "🛠️",
	git.StatusWtNew:           "🆕",
	git.StatusWtModified:      "📝",
	git.StatusWtDeleted:       "💥",
	git.StatusWtRenamed:       "➡️",
	git.StatusWtTypeChange:    "🛠️",
}

func (s *Status) ChangedFilesPaths() []string {
	paths := s.Untracked
	for _, unstagedChange := range s.Unstaged {
		if unstagedChange.NewPath != "" {
			paths = append(paths, unstagedChange.NewPath)
		} else {
			paths = append(paths, unstagedChange.Path)
		}
	}
	return paths
}

func (s *Status) StagedFilesPaths() []string {
	paths := []string{}
	for _, stagedChange := range s.Staged {
		if stagedChange.NewPath != "" {
			paths = append(paths, stagedChange.NewPath)
		} else {
			paths = append(paths, stagedChange.Path)
		}
	}
	return paths
}

func (s *Status) FormatStaged() string {
	// TODO: use strings.Builder here.
	out := ""
	for _, stagedChange := range s.Staged {
		out += fmt.Sprintf(
			"%s %s",
			gitStatusToSymbol[stagedChange.Change],
			stagedChange.Path,
		)
		if stagedChange.NewPath != "" {
			out += fmt.Sprintf(" -> %s", stagedChange.NewPath)
		}
		out += "\n"
	}
	return out
}

// TODO: this is basically identical to `FormatStaged`.
func (s *Status) FormatUnstaged() string {
	// TODO: use strings.Builder here.
	out := ""
	for _, unstagedChange := range s.Unstaged {
		out += fmt.Sprintf(
			"%s %s",
			gitStatusToSymbol[unstagedChange.Change],
			unstagedChange.Path,
		)
		if unstagedChange.NewPath != "" {
			out += fmt.Sprintf(" -> %s", unstagedChange.NewPath)
		}
		out += "\n"
	}
	return out
}

func (s *Status) FormatUntracked() string {
	out := ""
	for _, untrackedChange := range s.Untracked {
		out += fmt.Sprintf("🆕 %s\n", untrackedChange)
	}
	return out
}

func (r *Repository) Status(branch string) (Status, error) {
	staged := []BaseStatus{}
	unstaged := []BaseStatus{}
	untracked := []string{}

	branchDir := BranchPath(branch, r.Name)

	repo, err := git.OpenRepository(branchDir)
	if err != nil {
		return Status{}, err
	}
	defer repo.Free()

	workdirStatusList, err := repo.StatusList(
		&git.StatusOptions{
			Show:  git.StatusShowWorkdirOnly,
			Flags: git.StatusOptIncludeUntracked,
		},
	)
	if err != nil {
		return Status{}, err
	}
	defer workdirStatusList.Free()

	workdirStatusListCount, err := workdirStatusList.EntryCount()
	if err != nil {
		return Status{}, err
	}

	for i := range workdirStatusListCount {
		entry, err := workdirStatusList.ByIndex(i)
		if err != nil {
			return Status{}, err
		}

		switch entry.Status {
		case git.StatusWtNew:
			path := entry.IndexToWorkdir.NewFile.Path
			untracked = append(untracked, path)
		case git.StatusWtRenamed:
			oldPath := entry.IndexToWorkdir.OldFile.Path
			newPath := entry.IndexToWorkdir.NewFile.Path
			unstaged = append(unstaged, BaseStatus{
				Change:  entry.Status,
				Path:    oldPath,
				NewPath: newPath,
			})
		default:
			unstaged = append(unstaged, BaseStatus{
				Change: entry.Status,
				Path:   entry.IndexToWorkdir.OldFile.Path,
			})
		}
	}

	indexStatusList, err := repo.StatusList(
		&git.StatusOptions{
			Show: git.StatusShowIndexOnly,
		},
	)
	if err != nil {
		return Status{}, err
	}
	defer indexStatusList.Free()

	indexStatusListCount, err := indexStatusList.EntryCount()
	if err != nil {
		return Status{}, err
	}

	for i := range indexStatusListCount {
		entry, err := indexStatusList.ByIndex(i)
		if err != nil {
			return Status{}, err
		}

		switch entry.Status {
		case git.StatusIndexNew:
			path := entry.HeadToIndex.NewFile.Path
			staged = append(staged, BaseStatus{
				Change: entry.Status,
				Path:   path,
			})
		case git.StatusWtRenamed:
			oldPath := entry.HeadToIndex.OldFile.Path
			newPath := entry.HeadToIndex.NewFile.Path
			staged = append(staged, BaseStatus{
				Change:  entry.Status,
				Path:    oldPath,
				NewPath: newPath,
			})
		default:
			staged = append(staged, BaseStatus{
				Change: entry.Status,
				Path:   entry.HeadToIndex.OldFile.Path,
			})
		}
	}

	return Status{
		Staged:    staged,
		Unstaged:  unstaged,
		Untracked: untracked,
	}, nil
}
