package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

func getDataDir() string {
	return "data/pages"
}

func homepage(c *gin.Context) {
	var pages []string
	dataDir := getDataDir()

	err := filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		// Ignore directories and dot files
		if d.IsDir() || d.Name()[0:1] == "." {
			return nil
		}

		pages = append(pages, path[len(dataDir):])
		return nil
	})

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"pages": pages,
	})
}

func viewWiki(c *gin.Context) {
	page := c.Param("page")
	file := getDataDir() + page
	wikiContent, err := os.ReadFile(file)

	if err != nil {
		fmt.Printf(err.Error())
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
	file := getDataDir() + page
	wikiContent, err := os.ReadFile(file)

	if err != nil {
		fmt.Printf(err.Error())
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"title":       page,
		"page":        page,
		"wikiContent": string(wikiContent),
	})
}

func saveWiki(c *gin.Context) {
	page := c.PostForm("page")
	wikiContent := c.PostForm("content")
	filename := getDataDir() + page

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = file.WriteString(wikiContent)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = file.Sync()
	if err != nil {
		return
	}

	c.Redirect(http.StatusSeeOther, "wiki/"+page)
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", homepage)
	router.GET("/wiki/*page", viewWiki)
	router.GET("/edit/*page", editWiki)
	router.POST("/submitWiki", saveWiki)

	err := router.Run()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}
