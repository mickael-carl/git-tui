package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"

	fzf "github.com/junegunn/fzf/src"

	"github.com/mickael-carl/git-tui/pkg/git"
)

func pickRepo(repos []git.Repository) (string, []git.Branch, error) {
	reposMap := map[string][]git.Branch{}

	input := make(chan string, len(repos))
	for _, repo := range repos {
		input <- repo.Name
		reposMap[repo.Name] = repo.Branches
	}
	close(input)

	output := make(chan string, 1)

	options, err := fzf.ParseOptions(true, []string{})
	if err != nil {
		log.Fatal(err)
	}
	options.Input = input
	options.Output = output

	_, err = fzf.Run(options)
	if err != nil {
		log.Fatal(err)
	}

	name := <-output

	return name, reposMap[name], nil
}

func pickBranch(branches []git.Branch) (string, error) {
	input := make(chan string, len(branches))
	for _, branch := range branches {
		input <- branch.Name
	}
	close(input)

	output := make(chan string, 1)

	options, err := fzf.ParseOptions(true, []string{})
	if err != nil {
		log.Fatal(err)
	}
	options.Input = input
	options.Output = output

	_, err = fzf.Run(options)
	if err != nil {
		log.Fatal(err)
	}

	name := <-output

	return name, nil
}

func main() {
	repos, err := git.GetRepositories()
	if err != nil {
		log.Fatal(err)
	}

	repo, branches, err := pickRepo(repos)
	if err != nil {
		log.Fatal(err)
	}

	branch, err := pickBranch(branches)
	if err != nil {
		log.Fatal(err)
	}

	fishPath, err := exec.LookPath("fish")
	if err != nil {
		log.Fatal(err)
	}

	_ = syscall.Exec(
		fishPath,
		[]string{
			fishPath,
			"-iC",
			// TODO: make this configurable too.
			"cd " + path.Join(os.Getenv("HOME"), "dev", repo, branch),
		},
		os.Environ(),
	)
}
