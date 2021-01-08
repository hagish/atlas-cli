package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	root := "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources\\Biomes\\Campground"
	outputDir := "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources\\Atlas\\Biomes_Campground"

	files_cfill := ListPngFiles(root, "_cfill.png")
	files_cudf := ListPngFiles(root, "_cudf.png")

	var files_cfill_cleaned []string
	var files_cudf_cleaned []string

	for _, s := range files_cfill {
		other := strings.Replace(s, "_cfill.png", "_cudf.png", -1)
		if Contains(files_cudf, other) {
			files_cfill_cleaned = append(files_cfill_cleaned, s)
		}
	}

	for _, s := range files_cudf {
		other := strings.Replace(s, "_cudf.png", "_cfill.png", -1)
		if Contains(files_cfill, other) {
			files_cudf_cleaned = append(files_cudf_cleaned, s)
		}
	}

	GenerateAtlas(files_cfill_cleaned, outputDir, "atlas-cfill",2048)
	GenerateAtlas(files_cudf_cleaned, outputDir, "atlas-cudf",2048 * 2)
}

func ListPngFiles(root string, pattern string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".png") && strings.Contains(path, pattern) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return files
}

func GenerateAtlas(files []string, outputDir string, name string, maxSize int) {
	params := GenerateParams {
		Name   	   : name, // The base name of the outputted files
		Descriptor : DESC_KIWI, // The format of the data file for the atlases
		Packer     : PackGrowing, // The algorithm to use when packing
		Sorter	   : SortMaxSide, // The order to sort files by
		MaxWidth   : maxSize, // Maximum width/height of the atlas images
		MaxHeight  : maxSize,
		MaxAtlases : 0, // Indicates no maximum
		Padding    : 0, // The amount of blank space to add around each image
		Gutter     : 0, // The amount to bleed the outer pixels of each image
	}

	_, err := Generate(files, outputDir, &params)
	if err != nil {
		log.Fatal(err)
	}
}
