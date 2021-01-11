package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func listPngFiles(root string, pattern string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		path = filepath.ToSlash(path)
		if strings.HasSuffix(path, ".png") && strings.Contains(path, pattern) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return nil
	}

	return files
}
