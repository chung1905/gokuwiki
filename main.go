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

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	dataDir := "data/pages"

	router.GET("/", func(c *gin.Context) {
		pages := []string{}
		filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			pages = append(pages, path[len(dataDir):])
			return nil
		})

		c.HTML(http.StatusOK, "list.html", gin.H{
			"pages": pages,
		})
	})

	router.GET("/wiki/*page", func(c *gin.Context) {
		page := c.Param("page")
		file := dataDir + page
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
	})

	err := router.Run()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}
