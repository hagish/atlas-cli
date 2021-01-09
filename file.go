package main

// Represents a File to be outputted
type File struct {
	Atlas            *Atlas
	FileName         string
	FileNameRelative string
	X                int
	Y                int
	Width            int
	Height           int
}
