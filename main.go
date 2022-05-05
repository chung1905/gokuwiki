package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func getRepoDir() string {
	return "data/repo/"
}

func getPageDirName() string {
	return "pages"
}

func getPagesDir() string {
	return getRepoDir() + getPageDirName()
}

func isAllow(path string, d fs.DirEntry) bool {
	// Ignore directories and dot files
	if d.IsDir() || path[0:1] == "." || d.Name()[0:1] == "." {
		return false
	}

	return true
}

func homepage(c *gin.Context) {
	var pages []string
	dataDir := getPagesDir()
	dataDirLen := len(dataDir)

	err := filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		if !isAllow(path[dataDirLen:], d) {
			return nil
		}

		pages = append(pages, path[len(dataDir)+1:])
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"pages": pages,
	})
}

func viewWiki(c *gin.Context) {
	page := c.Param("page")
	file := getPagesDir() + page
	wikiContent, err := os.ReadFile(file)

	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}

	extensions := parser.CommonExtensions | parser.HardLineBreak | parser.FencedCode
	parserModel := parser.NewWithExtensions(extensions)
	wikiContent = markdown.NormalizeNewlines(wikiContent)
	output := markdown.ToHTML(wikiContent, parserModel, nil)

	c.HTML(http.StatusOK, "wiki.html", gin.H{
		"title":       page,
		"wikiContent": template.HTML(output),
	})
}

func editWiki(c *gin.Context) {
	page := c.Param("page")
	file := getPagesDir() + page
	wikiContent, _ := os.ReadFile(file)

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"title":       page,
		"page":        page,
		"wikiContent": string(wikiContent),
	})
}

func commitFile(page string) {
	repo, err := git.PlainOpen(getRepoDir())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = worktree.Add(getPageDirName() + page)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = worktree.Commit("Update "+getPageDirName()+page, getGitCommitOptions())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func deleteWiki(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func saveWiki(c *gin.Context) {
	page := c.PostForm("page")
	wikiContent := c.PostForm("content")
	wikiContentBytes := markdown.NormalizeNewlines([]byte(wikiContent))

	filename := getPagesDir() + page

	if len(wikiContentBytes) == 0 {
		deleteWiki(filename)
		go commitFile(page)
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = file.Write(wikiContentBytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = file.Sync()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go commitFile(page)
	c.Redirect(http.StatusSeeOther, "wiki/"+page)
}

func getRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/pub", "./pub") // use the loaded source
	router.GET("/", homepage)
	router.GET("/wiki/*page", viewWiki)
	router.GET("/edit/*page", editWiki)
	router.POST("/submitWiki", saveWiki)

	return router
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

func prepareGitRepo(repoDir string) {
	repo := initRepoIfNotExist(repoDir)
	commitOldData(repo)
}

func prepareDirectory(pagesDir string) {
	err := os.MkdirAll(pagesDir, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	prepareDirectory(getPagesDir())
	prepareGitRepo(getRepoDir())
	router := getRouter()
	err := router.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
