package main

import (
	"path/filepath"
	"text/template"
)

// Get the template for the given descriptor format
// Returns an error if the template can not be parsed
func GetTemplateForFormat(templateFile string) (*template.Template, error) {
	t := template.New(filepath.Base(templateFile))
	return t.ParseFiles(templateFile)
}
