Go Atlas
========

Texture packer command line tool written in Go.
This is mostly a command line tool around https://github.com/ikkeps/atlas with the option to add additional atlases with the same layout.
Not all the original features are exposed in the config. Beware! The code may be quick and dirty.

### Features

* Set maximum width/height of atlases for platform constraints
* Generate as many atlases as you need with a single command
* Generate descriptor files in a range of formats
* Generate a set of additional texture atlases with the same layout for additional images (e.g. normal maps)
* Atlases are power-of-two
* Only supports `.png` files (and probably not with `.PNG`)

### Example Usage

Config file example
```
{
  "templateFile": "template.txt",
  "templateExt": "json",
  "inputDir": "Sprites",
  "outputDir": "Atlas",
  "relativeFileNameBase": "Sprites",
  "maxSize": 4096,
  "atlasName": "atlas-main",
  "pattern": "_main.png",
  "initialColor": {"r": 0, "g": 0, "b": 0, "a": 0},
  "additional": [
    {
      "atlasName": "atlas-normal",
      "scale": 1,
      "pattern": "_normal.png",
      "initialColor": { "r": 127, "g": 127, "b": 255, "a": 255}
    }
  ]
}
```
Run the command with `--config yourconfigfile.json`.
It will search for `.png` files in the `inputDir` that have `pattern` in the filename.
The additional atlases take the those input filenames and replace the main `pattern` with `pattern` from the additional.  

All paths within the config file are relative to the directory of the config file.

Template file example
```
{
	"name": "{{.Name}}",
	"cells": [
		{{with .Files}}{{range $index, $el := .}}{{if $index}},{{end}}{
	        "x": {{$el.X}},
	        "y": {{$el.Y}},
	        "w": {{$el.Width}},
	        "h": {{$el.Height}},
	        "fileName": "{{$el.FileName}}",
	        "relativeFileName": "{{$el.FileNameRelative}}",
	        "fileNameOnly": "{{$el.FileNameOnly}}",
	        "name": "{{$el.Name}}",
	    }{{end}}{{end}}
    ]
}
```

Example file structure
```
atlas.json
template.txt
Sprites/s1_main.png
Sprites/s1_normal.png
Sprites/s2_main.png
Sprites/s2_normal.png
Atlas/
```

### License

> This is free and unencumbered software released into the public domain.

> Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

> In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

> For more information, please refer to <http://unlicense.org/>