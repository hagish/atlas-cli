package main

import "testing"

var (
	file_25x100  = &File{FileName: "25x100", Width: 25, Height: 100}
	file_50x50   = &File{FileName: "50x50", Width: 50, Height: 50}
	file_50x300  = &File{FileName: "50x300", Width: 50, Height: 300}
	file_100x100 = &File{FileName: "100x100", Width: 100, Height: 100}
	file_200x200 = &File{FileName: "200x200", Width: 200, Height: 200}
	file_500x50  = &File{FileName: "500x50", Width: 500, Height: 50}
)

var files = []*File{
	file_25x100,
	file_50x50,
	file_50x300,
	file_100x100,
	file_200x200,
	file_500x50,
}

func TestSortWidth(t *testing.T) {
	run(t, SortWidth, []*File{
		file_500x50,
		file_200x200,
		file_100x100,
		file_50x300,
		file_50x50,
		file_25x100,
	})
}

func TestSortHeight(t *testing.T) {
	run(t, SortHeight, []*File{
		file_50x300,
		file_200x200,
		file_100x100,
		file_25x100,
		file_500x50,
		file_50x50,
	})
}

func TestSortAreaWidth(t *testing.T) {
	run(t, SortAreaWidth, []*File{
		file_200x200, // 40,000
		file_500x50,  // 25,000
		file_50x300,  // 15,000
		file_100x100, // 10,000
		file_50x50,   // 2,500
		file_25x100,  // 2,500
	})
}

func TestSortAreaHeight(t *testing.T) {
	run(t, SortAreaHeight, []*File{
		file_200x200, // 40,000
		file_500x50,  // 25,000
		file_50x300,  // 15,000
		file_100x100, // 10,000
		file_25x100,  // 2,500
		file_50x50,   // 2,500
	})
}

func TestSortMaxSide(t *testing.T) {
	run(t, SortMaxSide, []*File{
		file_500x50,
		file_50x300,
		file_200x200,
		file_100x100,
		file_25x100,
		file_50x50,
	})
}

func run(t *testing.T, sorter Sorter, expect []*File) {
	res := sorter(files)
	for i, file := range expect {
		if res[i] != file {
			t.Errorf("File not sorted correctly: want %s, got %s", file.FileName, res[i].FileName)
		}
	}
}
