package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
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

	output := markdown.ToHTML(wikiContent, nil, nil)

	c.HTML(http.StatusOK, "wiki.html", gin.H{
		"title":       page,
		"wikiContent": template.HTML(output),
	})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", homepage)
	router.GET("/wiki/*page", viewWiki)

	err := router.Run()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}
