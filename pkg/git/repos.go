package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	git "github.com/libgit2/git2go/v36"
)

// TODO: move this to some configuration file.
var DevDir = path.Join(os.Getenv("HOME"), "dev")

// TODO: does this need to be public?
type Remote struct {
	URL  string
	Name string
}

type Config struct {
	MainBranch     string
	Remote         Remote
	PubKeyPath     string
	PrivateKeyPath string
}

type Branch struct {
	Name          string
	CommitsAhead  int
	CommitsBehind int
	LastCommit    time.Time
}

type Repository struct {
	Name     string
	Config   Config
	Branches []Branch
}

func NewRepository(name, url, mainBranch, pubKeyPath, privateKeyPath string) error {
	config := Config{
		MainBranch: mainBranch,
		Remote: Remote{
			Name: "origin",
			URL:  url,
		},
		PubKeyPath:     pubKeyPath,
		PrivateKeyPath: privateKeyPath,
	}

	_, err := os.Open(repoPath(name))
	if err == nil {
		return fmt.Errorf("repository %s already exists", name)
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(repoPath(name), 0o755); err != nil {
		return err
	}

	json, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(repoPath(name), "config.json"), json, 0o644); err != nil {
		return err
	}

	repo := Repository{
		Config: config,
	}

	remoteCallbacks := git.RemoteCallbacks{
		CredentialsCallback: repo.credentialsCallback,
	}

	options := git.FetchOptions{
		RemoteCallbacks: remoteCallbacks,
	}

	cloneOptions := git.CloneOptions{
		FetchOptions:    options,
		CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutForce},
	}

	// TODO: add a progress callback for feedback in the UI.
	gitRepo, err := git.Clone(url, BranchPath(mainBranch, name), &cloneOptions)
	if err != nil {
		return err
	}
	defer gitRepo.Free()

	return nil
}

func DeleteRepository(name string) error {
	return os.RemoveAll(repoPath(name))
}

func GetRepositories() ([]Repository, error) {
	files, err := os.ReadDir(DevDir)
	if err != nil {
		return []Repository{}, err
	}

	repos := []Repository{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		data, err := os.ReadFile(filepath.Join(DevDir, file.Name(), "config.json"))
		if err != nil {
			// Ignore repositories not managed by this CLI.
			if errors.Is(err, os.ErrNotExist) {
				continue
			} else {
				return []Repository{}, err
			}
		}

		var config Config
		if err := json.Unmarshal(data, &config); err != nil {
			return []Repository{}, err
		}

		branches, err := getBranches(file.Name(), config.MainBranch)
		if err != nil {
			return []Repository{}, err
		}

		repos = append(repos, Repository{
			Name:     file.Name(),
			Config:   config,
			Branches: branches,
		})
	}

	return repos, nil
}

func repoPath(name string) string {
	return path.Join(DevDir, name)
}
