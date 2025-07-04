package internal

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GenerateStaticSite(outputDir, pagesDir, siteBaseURL string) error {
	// Create output directory
	if err := CreateDir(outputDir); err != nil {
		return logErrorf("failed to create output directory: %w", err)
	}

	// Copy static assets (CSS, JS, etc.)
	if err := copyStaticAssets("web/pub", outputDir+"/"); err != nil {
		return logErrorf("failed to copy static assets: %w", err)
	}

	// Get all wiki pages
	pages := ListFiles(pagesDir)

	// Generate index page
	if err := generateStaticIndexPage(pages, outputDir+"/index.html"); err != nil {
		return logErrorf("failed to generate index: %w", err)
	}

	// Generate each wiki page
	wikiDir := outputDir + "/wiki/"
	if err := CreateDir(wikiDir); err != nil {
		return logErrorf("failed to create wiki directory: %w", err)
	}

	for _, page := range pages {
		if err := generateStaticWikiPage(page, wikiDir, pagesDir); err != nil {
			return logErrorf("failed to generate page %s: %w", page, err)
		}
	}

	// Generate sitemap.xml
	if err := GenerateSitemap(pagesDir, pages, outputDir, siteBaseURL); err != nil {
		return logErrorf("failed to generate sitemap: %w", err)
	}

	log.Printf("Static site generated at: %s", outputDir)
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

func generateStaticIndexPage(pages []string, outputPath string) error {
	// Load all templates like Gin does
	templatesGlob := "web/templates/*/*.gohtml"
	tmpl, err := template.ParseGlob(templatesGlob)
	if err != nil {
		return logErrorf("failed to parse templates: %w", err)
	}

	// Remove .md suffix from page names
	cleanedPages := make([]string, len(pages))
	for i, page := range pages {
		cleanedPages[i] = strings.TrimSuffix(page, ".md")
	}

	// Create data structure matching what homepage uses
	data := gin.H{
		"title":  "Wiki",
		"pages":  cleanedPages,
		"result": GetMessage(""),
	}

	// Render the template to a buffer
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "list.gohtml", data); err != nil {
		return logErrorf("failed to execute template: %w", err)
	}

	// Write the content directly as it already includes full HTML structure
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return logErrorf("failed to write index file: %w", err)
	}

	log.Printf("Index file written to: %s", outputPath)
	return nil
}

func generateStaticWikiPage(page, outputDir, pagesDir string) error {
	// Load all templates like Gin does
	templatesGlob := "web/templates/*/*.gohtml"
	tmpl, err := template.ParseGlob(templatesGlob)
	if err != nil {
		return logErrorf("failed to parse templates: %w", err)
	}

	// Read page content
	file := pagesDir + page
	content, err := ReadFile(file)
	if err != nil {
		return logErrorf("failed to read page content: %w", err)
	}

	// Convert markdown to HTML
	htmlContent := Md2html(content)

	// Get last modified time
	fileStat, _ := os.Stat(file)
	lastModifiedTime := fileStat.ModTime().Format(time.UnixDate)

	// Create data structure matching what viewWiki uses
	pageName := strings.TrimSuffix(page, ".md")
	data := gin.H{
		"page":             pageName,
		"title":            pageName,
		"wikiContent":      template.HTML(htmlContent),
		"buttonText":       "Edit",
		"result":           GetMessage(""),
		"lastModifiedTime": lastModifiedTime,
	}

	// Create directory structure if needed
	pageOutputPath := outputDir + pageName + ".html"
	dir := filepath.Dir(pageOutputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Render the template to a buffer
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "wiki.gohtml", data); err != nil {
		return logErrorf("failed to execute template: %w", err)
	}

	// Write the content directly
	if err := os.WriteFile(pageOutputPath, buf.Bytes(), 0644); err != nil {
		return logErrorf("failed to write wiki page file: %w", err)
	}

	log.Printf("Wiki page written to: %s", pageOutputPath)
	return nil
}

func logErrorf(format string, args ...interface{}) error {
	log.Printf(format, args...)
	return fmt.Errorf(format, args...)
}
