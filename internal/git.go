package internal

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"time"
)

func PrepareGitRepo(repoDir string, remoteUrl string) {
	repo := initRepoIfNotExist(repoDir)
	commitOldData(repo)
	if len(remoteUrl) > 0 {
		addRemote(repo, remoteUrl)
	}
}

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

	accessToken := os.Getenv("GOKUWIKI_ACCESS_TOKEN")
	if len(accessToken) > 0 {
		push(repo, accessToken)
	}
}

func push(repo *git.Repository, accessToken string) {
	err := repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "gokuwiki",
			Password: accessToken,
		},
		InsecureSkipTLS: true,
	})
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

func addRemote(repo *git.Repository, url string) {
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	if errors.Is(err, git.ErrRemoteExists) {
		fmt.Println("Remote \"origin\" already exists")
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
