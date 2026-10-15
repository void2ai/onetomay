package unpacker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// markerRegex matches lines like: ##-## path/to/file.txt ##-##
var markerRegex = regexp.MustCompile(`^##-##\s+(.+?)\s+##-##$`)

// maxLineSize caps a single input line at 16 MiB. The bufio.Scanner default
// (64 KiB) is too small for packed files with long lines such as minified
// assets or base64 blobs.
const maxLineSize = 16 * 1024 * 1024

// FileResult holds the outcome of writing one extracted file.
type FileResult struct {
	Path    string
	Success bool
	Error   error
}

// Unpack reads a packed stream from r, parses file blocks, and writes them to
// the file system under baseDir. It returns a slice of results, one per
// extracted file.
func Unpack(r io.Reader, baseDir string) []FileResult {
	var results []FileResult
	var currentPath string
	var currentLines []string

	// flush writes the accumulated content to disk for the current file.
	flush := func() {
		if currentPath == "" {
			return
		}

		// Trim blank lines adjacent to the markers. All other whitespace
		// (indentation, internal blank lines) is preserved.
		lines := trimBlankLines(currentLines)
		content := ""
		if len(lines) > 0 {
			content = strings.Join(lines, "\n") + "\n"
		}

		// Resolve the output path, rejecting anything that would escape baseDir.
		fullPath, err := safeJoin(baseDir, currentPath)
		if err != nil {
			results = append(results, FileResult{Path: currentPath, Error: err})
			return
		}

		// Ensure the parent directory exists.
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			results = append(results, FileResult{
				Path:  currentPath,
				Error: fmt.Errorf("failed to create directory %s: %w", dir, err),
			})
			return
		}

		// Write the file.
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			results = append(results, FileResult{
				Path:  currentPath,
				Error: fmt.Errorf("failed to write file %s: %w", fullPath, err),
			})
			return
		}

		results = append(results, FileResult{Path: currentPath, Success: true})
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineSize)

	seenMarker := false
	reportedPreamble := false
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()

		// Check if this line is a marker.
		if matches := markerRegex.FindStringSubmatch(line); matches != nil {
			// We found a new marker, so flush the previous file first.
			flush()

			// Start a new file block.
			currentPath = matches[1]
			currentLines = nil
			seenMarker = true
			continue
		}

		// Non-blank content before the first marker is invalid input.
		if !seenMarker && !reportedPreamble && strings.TrimSpace(line) != "" {
			results = append(results, FileResult{
				Error: fmt.Errorf("line %d: content found before the first marker", lineNo),
			})
			reportedPreamble = true
		}

		// Accumulate content lines.
		currentLines = append(currentLines, line)
	}

	// Flush the last file after the loop ends.
	flush()

	// Check for scanner errors.
	if err := scanner.Err(); err != nil {
		results = append(results, FileResult{
			Error: fmt.Errorf("input read error: %w", err),
		})
	}

	return results
}

// safeJoin joins baseDir and name and verifies the result stays inside
// baseDir. Absolute paths and paths that escape baseDir via ".." are rejected.
func safeJoin(baseDir, name string) (string, error) {
	if baseDir == "" {
		baseDir = "."
	}
	if filepath.IsAbs(name) {
		return "", fmt.Errorf("invalid path %q: absolute paths are not allowed", name)
	}
	full := filepath.Join(baseDir, name)
	rel, err := filepath.Rel(baseDir, full)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid path %q: escapes the output directory", name)
	}
	return full, nil
}

// trimBlankLines removes leading and trailing blank (empty or whitespace-only)
// lines, preserving everything in between.
func trimBlankLines(lines []string) []string {
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	return lines[start:end]
}
