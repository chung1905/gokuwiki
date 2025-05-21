package internal

import (
	"fmt"
	"os"
	"strings"
)

func GenerateSitemap(pagesDir string, pages []string, outputDir string, siteBaseURL string) error {
	// Create the XML structure for sitemap
	xmlStart := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
        xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	xmlEnd := `
</urlset>`

	var sitemap strings.Builder
	sitemap.WriteString(xmlStart)

	// Add home page
	sitemap.WriteString(`
  <url>
    <loc>` + siteBaseURL + `/</loc>
    <changefreq>weekly</changefreq>
    <priority>1.0</priority>
  </url>`)

	// Add each wiki page
	for _, page := range pages {
		pagePath := strings.TrimSuffix(page, ".md")

		// Get last modified time
		file := pagesDir + page
		fileStat, err := os.Stat(file)
		lastMod := ""
		if err == nil {
			lastMod = fileStat.ModTime().Format("2006-01-02")
		}

		sitemap.WriteString(`
  <url>
    <loc>` + siteBaseURL + `/wiki/` + pagePath + `.html</loc>`)

		if lastMod != "" {
			sitemap.WriteString(`
    <lastmod>` + lastMod + `</lastmod>`)
		}

		sitemap.WriteString(`
    <changefreq>monthly</changefreq>
    <priority>0.8</priority>
  </url>`)
	}

	sitemap.WriteString(xmlEnd)

	// Write the sitemap to file
	if err := os.WriteFile(outputDir+"/sitemap.xml", []byte(sitemap.String()), 0644); err != nil {
		return fmt.Errorf("failed to write sitemap.xml: %w", err)
	}

	fmt.Println("Sitemap generated at:", outputDir+"/sitemap.xml")
	return nil
}
