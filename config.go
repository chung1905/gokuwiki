package main

import (
	"os"
	"strings"
)

func getOutputDir() string {
	dir := os.Getenv("OUTPUT_DIR")
	if dir == "" {
		return "output"
	}
	return dir
}

func getSiteBaseURL() string {
	return strings.TrimRight(os.Getenv("GOKUWIKI_SITE_BASE_URL"), "/")
}

func getRepoDir() string {
	return "data/repo/"
}

func getPageDirName() string {
	return "pages/"
}

func getPagesDir() string {
	return getRepoDir() + getPageDirName()
}

func getRepoURL() string {
	return os.Getenv("GOKUWIKI_REPO_URL")
}

func getGitAccessToken() string {
	return os.Getenv("GOKUWIKI_ACCESS_TOKEN")
}

func getTurnstileEnabled() bool {
	return os.Getenv("GOKUWIKI_TURNSTILE_ENABLED") == "true"
}

func getTurnstileSiteKey() string {
	return os.Getenv("GOKUWIKI_TURNSTILE_SITE_KEY")
}

func getTurnstileSecretKey() string {
	return os.Getenv("GOKUWIKI_TURNSTILE_SECRET_KEY")
}
