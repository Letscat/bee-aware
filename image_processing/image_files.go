package image_processing

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
)

func OpenImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func SaveImage(greyPixels [][]uint8, path string) {
	rect := image.Rect(0, 0, len(greyPixels), len(greyPixels[0]))
	nImg := image.NewGray(rect)
	for x := 0; x < len(greyPixels); x++ {
		for y := 0; y < len(greyPixels[0]); y++ {
			q := greyPixels[x]
			if q == nil {
				continue
			}
			nImg.Set(x, y, color.Gray{Y: greyPixels[x][y]})

		}
	}
	f, createErr := os.Create(path)
	if createErr != nil {
		fmt.Println("Creating file:", createErr)
	}
	defer f.Close()
	err := png.Encode(f, nImg)
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}
}
