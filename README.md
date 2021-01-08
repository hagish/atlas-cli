Go Atlas
========

Texture packer written in Go.

### Features

* Set maximum width/height of atlases for platform constraints
* Generate as many atlases as you need with a single command
* Add gutter to the images to prevent join lines between sprites
* Generate descriptor files in a range of formats (Currently only Kiwi.js supported)
* Specify assets that must be grouped together to ensure maximum runtime performance (TODO)

### Example Usage

Basic example
```
inFiles := []string{
	"./assets/sprite1.png",
	"./assets/sprite2.png"
}
outputDir := "./assets/spritesheets"
res, err := atlas.Generate(inFiles, outputDir, nil)
```

You can pass params to the Generate function
```
params := atlas.GenerateParams {
	Name   	   : "atlas" // The base name of the outputted files
	Descriptor : atlas.DESC_KIWI // The format of the data file for the atlases
	Packer     : atlas.PackGrowing // The algorithm to use when packing
	Sorter	   : atlas.SortMaxSide // The order to sort files by
	MaxWidth   : 2048 // Maximum width/height of the atlas images
	MaxHeight  : 2048 
	MaxAtlases : 0 // Indicates no maximum
	Padding    : 0 // The amount of blank space to add around each image
	Gutter     : 0 // The amount to bleed the outer pixels of each image
}
res, err := atlas.Generate(inFiles, outputDir, &params)
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