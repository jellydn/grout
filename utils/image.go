package utils

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

// CalculateResizeHeight calculates the new height for an image when resizing to a target width
// while maintaining the original aspect ratio. Returns the new height and any error encountered.
func CalculateResizeHeight(imagePath string, targetWidth int) (int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return 0, fmt.Errorf("failed to decode PNG: %w", err)
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	if originalWidth == 0 {
		return 0, fmt.Errorf("image has zero width")
	}

	aspectRatio := float64(originalHeight) / float64(originalWidth)
	newHeight := int(float64(targetWidth) * aspectRatio)

	return newHeight, nil
}

// ResizePNG resizes a PNG image to the specified width while maintaining aspect ratio.
// It reads from inputPath and writes the resized image to outputPath.
func ResizePNG(inputPath, outputPath string, targetWidth int) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input image: %w", err)
	}
	defer inputFile.Close()

	srcImg, err := png.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Calculate new dimensions
	bounds := srcImg.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	if originalWidth == 0 {
		return fmt.Errorf("image has zero width")
	}

	aspectRatio := float64(originalHeight) / float64(originalWidth)
	newHeight := int(float64(targetWidth) * aspectRatio)

	dstImg := image.NewRGBA(image.Rect(0, 0, targetWidth, newHeight))

	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, bounds, draw.Over, nil)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, dstImg); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}
