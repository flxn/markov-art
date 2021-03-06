package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func printUsage() {
	fmt.Printf("Usage: %v [flags] inputFile outputFile\n", os.Args[0])
	fmt.Println("Flags:")
	fmt.Println(" -o Markov chain of n-th order (Default: 1)")
	fmt.Println(" -horizontal Horizontal line scanning (Default: Vertical if omitted)")
	fmt.Println(" -h Print help")
	fmt.Println("Arguments:")
	fmt.Println(" inputFile Path to input image")
	fmt.Println(" outputFile Path to output image (must be *.png)")
}

func main() {
	flag.Usage = printUsage

	lookupPtr := flag.Int("o", 1, "Markov chain of n-th order")
	horizontalPtr := flag.Bool("horizontal", false, "Orientation horizontal (Default: vertical)")

	flag.Parse()

	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	inputPath := flag.Arg(0)
	outFileName := flag.Arg(1)

	file, err := os.Open(inputPath)
	handle(err)
	defer file.Close()

	fmt.Printf("File: %s\n", inputPath)

	img, _, err := image.Decode(file)
	handle(err)

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	fmt.Printf("Dimensions: %vx%v\n", width, height)

	orientation := 0
	if *horizontalPtr {
		orientation = 1
	}

	m := createMarkovChain(img, *lookupPtr, orientation)

	markovImage := m.generateImage(width, height)

	outFile, err := os.Create(outFileName)
	handle(err)
	defer outFile.Close()

	png.Encode(outFile, markovImage)
	fmt.Printf("\nImage saved as %v\n", outFileName)
}
