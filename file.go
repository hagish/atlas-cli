package main

import (
	"path"
	"strings"
)

// Represents a File to be outputted
type File struct {
	Atlas            *Atlas
	FileName         string
	FileNameRelative string
	FileNameOnly     string
	Name             string
	X                int
	Y                int
	Width            int
	Height           int
}

func (f *File) Complete () {
	f.FileNameOnly = path.Base(f.FileName)
	f.Name = strings.TrimSuffix(f.FileNameOnly, ".png")
}
