package main

func getRepoDir() string {
	return "data/repo/"
}

func getPageDirName() string {
	return "pages"
}

func getPagesDir() string {
	return getRepoDir() + getPageDirName()
}
