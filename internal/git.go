package internal

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"time"
)

func CommitFile(filepath string, repoDir string, editComment string) {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = worktree.Add(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = worktree.Commit(filepath+": "+editComment, getGitCommitOptions())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func initRepoIfNotExist(repoDir string) *git.Repository {
	repo, err := git.PlainOpenWithOptions(repoDir, &git.PlainOpenOptions{DetectDotGit: false})
	if errors.Is(err, git.ErrRepositoryNotExists) {
		repo, _ = git.PlainInit(repoDir, false)
	}

	return repo
}

func getGitCommitOptions() *git.CommitOptions {
	return &git.CommitOptions{
		Author: &object.Signature{
			Name:  "gokuwiki web",
			Email: "gokuwiki+web@chungn.com",
			When:  time.Now(),
		},
	}
}

func commitOldData(repo *git.Repository) {
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = worktree.Add(".")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = worktree.Commit("Initial commit ", getGitCommitOptions())
}

func PrepareGitRepo(repoDir string) {
	repo := initRepoIfNotExist(repoDir)
	commitOldData(repo)
}
