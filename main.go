package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"html/template"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/wiki/*page", func(c *gin.Context) {
		page := c.Param("page")
		file := "./data/pages" + page
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
