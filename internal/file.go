package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func CreateDir(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func DeleteFile(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func ReadFile(filepath string) ([]byte, error) {
	content, err := os.ReadFile(filepath)
	return content, err
}

func ListFiles(dir string) []string {
	var pages []string
	dataDirLen := len(dir)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		if !isAllow(path[dataDirLen:], d) {
			return nil
		}

		pages = append(pages, path[len(dir)+1:])
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return pages
}

func SaveFile(content []byte, filepath string) {
	CreateDir(path.Dir(filepath))
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = file.Write(content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = file.Sync()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func isAllow(path string, d fs.DirEntry) bool {
	// Ignore directories and dot files
	if d.IsDir() || path[0:1] == "." || d.Name()[0:1] == "." {
		return false
	}

	return true
}
