package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func readFile(s string) (image.Image, error) {
	//TODO check file type
	file, err := os.Open(s)
	if err != nil {
		return nil, fmt.Errorf("Cannot open file: %v", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode image: %v", err)
	}
	return img, nil
}

func resizeImage(img image.Image, width uint) image.Image {
	bounds := img.Bounds()
	w, h := float64(bounds.Max.X), float64(bounds.Max.Y)
	ratio := float64(width) / w
	newHeight := uint((ratio * h) / 1.75)
	resized := resize.Resize(width, newHeight, img, resize.Lanczos3)
	return resized
}

func valueChar(val float64) string {
	var outChar string
	if val < 51 {
		outChar = "#"
	} else if val > 51 && val < 102 {
		outChar = "%"
	} else if val > 102 && val < 153 {
		outChar = "="
	} else if val > 153 && val < 204 {
		outChar = "-"
	} else if val > 204 && val < 256 {
		outChar = " "
	}
	return outChar
}

func main() {

	inPath := flag.String("i", "", "input file")
	width := flag.Uint("w", 80, "column width in output")
	flag.Parse()

	img, err := readFile(*inPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fileName := strings.TrimSuffix(*inPath, filepath.Ext(*inPath))
	smallImage := resizeImage(img, *width)
	bounds := smallImage.Bounds()
	out := make([][]string, bounds.Max.Y)

	for i := 0; i < bounds.Max.Y; i++ {
		out[i] = make([]string, *width)
	}

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			rgba := smallImage.At(x, y)
			r, g, b, _ := rgba.RGBA()
			value := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			asciiChar := valueChar(value / 256)
			out[y][x] = asciiChar
		}
	}
	outFileName := fileName + ".txt"
	outTxt, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(outTxt)
	defer outTxt.Close()
	for i := range out {
		row := fmt.Sprintf("%v\n", strings.Join(out[i], ""))
		w.WriteString(row)
	}
	w.Flush()
}
