package main

import (
	"fmt"
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

func main() {
	root := "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources\\Biomes\\Campground"
	outputDir := "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources\\Atlas\\Biomes_Campground"

	files_cfill := listPngFiles(root, "_cfill.png")

	maxSize := 2048

	res, err := generateAtlas(files_cfill, outputDir, "atlas-cfill", maxSize)
	if err != nil {
		log.Fatal(err)
	}

	generateAdditionalAtlas(res, "atlas-cudf", 2,"_cfill.png", "_cudf.png", outputDir)
	generateAdditionalAtlas(res, "atlas-cn", 1,"_cfill.png", "_cn.png", outputDir)
}

func generateAdditionalAtlas(res *GenerateResult, name string, scale int, oldPattern string, newPattern string, outputDir string) {
	for i, atlas := range res.Atlases {
		additionalAtlas := &Atlas{
			Name:       fmt.Sprintf("%s-%d", name, (i + 1)),
			Width:      atlas.Width * scale,
			Height:     atlas.Height * scale,
			MaxWidth:   atlas.MaxWidth * scale,
			MaxHeight:  atlas.MaxHeight * scale,
			Descriptor: DESC_KIWI,
			Padding:    0,
			Gutter:     0,
			Files:      make([]*File, 0),
		}

		// remap additional files
		for _, file := range atlas.Files {
			udfFile := strings.Replace(file.FileName, oldPattern, newPattern, -1)
			if fileExists(udfFile) {
				additionalAtlas.Files = append(additionalAtlas.Files, &File{
					X:        file.X * scale,
					Y:        file.Y * scale,
					FileName: udfFile,
					Width:    file.Width * scale,
					Height:   file.Height * scale,
					Atlas:    additionalAtlas,
				})
			}
		}

		fmt.Printf("Writing atlas named %s to %s\n", additionalAtlas.Name, outputDir)
		err := additionalAtlas.Write(outputDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func listPngFiles(root string, pattern string) []string {
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

func generateAtlas(files []string, outputDir string, name string, maxSize int) (*GenerateResult, error) {
	params := GenerateParams{
		Name:       name,        // The base name of the outputted files
		Descriptor: DESC_KIWI,   // The format of the data file for the atlases
		Packer:     PackGrowing, // The algorithm to use when packing
		Sorter:     SortMaxSide, // The order to sort files by
		MaxWidth:   maxSize,     // Maximum width/height of the atlas images
		MaxHeight:  maxSize,
		MaxAtlases: 0, // Indicates no maximum
		Padding:    0, // The amount of blank space to add around each image
		Gutter:     0, // The amount to bleed the outer pixels of each image
	}

	res, err := Generate(files, outputDir, &params)
	return res, err
}
