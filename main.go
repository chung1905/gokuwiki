package main

import (
	"chungn/gokuwiki/internal"
	"chungn/gokuwiki/internal/captcha"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func getRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*/*.gohtml")
	router.Use(static.Serve("/", static.LocalFile("./web/pub", false)))
	router.GET("/", homepage)
	router.GET("/wiki/*page", viewWiki)
	router.GET("/edit/*page", editWiki)
	router.POST("/submitWiki", saveWiki)

	return router
}

func homepage(c *gin.Context) {
	c.HTML(http.StatusOK, "list.gohtml", gin.H{
		"title":  "Wiki",
		"pages":  internal.ListFiles(getPagesDir()),
		"result": internal.GetMessage(c.Query("m")),
	})
}

func viewWiki(c *gin.Context) {
	page := c.Param("page")
	file := getPagesDir() + page
	wikiContent, err := internal.ReadFile(file)
	var buttonText string
	var lastModifiedTime string

	if errors.Is(err, fs.ErrNotExist) {
		wikiContent = ([]byte)("Empty Page")
		buttonText = "Create"
	} else {
		buttonText = "Edit"
		fileStat, _ := os.Stat(file)
		lastModifiedTime = fileStat.ModTime().Format(time.UnixDate)
	}

	output := internal.Md2html(wikiContent)

	c.HTML(http.StatusOK, "wiki.gohtml", gin.H{
		"page":             page,
		"title":            page,
		"wikiContent":      template.HTML(output),
		"buttonText":       buttonText,
		"result":           internal.GetMessage(c.Query("m")),
		"lastModifiedTime": lastModifiedTime,
	})
}

func editWiki(c *gin.Context) {
	page := c.Param("page")
	file := getPagesDir() + page
	wikiContent, _ := internal.ReadFile(file)

	c.Header("X-Robots-Tag", "noindex")
	c.HTML(http.StatusOK, "edit.gohtml", gin.H{
		"title":            page,
		"page":             page,
		"wikiContent":      string(wikiContent),
		"turnstileEnabled": getTurnstileEnabled(),
		"turnstileSiteKey": getTurnstileSiteKey(),
	})
}

func saveWiki(c *gin.Context) {
	var requestJson struct {
		Page    string `json:"page"`
		Content string `json:"content"`
		Comment string `json:"comment"`
		Captcha string `json:"captcha"`
	}

	e := c.Bind(&requestJson)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": e.Error()})
		log.Fatal(e)
		return
	}

	if getTurnstileEnabled() {
		captchaResult := captcha.Validate(requestJson.Captcha, getTurnstileSecretKey())
		if !captchaResult {
			return
		}
	}

	page := requestJson.Page
	if page[0:1] != "/" {
		page = "/" + page
	}

	if page == "/" {
		c.JSON(http.StatusBadRequest, gin.H{"result": internal.GetMessage("missing-path")})
		return
	}

	editComment := requestJson.Comment
	if len(editComment) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"result": internal.GetMessage("missing-comment")})
		return
	}

	wikiContent := requestJson.Content
	wikiContentBytes := internal.NormalizeNewlines([]byte(wikiContent))

	filepath := getPagesDir() + page

	if len(wikiContentBytes) == 0 {
		internal.DeleteFile(filepath)
		go internal.CommitFile(getPageDirName()+page, getRepoDir(), editComment, getGitAccessToken())
		c.JSON(http.StatusOK, gin.H{"result": internal.GetMessage("wiki-removed")})
		return
	}

	err := internal.SaveFile(wikiContentBytes, filepath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": internal.GetMessage("save-error")})
		return
	}

	go internal.CommitFile(getPageDirName()+page, getRepoDir(), editComment, getGitAccessToken())

	c.JSON(http.StatusOK, gin.H{"result": internal.GetMessage("wiki-saved")})
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
