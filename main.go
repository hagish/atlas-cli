package main

import (
	"fmt"
	"image/color"
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
	inputDir := "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources"
	outputDir := "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources\\Atlas\\Biomes_Campground"
	relativeFileNameBase := "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources"

	inputDir = filepath.ToSlash(inputDir)
	outputDir = filepath.ToSlash(outputDir)
	relativeFileNameBase = strings.TrimRight(filepath.ToSlash(relativeFileNameBase), "/") + "/"

	files_cfill := listPngFiles(inputDir, "_cfill.png")

	maxSize := 2048 * 2

	res, err := generateAtlas(files_cfill, outputDir, "atlas-cfill", maxSize, relativeFileNameBase)
	if err != nil {
		log.Fatal(err)
	}

	param := new(AdditionalAtlasParams)
	param.name = "atlas-cudf"
	param.scale = 2
	param.oldPattern = "_cfill.png"
	param.newPattern = "_cudf.png"
	param.outputDir = outputDir
	param.initialColor = color.RGBA{R: 0, G: 0, B: 0, A: 0}
	param.RelativeFileNameBase = relativeFileNameBase

	generateAdditionalAtlas(res, param)

	param.name = "atlas-cn"
	param.newPattern = "_cn.png"
	param.scale = 1
	param.initialColor = color.RGBA{R: 127, G: 127, B: 255, A: 255}

	generateAdditionalAtlas(res, param)
}

type AdditionalAtlasParams struct {
	name                 string
	scale                int
	oldPattern           string
	newPattern           string
	outputDir            string
	initialColor         color.RGBA
	RelativeFileNameBase string
}

func generateAdditionalAtlas(res *GenerateResult, param *AdditionalAtlasParams) {
	for i, atlas := range res.Atlases {
		additionalAtlas := &Atlas{
			Name:       fmt.Sprintf("%s-%d", param.name, (i + 1)),
			Width:      atlas.Width * param.scale,
			Height:     atlas.Height * param.scale,
			MaxWidth:   atlas.MaxWidth * param.scale,
			MaxHeight:  atlas.MaxHeight * param.scale,
			Descriptor: DESC_KIWI,
			Padding:    0,
			Gutter:     0,
			Files:      make([]*File, 0),
		}

		// remap additional files
		for _, file := range atlas.Files {
			udfFile := strings.Replace(file.FileName, param.oldPattern, param.newPattern, -1)
			if fileExists(udfFile) {
				f := &File{
					X:                file.X * param.scale,
					Y:                file.Y * param.scale,
					FileName:         udfFile,
					FileNameRelative: strings.Replace(udfFile, param.RelativeFileNameBase, "", 1),
					Width:            file.Width * param.scale,
					Height:           file.Height * param.scale,
					Atlas:            additionalAtlas,
				}
				f.Complete()
				additionalAtlas.Files = append(additionalAtlas.Files, f)
			}
		}

		fmt.Printf("Writing atlas named %s to %s\n", additionalAtlas.Name, param.outputDir)
		err := additionalAtlas.Write(param.outputDir, param.initialColor)
		if err != nil {
			log.Fatal(err)
		}
	}
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
	}

	return files
}

func generateAtlas(files []string, outputDir string, name string, maxSize int, relativeFileNameBase string) (*GenerateResult, error) {
	params := GenerateParams{
		Name:                 name,        // The base name of the outputted files
		Descriptor:           DESC_KIWI,   // The format of the data file for the atlases
		Packer:               PackGrowing, // The algorithm to use when packing
		Sorter:               SortMaxSide, // The order to sort files by
		MaxWidth:             maxSize,     // Maximum width/height of the atlas images
		MaxHeight:            maxSize,
		MaxAtlases:           0, // Indicates no maximum
		Padding:              0, // The amount of blank space to add around each image
		Gutter:               0, // The amount to bleed the outer pixels of each image
		PowerOfTwo:           true,
		RelativeFileNameBase: relativeFileNameBase,
	}

	res, err := generate(files, outputDir, &params)
	return res, err
}
