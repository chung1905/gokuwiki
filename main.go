package main

import (
	"chungn/gokuwiki/internal"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
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
		"pages": internal.ListFiles(getPagesDir()),
	})
}

func viewWiki(c *gin.Context) {
	page := c.Param("page")
	file := getPagesDir() + page
	wikiContent, err := internal.ReadFile(file)

	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}

	output := internal.Md2html(wikiContent)

	c.HTML(http.StatusOK, "wiki.html", gin.H{
		"title":       page,
		"wikiContent": template.HTML(output),
	})
}

func editWiki(c *gin.Context) {
	page := c.Param("page")
	file := getPagesDir() + page
	wikiContent, _ := internal.ReadFile(file)

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"title":       page,
		"page":        page,
		"wikiContent": string(wikiContent),
	})
}

func saveWiki(c *gin.Context) {
	page := c.PostForm("page")
	wikiContent := c.PostForm("content")
	wikiContentBytes := internal.NormalizeNewlines([]byte(wikiContent))

	filepath := getPagesDir() + page

	if len(wikiContentBytes) == 0 {
		internal.DeleteFile(filepath)
		go internal.CommitFile(page, getRepoDir())
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	internal.SaveFile(wikiContentBytes, filepath)
	go internal.CommitFile(page, getRepoDir())

	c.Redirect(http.StatusSeeOther, "wiki/"+page)
}

func main() {
	internal.CreateDir(getPagesDir())
	internal.PrepareGitRepo(getRepoDir())
	router := getRouter()
	err := router.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
