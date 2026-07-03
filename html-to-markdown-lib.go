// Package main provides a utility for converting HTML to Markdown format
// This program can be used both as a CLI tool and as a library via CGO
package main

import (
	"C" // Required for CGO integration
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ConvertHTMLToMarkdown converts an HTML string to Markdown format
// This function is exported for use with CGO, allowing it to be called from other languages
//
//export ConvertHTMLToMarkdown
func ConvertHTMLToMarkdown(html *C.char) *C.char {
	// Build a converter configured with GitHub Flavored Markdown (shared helper)
	converter := newMarkdownConverter()

	// Convert the HTML string to Markdown
	markdown, err := converter.ConvertString(C.GoString(html))
	if err != nil {
		// Error handling is commented out to prevent crash in CGO context
		// log.Fatal(err)
		return C.CString("Error: Failed to convert HTML to Markdown")
	}

	// Return the result as a C string for CGO compatibility
	return C.CString(markdown)
}

// ConvertHTMLFileToMarkdown converts HTML content from a file to Markdown and returns it as a C string
//
//export ConvertHTMLFileToMarkdown
func ConvertHTMLFileToMarkdown(filepath *C.char) *C.char {
	content, err := os.ReadFile(C.GoString(filepath))
	if err != nil {
		return C.CString("Error: Failed to read file")
	}

	return ConvertHTMLToMarkdown(C.CString(string(content)))
}

// convertHTMLFileToMarkdown is a Go helper function to convert HTML content from a file
func convertHTMLFileToMarkdown(inputFile string, outputFile string) error {
	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	// Convert HTML to Markdown
	markdown, err := convertHTMLToMarkdownString(string(content))
	if err != nil {
		return fmt.Errorf("error converting HTML to Markdown: %w", err)
	}

	// Write to output file or stdout
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(markdown), 0644)
		if err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
	} else {
		fmt.Print(markdown)
	}

	return nil
}

// Required for building as a C shared library
func main() {
	// When used as a CLI tool
	if len(os.Args) > 1 {
		// Simple CLI mode
		inputFile := os.Args[1]
		outputFile := ""

		// Check if output file is specified
		if len(os.Args) > 2 {
			outputFile = os.Args[2]
		}

		err := convertHTMLFileToMarkdown(inputFile, outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else if !isPipedInput() {
		// Print usage if no arguments and no piped input
		fmt.Println("Usage: html-to-markdown-lib <input_file> [output_file]")
		fmt.Println("       cat file.html | html-to-markdown-lib")
	} else {
		// Process HTML from stdin
		reader := bufio.NewReader(os.Stdin)
		var html strings.Builder

		// Read all stdin content
		for {
			line, err := reader.ReadString('\n')
			html.WriteString(line)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				os.Exit(1)
			}
		}

		// Convert and print to stdout
		markdown, err := convertHTMLToMarkdownString(html.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting HTML: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(markdown)
	}
}

// Helper function to check if input is piped
func isPipedInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}
