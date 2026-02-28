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

	// Third-party library for HTML to Markdown conversion
	md "github.com/tomkosm/html-to-markdown"
	"github.com/tomkosm/html-to-markdown/plugin"
)

// ConvertHTMLToMarkdown converts an HTML string to Markdown format
// This function is exported for use with CGO, allowing it to be called from other languages
//export ConvertHTMLToMarkdown
func ConvertHTMLToMarkdown(html *C.char) *C.char {
	// Create a new converter with default settings
	// Parameters: "" = no base URL, true = enable smart symbol conversion, nil = no additional options
	converter := md.NewConverter("", true, nil)
	
	// Use GitHub Flavored Markdown plugins for better compatibility
	converter.Use(plugin.GitHubFlavored())

	// Convert the HTML string to Markdown
	markdown, err := converter.ConvertString(C.GoString(html))
	if err != nil {
		// Error handling is commented out to prevent crash in CGO context
		// log.Fatal(err)
	}
	
	// Return the result as a C string for CGO compatibility
	return C.CString(markdown)
}

// convertHTMLToMarkdownString is a Go-friendly wrapper for HTML to Markdown conversion
// It provides the same functionality as ConvertHTMLToMarkdown but with native Go types
func convertHTMLToMarkdownString(html string) (string, error) {
	// Create converter with same settings as the exported function
	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())

	// Return both the converted markdown and any error that occurred
	return converter.ConvertString(html)
}

// main function serves as the entry point for CLI usage
func main() {
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
