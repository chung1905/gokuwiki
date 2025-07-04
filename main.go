package main

import (
	"chungn/gokuwiki/internal"
	"chungn/gokuwiki/internal/captcha"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func getRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*/*.gohtml")
	router.Use(static.Serve("/", static.LocalFile(getOutputDir(), false)))
	router.GET("/edit/*page", editWiki)
	router.POST("/submitWiki", saveWiki)

	return router
}

func editWiki(c *gin.Context) {
	page := c.Param("page")
	pageFile := page + ".md"
	file := getPagesDir() + pageFile
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
		OriginalPage string `json:"original-page"`
		Page         string `json:"page"`
		Content      string `json:"content"`
		Comment      string `json:"comment"`
		Captcha      string `json:"captcha"`
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
	originalPage := requestJson.OriginalPage
	if originalPage[0:1] != "/" {
		originalPage = "/" + originalPage
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

	pageFilePath := getPagesDir() + page + ".md"
	originalPageFilePath := getPagesDir() + originalPage + ".md"

	// Handle page move
	if originalPage != page {
		internal.DeleteFile(originalPageFilePath)
	}

	if len(wikiContentBytes) == 0 {
		internal.DeleteFile(pageFilePath)
		go internal.CommitFile(pageFilePath, getRepoDir(), editComment, getGitAccessToken())
		c.JSON(http.StatusOK, gin.H{"result": internal.GetMessage("wiki-removed")})
		return
	}

	err := internal.SaveFile(wikiContentBytes, pageFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": internal.GetMessage("save-error")})
		return
	}

	// todo: only generate the changed page
	if err := internal.GenerateStaticSite(getOutputDir(), getPagesDir(), getSiteBaseURL()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": internal.GetMessage("save-error")})
		return
	}

	if originalPage != page {
		go internal.CommitFiles([]string{
			pageFilePath,
			originalPageFilePath,
		}, getRepoDir(), editComment, getGitAccessToken())
	} else {
		go internal.CommitFile(pageFilePath, getRepoDir(), editComment, getGitAccessToken())
	}

	c.JSON(http.StatusOK, gin.H{"result": internal.GetMessage("wiki-saved")})
}

func main() {
	_ = internal.CreateDir(getPagesDir())
	internal.PrepareGitRepo(getRepoDir(), getRepoURL(), getGitAccessToken())

	err := internal.GenerateStaticSite(getOutputDir(), getPagesDir(), getSiteBaseURL())
	if err != nil {
		log.Printf("Error generating static site: %v", err)
	}

	router := getRouter()
	if err := router.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
