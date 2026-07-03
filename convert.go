// Package main provides a utility for converting HTML to Markdown format
// This program can be used both as a CLI tool and as a library via CGO
package main

import (
	// Third-party library for HTML to Markdown conversion
	md "github.com/tomkosm/html-to-markdown"
	"github.com/tomkosm/html-to-markdown/plugin"
)

// newMarkdownConverter builds a converter configured the same way for every
// entry point (CLI stdin mode, CLI file mode, and the CGO exports): no base
// URL, smart symbol conversion enabled, with the GitHub Flavored Markdown
// plugin applied. Keeping this in one place guarantees the CLI and shared
// library always produce identical output.
func newMarkdownConverter() *md.Converter {
	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())
	return converter
}

// convertHTMLToMarkdownString is a Go-friendly wrapper for HTML to Markdown
// conversion. It is shared by both build targets (the stdin/stdout CLI in
// html-to-markdown.go and the CGO library in html-to-markdown-lib.go) so that
// the conversion logic exists in exactly one place.
func convertHTMLToMarkdownString(html string) (string, error) {
	// Create converter with default settings shared across all entry points
	converter := newMarkdownConverter()

	// Return both the converted markdown and any error that occurred
	return converter.ConvertString(html)
}
