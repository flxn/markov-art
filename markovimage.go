package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strconv"
	str "strings"
	"time"
)

type colorMap map[string][]color.Color

func colorToCString(col color.Color) string {
	r, g, b, a := col.RGBA()
	colorString := strconv.Itoa(int(r)) + strconv.Itoa(int(g)) + strconv.Itoa(int(b)) + strconv.Itoa(int(a))
	return colorString
}

func createMarkovChain(img image.Image) colorMap {
	chain := make(colorMap)
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	previousColors := make([]string, LOOKUP)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			newColor := img.At(x, y)

			key := str.Join(previousColors, "|")

			chain[key] = append(chain[key], newColor)

			previousColors = previousColors[1:]
			previousColors = append(previousColors, colorToCString(newColor))

		}
	}

	return chain
}

func (chain colorMap) generateImage(width int, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	rand.Seed(time.Now().UTC().UnixNano())
	pixelCount := width * height

	previousColors := make([]string, LOOKUP)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			key := str.Join(previousColors, "|")
			if len(chain[key]) == 0 {
				fmt.Printf("\nNo next color found for key %v.\nResetting previous color information...\n", key)
				previousColors = make([]string, LOOKUP)
				continue
			}
			newColor := chain[key][rand.Intn(len(chain[key]))]
			img.Set(x, y, newColor)
			previousColors = previousColors[1:]
			previousColors = append(previousColors, colorToCString(newColor))
			fmt.Printf("\r%.0f%% done", float32(x*y)/float32(pixelCount)*100)
		}
	}

	return img
}

/*
func (mi markovImage) generate() {
	for x := 0; x < theImage.width; x++ {
		for y := 0; y < theImage.height; y++ {
			mi.dict[theImage.imagePoints[x][y].color]["above"] = theImage.imagePoints[x][y].getColorAbove()
			mi.dict[theImage.imagePoints[x][y].color]["below"] = theImage.imagePoints[x][y].getColorBelow()
			mi.dict[theImage.imagePoints[x][y].color]["right"] = theImage.imagePoints[x][y].getColorRight()
			mi.dict[theImage.imagePoints[x][y].color]["left"] = theImage.imagePoints[x][y].getColorLeft()
		}
	}

	fmt.Print(len(mi.dict))

	theImage.dict = mi.dict
}
*/
