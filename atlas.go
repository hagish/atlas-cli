package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path"
)

// Represents a single atlas to be outputted
type Atlas struct {
	Name                string
	Files               []*File
	Width, Height       int
	MaxWidth, MaxHeight int
	Padding, Gutter     int
	Descriptor          DescriptorFormat
}

// Adds a file into the atlas at the given position
// Used by packers to add packed files to the atlas
func (a *Atlas) AddFile(file *File, x, y int) {
	file.Atlas = a
	file.X, file.Y = x, y
	a.Files = append(a.Files, file)
}

// Writes the atlas to the given output directory, this is shorthand
// for calling both WriteImage and WriteDescriptor
func (a *Atlas) Write(outputDir string) error {
	if err := a.WriteImage(outputDir); err == nil {
		return a.WriteDescriptor(outputDir)
	} else {
		return err
	}
}

// Writes the image for this atlas to the given output directory
// Returns an error if any IO operation fails
func (a *Atlas) WriteImage(outputDir string) (err error) {
	// Generate the image data
	im := image.NewRGBA(image.Rect(0, 0, a.Width, a.Height))
	// Set the background colour of the image
	for i, n := 0, len(im.Pix); i < n; i += 4 {
		im.Pix[i] = 0   // Red
		im.Pix[i+1] = 0 // Green
		im.Pix[i+2] = 0 // Blue
		im.Pix[i+3] = 0 // Alpha
	}

	if a.Gutter == 0 {
		err = compositeImage(a.Files, func(file *File, cim image.Image) {
			dp := image.Pt(file.X+a.Padding, file.Y+a.Padding)
			draw.Draw(im, image.Rectangle{dp, dp.Add(cim.Bounds().Size())}, cim, cim.Bounds().Min, draw.Src)
		})
	} else {
		err = compositeImage(a.Files, func(file *File, cim image.Image) {
			// Create a temp image with padding for the gutter
			cimSize := cim.Bounds().Size()
			tempRect := image.Rect(0, 0, cimSize.X+a.Gutter*2, cimSize.Y+a.Gutter*2)
			temp := image.NewRGBA(tempRect)
			dp := image.Pt(a.Gutter, a.Gutter)
			draw.Draw(temp, image.Rectangle{dp, dp.Add(cimSize)}, cim, cim.Bounds().Min, draw.Src)
			// Bleed the image into the gutter space
			bleed(temp, a.Gutter)
			// Now draw the image with the gutter onto the texture atlas
			dp = image.Pt(file.X+a.Padding, file.Y+a.Padding)
			draw.Draw(im, image.Rectangle{dp, dp.Add(tempRect.Size())}, temp, temp.Bounds().Min, draw.Src)
		})
	}

	if err != nil {
		return err
	}

	out, err := os.Create(path.Join(outputDir, fmt.Sprintf("%s.png", a.Name)))
	if err != nil {
		return err
	}

	err = png.Encode(out, im)
	if err != nil {
		return err
	}
	return nil
}

// Runs through all the given files, reading their image and then performs the op function on them
func compositeImage(files []*File, op func(file *File, cim image.Image)) error {
	for _, file := range files {
		// Open the given file
		r, err := os.Open(file.FileName)
		if err != nil {
			return err
		}
		// Decode the image
		cim, _, err := image.Decode(r)
		if err != nil {
			return err
		}
		op(file, cim)
		r.Close()
	}
	return nil
}

// Copy an images outer edge of pixels and "bleeds" them to the edge of the image
func bleed(im draw.Image, amount int) {
	// TODO refactor
	outer := im.Bounds()
	inner := image.Rect(outer.Min.X+amount, outer.Min.Y+amount, outer.Max.X-amount, outer.Max.Y-amount)

	// Top left
	col := im.At(inner.Min.X, inner.Min.Y)
	for x := 0; x < amount; x++ {
		for y := 0; y < amount; y++ {
			im.Set(x, y, col)
		}
	}
	// Top edge
	for x := inner.Min.X; x < inner.Max.X; x++ {
		col := im.At(x, amount)
		for y := 0; y < amount; y++ {
			im.Set(x, y, col)
		}
	}
	// Top right
	col = im.At(inner.Max.X-1, inner.Min.Y)
	for x := 1; x <= amount; x++ {
		for y := 0; y < amount; y++ {
			im.Set(outer.Max.X-x, y, col)
		}
	}
	// Right edge
	for y := amount; y < inner.Max.X; y++ {
		col := im.At(inner.Max.X-1, y)
		for x := 1; x <= amount; x++ {
			im.Set(outer.Max.X-x, y, col)
		}
	}
	// Bottom right
	col = im.At(inner.Max.X-1, inner.Max.Y-1)
	for x := 1; x <= amount; x++ {
		for y := 1; y <= amount; y++ {
			im.Set(outer.Max.X-x, outer.Max.Y-y, col)
		}
	}
	// Bottom edge
	for x := amount; x < inner.Max.X; x++ {
		col := im.At(x, inner.Max.Y-1)
		for y := 1; y <= amount; y++ {
			im.Set(x, outer.Max.Y-y, col)
		}
	}
	// Bottom left
	col = im.At(inner.Min.X, inner.Max.Y-1)
	for x := 0; x < amount; x++ {
		for y := 1; y <= amount; y++ {
			im.Set(x, outer.Max.Y-y, col)
		}
	}
	// Left edge
	for y := amount; y < inner.Max.Y; y++ {
		col := im.At(amount, y)
		for x := 0; x < amount; x++ {
			im.Set(x, y, col)
		}
	}
}

// Writes the descriptor file for this atlas to the given output directory
// Returns an error if any IO operation fails
func (a *Atlas) WriteDescriptor(outputDir string) error {
	t, err := GetTemplateForFormat(a.Descriptor)
	if err != nil {
		return err
	}
	ext := GetFileExtForFormat(a.Descriptor)
	out, err := os.Create(path.Join(outputDir, fmt.Sprintf("%s.%s", a.Name, ext)))
	if err != nil {
		return err
	}
	return t.Execute(out, a)
}
