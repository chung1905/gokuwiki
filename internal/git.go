package internal

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func PrepareGitRepo(repoDir string, remoteUrl string, accessToken string) {
	repo := initRepoIfNotExist(repoDir)
	if len(remoteUrl) > 0 {
		addRemote(repo, remoteUrl)
		if len(accessToken) > 0 {
			pull(repo, accessToken)
		}
	}
	commitOldData(repo, accessToken)
}

func CommitFiles(filepaths []string, repoDir string, editComment string, accessToken string) {
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

	for _, filepath := range filepaths {
		_, err = worktree.Add(filepath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	// Create commit message based on files
	var commitMsg string
	if len(filepaths) == 1 {
		commitMsg = filepaths[0] + ": " + editComment
	} else {
		commitMsg = fmt.Sprintf("Modified %d files: %s", len(filepaths), editComment)
	}

	_, err = worktree.Commit(commitMsg, getGitCommitOptions())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(accessToken) > 0 {
		push(repo, accessToken)
	}
}

func CommitFile(filepath string, repoDir string, editComment string, accessToken string) {
	CommitFiles([]string{filepath}, repoDir, editComment, accessToken)
}

func push(repo *git.Repository, accessToken string) {
	err := repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: "gokuwiki",
			Password: accessToken,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func pull(repo *git.Repository, accessToken string) {
	worktree, err := repo.Worktree()
	if err == nil {
		err = worktree.Pull(&git.PullOptions{
			Auth: &http.BasicAuth{
				Username: "gokuwiki",
				Password: accessToken,
			},
		})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			fmt.Println("Pull error:", err.Error())
		}
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

func commitOldData(repo *git.Repository, accessToken string) {
	worktree, err := repo.Worktree()
	if err != nil {
		log.Println(err.Error())
		return
	}

	status, _ := worktree.Status()
	if status.IsClean() {
		return
	}

	_, err = worktree.Add(".")
	if err != nil {
		log.Println(err.Error())
		return
	}

	_, _ = worktree.Commit("Commit unstaged files", getGitCommitOptions())
	if len(accessToken) > 0 {
		push(repo, accessToken)
	}
}

func addRemote(repo *git.Repository, url string) {
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	if errors.Is(err, git.ErrRemoteExists) {
		log.Println("Remote \"origin\" already exists")
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
