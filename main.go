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

	generateParams := GenerateParams{
		Name:                 config.AtlasName,
		RelativeFileNameBase: relativeFileNameBase,
		TemplateFile:         config.TemplateFile,
		TemplateExt:          config.TemplateExt,
		Padding:              config.Padding,
		Gutter:               config.Gutter,
		PowerOfTwo:           true,
		MaxWidth:             config.MaxSize,
		MaxHeight:            config.MaxSize,
		MaxAtlases:           0,
	}

	res, err := generateMainAtlas(files_main, outputDir, &generateParams)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	for _, addConfig := range config.Additional {
		param := new(AdditionalAtlasParams)
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
		generateParams.Name = addConfig.AtlasName
		param.generateParams = &generateParams
		generateAdditionalAtlas(res, param)
	}
}

type AdditionalAtlasParams struct {
	scale          int
	oldPattern     string
	newPattern     string
	outputDir      string
	initialColor   color.RGBA
	generateParams *GenerateParams
}

func generateMainAtlas(files []string, outputDir string, params *GenerateParams) (*GenerateResult, error) {
	res, err := generate(files, outputDir, params)
	return res, err
}

func generateAdditionalAtlas(res *GenerateResult, params *AdditionalAtlasParams) {
	for i, atlas := range res.Atlases {
		additionalAtlas := &Atlas{
			Name:         fmt.Sprintf("%s-%d", params.generateParams.Name, (i + 1)),
			Width:        atlas.Width * params.scale,
			Height:       atlas.Height * params.scale,
			MaxWidth:     atlas.MaxWidth * params.scale,
			MaxHeight:    atlas.MaxHeight * params.scale,
			TemplateFile: params.generateParams.TemplateFile,
			TemplateExt:  params.generateParams.TemplateExt,
			Padding:      params.generateParams.Padding,
			Gutter:       params.generateParams.Gutter,
			Files:        make([]*File, 0),
		}

		// remap Additional files
		for _, file := range atlas.Files {
			udfFile := strings.Replace(file.FileName, params.oldPattern, params.newPattern, -1)
			if fileExists(udfFile) {
				f := &File{
					X:                file.X * params.scale,
					Y:                file.Y * params.scale,
					FileName:         udfFile,
					FileNameRelative: strings.Replace(udfFile, params.generateParams.RelativeFileNameBase, "", 1),
					Width:            file.Width * params.scale,
					Height:           file.Height * params.scale,
					Atlas:            additionalAtlas,
				}
				f.Complete()
				additionalAtlas.Files = append(additionalAtlas.Files, f)
			}
		}

		fmt.Printf("Writing atlas named %s to %s\n", additionalAtlas.Name, params.outputDir)
		err := additionalAtlas.Write(params.outputDir, params.initialColor)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
			return
		}
	}
}
