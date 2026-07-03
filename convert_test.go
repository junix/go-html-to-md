package main

import (
	"strings"
	"testing"
)

// TestConvertHTMLToMarkdownStringBasics exercises the shared pure helper that
// both build targets delegate to. These are the first unit tests in the repo;
// they guard the core conversion logic without needing the CLI or CGO.
func TestConvertHTMLToMarkdownStringBasics(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string // substring expected in output
	}{
		{"bold", "<b>hello</b>", "hello"},
		{"strong", "<strong>world</strong>", "world"},
		{"paragraph", "<p>some text</p>", "some text"},
		{"heading", "<h1>Title</h1>", "Title"},
		{"link", `<a href="https://example.com">link</a>`, "link"},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := convertHTMLToMarkdownString(c.in)
			if err != nil {
				t.Fatalf("convertHTMLToMarkdownString(%q) returned unexpected error: %v", c.in, err)
			}
			if c.want != "" && !strings.Contains(got, c.want) {
				t.Errorf("convertHTMLToMarkdownString(%q) = %q, want substring %q", c.in, got, c.want)
			}
		})
	}
}

// TestConvertHTMLToMarkdownStringStripsTags confirms that surrounding HTML
// tags are removed and only the inner text survives in the markdown output.
func TestConvertHTMLToMarkdownStringStripsTags(t *testing.T) {
	got, err := convertHTMLToMarkdownString("<p>line one</p>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(got, "<p>") || strings.Contains(got, "</p>") {
		t.Errorf("output still contains <p> tags: %q", got)
	}
}

// TestNewMarkdownConverterIsStable ensures the shared factory is non-nil and
// deterministic across calls so that both entry points behave identically.
func TestNewMarkdownConverterIsStable(t *testing.T) {
	c1 := newMarkdownConverter()
	c2 := newMarkdownConverter()
	if c1 == nil || c2 == nil {
		t.Fatalf("newMarkdownConverter() returned nil")
	}
	out1, err := c1.ConvertString("<b>x</b>")
	if err != nil {
		t.Fatalf("c1.ConvertString error: %v", err)
	}
	out2, err := c2.ConvertString("<b>x</b>")
	if err != nil {
		t.Fatalf("c2.ConvertString error: %v", err)
	}
	if out1 != out2 {
		t.Errorf("converter not stable across calls: %q vs %q", out1, out2)
	}
}
