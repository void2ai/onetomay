# onetomany

A CLI tool that unpacks a single text file containing multiple file payloads into individual files on disk.

## Input Format

Files are separated by block markers:

```
##-## path/to/file.txt ##-##

file content goes here

##-## another/file.go ##-##

more content
```

Example file (`example.otm`):

```
##-## out/config.go ##-##

package main

import "fmt"

func main() {
	fmt.Println("Hello from config!")
}

##-## out/utils/helper.go ##-##

package utils

func Add(a, b int) int {
	return a + b
}

##-## out/README.md ##-##

# Example Output

This file was extracted by the unpacker.
```

## Usage

```bash
# Build
go build -o onetomany .

# Run: extract relative to the current directory
./onetomany example.otm

# Extract under a specific base directory
./onetomany example.otm -out build

# Read from standard input
cat example.otm | ./onetomany -
```

Output:

```
Created: out/config.go
Created: out/utils/helper.go
Created: out/README.md

Done. 3 file(s) created, 0 error(s).
```

Marker paths must be relative and stay inside the output directory: absolute
paths and `..` traversal are rejected with an error.

## Architecture

- `unpacker/unpacker.go` — Core parsing and file writing logic. Accepts any `io.Reader`, making it reusable for Web APIs or GUIs.
- `main.go` — CLI entry point. Handles arguments, file opening, and result display.

## In prompt

```
Put all the code files in one code block with:

Marker Specification

files are defined by a starting marker line. The marker format is:

##-## [FILE_PATH] ##-##

- Example: ##-## out/config.go ##-##
```
