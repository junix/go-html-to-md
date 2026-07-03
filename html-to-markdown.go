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
	}

	// Return the result as a C string for CGO compatibility
	return C.CString(markdown)
}

// main function serves as the entry point for CLI usage
func main() {
	// Print version and exit
	if len(os.Args) > 1 && (os.Args[1] == "-V" || os.Args[1] == "--version") {
		fmt.Println("html-to-markdown 0.1.0")
		return
	}

	// Display help information if requested
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("HTML to Markdown Converter")
		fmt.Println("Usage: cat input.html | html-to-markdown > output.md")
		fmt.Println("Reads HTML from stdin and outputs Markdown to stdout")
		return
	}

	// Read HTML input from stdin line by line
	reader := bufio.NewReader(os.Stdin)
	var htmlBuilder strings.Builder

	// Continue reading until EOF or error
	for {
		line, err := reader.ReadString('\n')
		htmlBuilder.WriteString(line)
		if err == io.EOF {
			// End of input reached
			break
		}
		if err != nil {
			// Handle other read errors
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
	}

	// Convert the accumulated HTML to Markdown
	markdown, err := convertHTMLToMarkdownString(htmlBuilder.String())
	if err != nil {
		// Handle conversion errors
		fmt.Fprintf(os.Stderr, "Error converting HTML to Markdown: %v\n", err)
		os.Exit(1)
	}

	// Output the converted Markdown to stdout
	fmt.Print(markdown)
}
