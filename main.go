package main

import (
	"bytes"
	"chungn/gokuwiki/internal"
	"chungn/gokuwiki/internal/captcha"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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

	pageFilePath := getPagesDir() + page
	originalPageFilePath := getPagesDir() + originalPage

	// Handle page move
	if originalPage != page {
		internal.DeleteFile(originalPageFilePath)
	}

	if len(wikiContentBytes) == 0 {
		internal.DeleteFile(pageFilePath)
		go internal.CommitFile(getPageDirName()+page, getRepoDir(), editComment, getGitAccessToken())
		c.JSON(http.StatusOK, gin.H{"result": internal.GetMessage("wiki-removed")})
		return
	}

	err := internal.SaveFile(wikiContentBytes, pageFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": internal.GetMessage("save-error")})
		return
	}

	// todo: only generate the changed page
	if err := generateStaticSite(getOutputDir()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": internal.GetMessage("save-error")})
		return
	}

	go internal.CommitFile(getPageDirName()+page, getRepoDir(), editComment, getGitAccessToken())

	c.JSON(http.StatusOK, gin.H{"result": internal.GetMessage("wiki-saved")})
}

func generateStaticSite(outputDir string) error {
	// Create output directory
	if err := internal.CreateDir(outputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy static assets (CSS, JS, etc.)
	if err := copyStaticAssets("web/pub", outputDir+"/"); err != nil {
		return fmt.Errorf("failed to copy static assets: %w", err)
	}

	// Get all wiki pages
	pages := internal.ListFiles(getPagesDir())

	// Generate index page
	if err := generateStaticIndexPage(pages, outputDir+"/index.html"); err != nil {
		return fmt.Errorf("failed to generate index: %w", err)
	}

	// Generate each wiki page
	wikiDir := outputDir + "/wiki/"
	if err := internal.CreateDir(wikiDir); err != nil {
		return fmt.Errorf("failed to create wiki directory: %w", err)
	}

	for _, page := range pages {
		if err := generateStaticWikiPage(page, wikiDir); err != nil {
			return fmt.Errorf("failed to generate page %s: %w", page, err)
		}
	}

	fmt.Printf("Static site generated at: %s\n", outputDir)
	return nil
}

func generateStaticIndexPage(pages []string, outputPath string) error {
	// Load all templates like Gin does
	templatesGlob := "web/templates/*/*.gohtml"
	tmpl, err := template.ParseGlob(templatesGlob)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Create data structure matching what homepage uses
	data := gin.H{
		"title":  "Wiki",
		"pages":  pages,
		"result": internal.GetMessage(""),
	}

	// Render the template to a buffer
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "list.gohtml", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write the content directly as it already includes full HTML structure
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	fmt.Printf("Index file written to: %s\n", outputPath)
	return nil
}

func generateStaticWikiPage(page string, outputDir string) error {
	// Load all templates like Gin does
	templatesGlob := "web/templates/*/*.gohtml"
	tmpl, err := template.ParseGlob(templatesGlob)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Read page content
	file := getPagesDir() + page
	content, err := internal.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read page content: %w", err)
	}

	// Convert markdown to HTML
	htmlContent := internal.Md2html(content)

	// Get last modified time
	fileStat, _ := os.Stat(file)
	lastModifiedTime := fileStat.ModTime().Format(time.UnixDate)

	// Create data structure matching what viewWiki uses
	data := gin.H{
		"page":             page,
		"title":            page,
		"wikiContent":      template.HTML(htmlContent),
		"buttonText":       "Edit",
		"result":           internal.GetMessage(""),
		"lastModifiedTime": lastModifiedTime,
	}

	// Create directory structure if needed
	pageOutputPath := outputDir + page + ".html"
	dir := filepath.Dir(pageOutputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Render the template to a buffer
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "wiki.gohtml", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write the content directly
	if err := os.WriteFile(pageOutputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write wiki page file: %w", err)
	}

	fmt.Printf("Wiki page written to: %s\n", pageOutputPath)
	return nil
}

func copyStaticAssets(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, 0644)
	})
}

func main() {
	_ = internal.CreateDir(getPagesDir())
	internal.PrepareGitRepo(getRepoDir(), getRepoURL())

	if err := generateStaticSite(getOutputDir()); err != nil {
		log.Printf("Error generating static site: %v", err)
	}

	router := getRouter()
	err := router.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
