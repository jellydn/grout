package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("BIOS Generator - Libretro Info File Parser")
	fmt.Println("==========================================")
	fmt.Println()

	// Determine output path (relative to constants/ directory)
	outputPath := filepath.Join("constants", "bios_generated.go")

	// Create components
	fetcher := NewFetcher()
	parser := NewParser()
	mapper := NewMapper()
	codegen := NewCodeGenerator(outputPath)

	// Step 1: Fetch all .info files from GitHub
	fmt.Println("Step 1: Fetching .info files from GitHub...")
	infoFiles, err := fetcher.FetchAllInfoFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching info files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	// Step 2: Parse all .info files
	fmt.Println("Step 2: Parsing .info files...")
	cores, err := parser.ParseAllInfoFiles(infoFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing info files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	// Step 3: Normalize core names and build mappings
	fmt.Println("Step 3: Building core and platform mappings...")
	normalizedCores := mapper.NormalizeCoreNames(cores)
	platformToCores := mapper.BuildPlatformToCoresMap(cores)

	fmt.Printf("Found %d platforms with BIOS requirements\n", len(platformToCores))
	fmt.Println()

	// Step 4: Generate Go source code
	fmt.Println("Step 4: Generating Go source code...")
	if err := codegen.Generate(normalizedCores, platformToCores); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("✓ BIOS mappings generated successfully!")
	fmt.Printf("  - %d cores with BIOS requirements\n", len(normalizedCores))
	fmt.Printf("  - %d platform mappings\n", len(platformToCores))
	fmt.Printf("  - Output: %s\n", outputPath)
}
