package image_processing

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"sync"
)

func Base64ToImage(base64String string) (image.Image, error) {

	decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewBuffer(decodedBytes))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ImageToGrayscalePixels(img image.Image) [][]uint8 {
	size := img.Bounds().Size()
	pixels := make([][]uint8, size.Y)
	for j := range pixels {
		pixels[j] = make([]uint8, size.X)
	}

	var wg sync.WaitGroup

	for y := 0; y < size.Y; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			row := make([]uint8, size.X)
			for x := 0; x < size.X; x++ {
				originalColor := img.At(x, y)
				grayColor := color.GrayModel.Convert(originalColor).(color.Gray)
				row[x] = grayColor.Y
			}
			pixels[y] = row
		}(y)
	}
	wg.Wait()
	return pixels
}

func StackBlur(pixels *[][]uint8, kernelSize int) {
	radius := kernelSize / 2

	for pixelX := 0; pixelX < len(*pixels); pixelX++ {
		for pixelY := 0; pixelY < len((*pixels)[pixelX]); pixelY++ {
			sum := calculateWeightedSum(pixels, pixelX, pixelY, radius, kernelSize)
			(*pixels)[pixelX][pixelY] = sum / uint8(kernelSize*kernelSize)
		}
	}
}

func calculateWeightedSum(pixels *[][]uint8, pixelX, pixelY, radius, kernelSize int) uint8 {
	sum := uint8(0)
	for kernelIndex := 0; kernelIndex < kernelSize*kernelSize; kernelIndex++ {
		kernelX := kernelIndex % kernelSize
		kernelY := kernelIndex / kernelSize
		if pixelX+kernelX-radius >= 0 && pixelX+kernelX-radius < len(*pixels) && pixelY+kernelY-radius >= 0 && pixelY+kernelY-radius < len((*pixels)[pixelX]) {
			sum += (*pixels)[pixelX+kernelX-radius][pixelY+kernelY-radius]
		}
	}
	return sum
}

func GlobalThresholding(image [][]uint8, threshold uint8) {
	for i := range image {
		for j := range image[i] {
			if image[i][j] > threshold {
				image[i][j] = 255
			} else {
				image[i][j] = 0
			}
		}
	}
}
func SubtractImages(img1, img2 [][]uint8) [][]uint8 {
	if len(img1) != len(img2) || len(img1[0]) != len(img2[0]) {
		panic("Images must be the same size")
	}

	result := make([][]uint8, len(img1))
	for i := range img1 {
		result[i] = make([]uint8, len(img1[0]))
		for j := range img1[0] {
			result[i][j] = uint8(max(0, int(img1[i][j])-int(img2[i][j])))
		}
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func FindContour(image [][]uint8, n int) bool {

	for i := range image {
		for j := range image[i] {
			if image[i][j] == 255 {
				contour := []int{}
				dfs(image, i, j, &contour)
				if len(contour) > n {
					return true
				}
			}
		}
	}
	return false
}

func dfs(image [][]uint8, x, y int, contour *[]int) {
	if x < 0 || x >= len(image) || y < 0 || y >= len(image[x]) || image[x][y] != 255 {
		return
	}

	image[x][y] = 0

	*contour = append(*contour, x, y)

	dfs(image, x-1, y, contour)
	dfs(image, x+1, y, contour)
	dfs(image, x, y-1, contour)
	dfs(image, x, y+1, contour)
}
