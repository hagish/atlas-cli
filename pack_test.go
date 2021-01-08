package main

import "testing"

func TestPackGrowing(t *testing.T) {
	FILES := []*File{
		&File{Width: 200, Height: 200},
		&File{Width: 100, Height: 100},
		&File{Width: 50, Height: 50},
	}

	type In struct {
		files []*File
		atlas *Atlas
	}

	type Want struct {
		files    []*File
		atlas    *Atlas
		numUnfit int
	}

	cases := []struct {
		in   In
		want Want
	}{
		{
			// This is a basic control test case to ensure
			// normal functionaly. We expect all assets to
			// fit and for the atlas to come out smaller than
			// the maximum size

			in: In{
				files: FILES,
				atlas: &Atlas{
					MaxWidth:  400,
					MaxHeight: 400,
				},
			},
			want: Want{
				files: []*File{
					&File{X: 0, Y: 0},
					&File{X: 200, Y: 0},
					&File{X: 200, Y: 100},
				},
				atlas: &Atlas{
					Width:  300,
					Height: 200,
				},
			},
		},
		{
			// In this case we are setting the atlas size to the
			// size of the largest asset. We expect that the other
			// two assets will be excluded

			in: In{
				files: FILES,
				atlas: &Atlas{
					MaxWidth:  200,
					MaxHeight: 200,
				},
			},
			want: Want{
				atlas: &Atlas{
					Width:  200,
					Height: 200,
				},
				numUnfit: 2,
			},
		},
		{
			// In this case we are expecting the packing
			// to fail as we are giving the atlas a maximum size
			// smaller than that of the assets that will fill it

			in: In{
				files: FILES,
				atlas: &Atlas{
					MaxWidth:  1,
					MaxHeight: 1,
				},
			},
			want: Want{
				numUnfit: 3,
			},
		},
	}

	for _, c := range cases {
		PackGrowing(c.in.atlas, c.in.files)
		numUnfit := 0
		for _, file := range c.in.files {
			if file.Atlas == nil {
				numUnfit += 1
			}
		}
		if numUnfit != c.want.numUnfit {
			t.Errorf("Unexpected number of unfit file(s): want %d, got %d", c.want.numUnfit, numUnfit)
		}
		if c.want.atlas != nil {
			if c.in.atlas.Width != c.want.atlas.Width || c.in.atlas.Height != c.want.atlas.Height {
				t.Errorf("Unexpected atlas size: want %dx%d, got %dx%d",
					c.want.atlas.Width, c.want.atlas.Height, c.in.atlas.Width, c.in.atlas.Height)
			}
		}
		if c.want.files != nil {
			for i, file := range c.in.files {
				expect := c.want.files[i]
				if file.X != expect.X || file.Y != expect.Y {
					t.Errorf("Unexpected file position: want %d,%d, got %d,%d",
						expect.X, expect.Y, file.X, file.Y)
				}
			}
		}
		for _, file := range c.in.files {
			file.Atlas = nil
		}
	}
}
