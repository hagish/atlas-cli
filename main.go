package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
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
	/*
		config := Config{
			TemplateFile:         "C:\\Users\\hagis\\projects1\\atlas-cli\\files\\template.txt",
			InputDir:             "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources",
			OutputDir:            "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources",
			RelativeFileNameBase: "C:\\Users\\hagis\\projects1\\Colors-of-the-Forest\\Assets\\Generated\\Resources",
			MaxSize:              1024 * 4,
			AtlasName:            "atlas-cfill",
			Pattern:              "_cfill.png",
			InitialColor: ConfigColor{
				R: 0, G: 0, B: 0, A: 0,
			},
			Additional: []ConfigAtlas{
				{
					AtlasName: "atlas-cudf",
					Scale:     2,
					Pattern:   "_cudf.png",
					InitialColor: ConfigColor{
						R: 0, G: 0, B: 0, A: 0,
					},
				},
				{
					AtlasName: "atlas-cn",
					Scale:     1,
					Pattern:   "_cn.png",
					InitialColor: ConfigColor{
						R: 127, G: 127, B: 255, A: 255,
					},
				},
			},
		}
	*/

	var configFile string
	flag.StringVar(&configFile, "config", "atlas.json", "Json config file that specifies the atlas params and input files. All paths within the json are relative to the json files directory.")
	flag.Parse()

	configDir := filepath.Dir(configFile)
	configJson, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	var config Config
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	os.Chdir(configDir)

	log.Printf("Using config: %v", configFile)
	log.Printf("Working directory: %v", configDir)
	log.Printf("Using template file: %v", config.TemplateFile)

	inputDir := filepath.ToSlash(config.InputDir)
	outputDir := filepath.ToSlash(config.OutputDir)
	relativeFileNameBase := strings.TrimRight(filepath.ToSlash(config.RelativeFileNameBase), "/") + "/"

	files_main := listPngFiles(inputDir, config.Pattern)

	maxSize := config.MaxSize

	res, err := generateAtlas(files_main, outputDir, config.AtlasName, maxSize, relativeFileNameBase, config.TemplateFile, config.TemplateExt)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	for _, addConfig := range config.Additional {
		param := new(AdditionalAtlasParams)
		param.name = addConfig.AtlasName
		param.scale = addConfig.Scale
		param.oldPattern = config.Pattern
		param.newPattern = addConfig.Pattern
		param.outputDir = outputDir
		param.initialColor = color.RGBA{
			R: addConfig.InitialColor.R,
			G: addConfig.InitialColor.G,
			B: addConfig.InitialColor.B,
			A: addConfig.InitialColor.A,
		}
		param.RelativeFileNameBase = relativeFileNameBase
		param.TemplateFile = config.TemplateFile
		param.TemplateExt = config.TemplateExt
		generateAdditionalAtlas(res, param)
	}
}

type AdditionalAtlasParams struct {
	name                 string
	scale                int
	oldPattern           string
	newPattern           string
	outputDir            string
	initialColor         color.RGBA
	RelativeFileNameBase string
	TemplateFile         string
	TemplateExt          string
}

func generateAdditionalAtlas(res *GenerateResult, param *AdditionalAtlasParams) {
	for i, atlas := range res.Atlases {
		additionalAtlas := &Atlas{
			Name:         fmt.Sprintf("%s-%d", param.name, (i + 1)),
			Width:        atlas.Width * param.scale,
			Height:       atlas.Height * param.scale,
			MaxWidth:     atlas.MaxWidth * param.scale,
			MaxHeight:    atlas.MaxHeight * param.scale,
			TemplateFile: param.TemplateFile,
			TemplateExt:  param.TemplateExt,
			Padding:      0,
			Gutter:       0,
			Files:        make([]*File, 0),
		}

		// remap Additional files
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
			os.Exit(1)
			return
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
		os.Exit(1)
		return nil
	}

	return files
}

func generateAtlas(files []string, outputDir string, name string, maxSize int, relativeFileNameBase string, templateFile string, templateExt string) (*GenerateResult, error) {
	params := GenerateParams{
		Name:                 name, // The base name of the outputted files
		TemplateFile:         templateFile,
		TemplateExt:          templateExt,
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
