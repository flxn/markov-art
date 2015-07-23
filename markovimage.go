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

// Alpha channel support
// Default false saves some memory
var ALPHA = false

type markovImage struct {
	lookup      int
	orientation int
	colorMap    map[string][]color.Color
}

func colorToCString(col color.Color) string {
	r, g, b, a := col.RGBA()
	colorString := strconv.Itoa(int(r)) + strconv.Itoa(int(g)) + strconv.Itoa(int(b))
	if ALPHA {
		colorString += strconv.Itoa(int(a))
	}
	return colorString
}

func createMarkovChain(img image.Image, lookup int, orientation int) markovImage {
	mi := new(markovImage)
	(*mi).lookup = lookup
	(*mi).orientation = orientation

	chain := make(map[string][]color.Color)

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	totalPixelCount := width * height
	pixelCount := 0

	previousColors := make([]string, mi.lookup)

	first := width
	second := height

	if mi.orientation == 1 {
		first = height
		second = width
	}

	for x := 0; x < first; x++ {
		for y := 0; y < second; y++ {
			var newColor color.Color
			if orientation == 0 {
				newColor = img.At(x, y)
			} else {
				newColor = img.At(y, x)
			}

			key := str.Join(previousColors, "|")

			chain[key] = append(chain[key], newColor)

			previousColors = previousColors[1:]
			previousColors = append(previousColors, colorToCString(newColor))
			if pixelCount%1000 == 0 {
				fmt.Printf("\rGenerating Markov Chain: %.0f%%", float32(pixelCount)/float32(totalPixelCount)*100)
			}
			pixelCount++
		}
	}
	fmt.Printf("\n")
	(*mi).colorMap = chain
	return *mi
}

func (mi markovImage) generateImage(width int, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	rand.Seed(time.Now().UTC().UnixNano())
	totalPixelCount := width * height
	pixelCount := 0

	previousColors := make([]string, mi.lookup)

	first := width
	second := height

	if mi.orientation == 1 {
		first = height
		second = width
	}

	for x := 0; x < first; x++ {
		for y := 0; y < second; y++ {
			key := str.Join(previousColors, "|")
			if len(mi.colorMap[key]) == 0 {
				fmt.Printf("\nNo color found for current key. Resetting past state information...\n")
				previousColors = make([]string, mi.lookup)
				continue
			}
			newColor := mi.colorMap[key][rand.Intn(len(mi.colorMap[key]))]

			if mi.orientation == 0 {
				img.Set(x, y, newColor)
			} else {
				img.Set(y, x, newColor)
			}
			previousColors = previousColors[1:]
			previousColors = append(previousColors, colorToCString(newColor))
			if pixelCount%1000 == 0 {
				fmt.Printf("\rGenerating Image: %.0f%%", float32(pixelCount)/float32(totalPixelCount)*100)
			}
			pixelCount++
		}
	}

	return img
}
