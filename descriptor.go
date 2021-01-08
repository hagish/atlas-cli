package main

import (
	"fmt"
	"text/template"
)

// Available descriptor templates
const (
	DESC_INVALID DescriptorFormat = ""
	DESC_KIWI    DescriptorFormat = "kiwi"
)

// Represents a Descriptor Format
type DescriptorFormat string

// Get the template for the given descriptor format
// Returns an error if the template can not be parsed
func GetTemplateForFormat(format DescriptorFormat) (*template.Template, error) {
	t := template.New(fmt.Sprintf("%s.template", format))
	return t.ParseFiles(fmt.Sprintf("templates/%s.template", format))
}

// Gets the file extension for the given descriptor format
// Extension is returned without the separator dot for eg:
// "xml", "json", "yaml"
func GetFileExtForFormat(format DescriptorFormat) string {
	switch format {
	case DESC_KIWI:
		return "json"
	default:
		return ""
	}
}
