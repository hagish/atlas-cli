package main

type Config struct {
	TemplateFile         string
	TemplateExt          string
	InputDir             string
	OutputDir            string
	RelativeFileNameBase string
	MaxSize              int

	Pattern      string
	AtlasName    string
	InitialColor ConfigColor

	Additional []ConfigAtlas
}

type ConfigColor struct {
	R, G, B, A uint8
}

type ConfigAtlas struct {
	Pattern      string
	Scale        int
	AtlasName    string
	InitialColor ConfigColor
}
