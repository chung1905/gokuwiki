package main

import (
	"chungn/gokuwiki/internal"
	"chungn/gokuwiki/internal/captcha"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
)

func getRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*/*.html")
	router.Static("/pub", "./web/pub")
	router.GET("/", homepage)
	router.GET("/wiki/*page", viewWiki)
	router.GET("/edit/*page", editWiki)
	router.POST("/submitWiki", saveWiki)

	return router
}

func homepage(c *gin.Context) {
	c.HTML(http.StatusOK, "list.html", gin.H{
		"title": "Wiki",
		"pages": internal.ListFiles(getPagesDir()),
	})
}

func viewWiki(c *gin.Context) {
	page := c.Param("page")
	file := getPagesDir() + page
	wikiContent, err := internal.ReadFile(file)
	var buttonText string

	if errors.Is(err, fs.ErrNotExist) {
		wikiContent = ([]byte)("Empty Page")
		buttonText = "Create"
	} else {
		buttonText = "Edit"
	}

	output := internal.Md2html(wikiContent)

	c.HTML(http.StatusOK, "wiki.html", gin.H{
		"page":        page,
		"title":       page,
		"wikiContent": template.HTML(output),
		"buttonText":  buttonText,
	})
}

func editWiki(c *gin.Context) {
	page := c.Param("page")
	file := getPagesDir() + page
	wikiContent, _ := internal.ReadFile(file)

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"title":            page,
		"page":             page,
		"wikiContent":      string(wikiContent),
		"turnstileEnabled": getTurnstileEnabled(),
		"turnstileSiteKey": getTurnstileSiteKey(),
	})
}

func saveWiki(c *gin.Context) {
	if getTurnstileEnabled() {
		captchaResult := captcha.Validate(c.PostForm("cf-turnstile-response"), getTurnstileSecretKey())
		if !captchaResult {
			return
		}
	}

	page := c.PostForm("page")
	if page[0:1] != "/" {
		page = "/" + page
	}

	editComment := c.PostForm("comment")
	if len(editComment) == 0 {
		c.Redirect(http.StatusSeeOther, "wiki/"+page)
		return
	}

	wikiContent := c.PostForm("content")
	wikiContentBytes := internal.NormalizeNewlines([]byte(wikiContent))

	filepath := getPagesDir() + page

	if len(wikiContentBytes) == 0 {
		internal.DeleteFile(filepath)
		go internal.CommitFile(getPageDirName()+page, getRepoDir(), editComment, getGitAccessToken())
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	internal.SaveFile(wikiContentBytes, filepath)
	go internal.CommitFile(getPageDirName()+page, getRepoDir(), editComment, getGitAccessToken())

	c.Redirect(http.StatusSeeOther, "wiki/"+page)
}

func main() {
	internal.CreateDir(getPagesDir())
	internal.PrepareGitRepo(getRepoDir(), getRepoURL())
	router := getRouter()
	err := router.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
