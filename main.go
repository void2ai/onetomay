package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"onetomany/unpacker"
)

func main() {
	// Define flags.
	outDir := flag.String("out", "", "Base output directory for extracted files (default: current directory)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: onetomany <input-file|-> [-out <dir>]")
		fmt.Println("  Reads a packed text file and extracts individual files.")
		fmt.Println("  Use '-' as the input file to read from standard input.")
		fmt.Println()
		fmt.Println("Marker Specification")
		fmt.Println("##-## [FILE_PATH] ##-##")
		fmt.Println("- Example: ##-## out/config.go ##-##")
		os.Exit(1)
	}

	inputPath := flag.Arg(0)

	// Files are extracted relative to the current directory unless -out is given.
	baseDir := "."
	if *outDir != "" {
		baseDir = *outDir
	}

	// Open the input: a file path, or standard input when the path is "-".
	var input io.Reader
	if inputPath == "-" {
		input = os.Stdin
	} else {
		file, err := os.Open(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot open file %s: %v\n", inputPath, err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	// Run the unpacker.
	results := unpacker.Unpack(input, baseDir)

	// Print results.
	successCount := 0
	failCount := 0

	for _, r := range results {
		if r.Path == "" {
			// This is a general input error.
			fmt.Fprintf(os.Stderr, "Error: %v\n", r.Error)
			failCount++
			continue
		}

		if r.Success {
			fmt.Printf("Created: %s\n", filepath.Join(baseDir, r.Path))
			successCount++
		} else {
			fmt.Fprintf(os.Stderr, "Failed: %s - %v\n", r.Path, r.Error)
			failCount++
		}
	}

	// Print summary.
	fmt.Printf("\nDone. %d file(s) created, %d error(s).\n", successCount, failCount)

	if failCount > 0 {
		os.Exit(1)
	}
}
