package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"strconv"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

//Lookup value
var LOOKUP = 2

func main() {
	if len(os.Args[1:]) < 2 {
		fmt.Printf("Usage: %v inputPath outputPath\n", os.Args[0])
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outFileName := os.Args[2]

	if len(os.Args[1:]) == 3 {
		lookupparam := os.Args[3]
		l, err := strconv.Atoi(lookupparam)
		handle(err)
		LOOKUP = l
	}

	file, err := os.Open(inputPath)
	handle(err)
	defer file.Close()

	fmt.Printf("File: %s\n", inputPath)

	img, _, err := image.Decode(file)
	handle(err)

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	fmt.Printf("Dimensions: %vx%v\n", width, height)

	m := createMarkovChain(img)

	markovImage := m.generateImage(width, height)

	outFile, err := os.Create(outFileName)
	handle(err)
	defer outFile.Close()

	png.Encode(outFile, markovImage)
	fmt.Printf("\nImage saved as %v\n", outFileName)
}
