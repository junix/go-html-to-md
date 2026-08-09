package main

import (
	"strings"
	"testing"
)

// TestConvertHTMLToMarkdownStringExact pins the EXACT Markdown output the
// shared converter produces for each HTML construct. An earlier version of
// this test only asserted the inner text survived as a substring (checking for
// "hello" given "<b>hello</b>"), which hides regressions that drop the
// formatting markers; the "empty" case asserted nothing at all. Every case
// here compares the full string so the formatting contract is locked down.
// All values were captured from the converter's real behavior.
func TestConvertHTMLToMarkdownStringExact(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string // exact Markdown output
	}{
		// Inline formatting: the marker chars are the contract.
		{"bold", "<b>hello</b>", "**hello**"},
		{"strong", "<strong>world</strong>", "**world**"},
		{"em", "<em>ital</em>", "_ital_"},
		{"i", "<i>ital</i>", "_ital_"},
		{"strikethrough_del", "<del>struck</del>", "~~struck~~"},
		{"strikethrough_s", "<s>struck</s>", "~~struck~~"},

		// Headings: level maps to N hash marks + a single space.
		{"h1", "<h1>Title</h1>", "# Title"},
		{"h2", "<h2>Sub</h2>", "## Sub"},
		{"h6", "<h6>Deep</h6>", "###### Deep"},
		{"uppercase_tag", "<H1>MIXED CASE TAG</H1>", "# MIXED CASE TAG"},

		// Paragraph / plain text passthrough.
		{"paragraph", "<p>some text</p>", "some text"},
		{"plain_text", "plain text", "plain text"},

		// Link & image: href/src must survive, not just the link text.
		{"link", `<a href="https://example.com">link</a>`, "[link](https://example.com)"},
		{"image", `<img src="x.png" alt="pic">`, "![pic](x.png)"},

		// Lists: marker style and numeric ordering are the contract.
		{"unordered_list", "<ul><li>one</li><li>two</li></ul>", "- one\n- two"},
		{"ordered_list", "<ol><li>a</li><li>b</li></ol>", "1. a\n2. b"},

		// Code: inline backticks vs GFM fenced block.
		{"inline_code", "<code>x</code>", "`x`"},
		{"fenced_code", "<pre><code>line1\nline2</code></pre>", "```\nline1\nline2\n```"},

		// Block-level constructs.
		{"blockquote", "<blockquote>quoted</blockquote>", "> quoted"},
		{"hr", "<hr>", "* * *"},

		// Multiple blocks are separated by a blank line.
		{"two_paragraphs", "<p>a</p><p>b</p>", "a\n\nb"},

		// Nesting: inner formatting is applied within the outer block.
		{"nested", "<div>nested <b>bold</b></div>", "nested **bold**"},

		// GFM table rendered only because plugin.GitHubFlavored() is applied.
		{"gfm_table",
			"<table><thead><tr><th>A</th><th>B</th></tr></thead><tbody><tr><td>1</td><td>2</td></tr></tbody></table>",
			"| A | B |\n| --- | --- |\n| 1 | 2 |"},

		// Empty / whitespace-only input: contract is empty output (not " " or "\n").
		{"empty", "", ""},
		{"whitespace_only", "   ", ""},
		{"empty_p", "<p></p>", ""},
		{"empty_h1", "<h1></h1>", ""},

		// HTML entity decoding (smart conversion enabled).
		{"entities", "&amp;&lt;&gt;", "&<>"},
		{"entity_in_text", "<p>1 &lt; 2 &amp;&amp; 3 &gt; 2</p>", "1 < 2 && 3 > 2"},

		// Unicode must pass through unchanged.
		{"vietnamese", "<p>tôi là bàn</p>", "tôi là bàn"},
		{"emoji", "<p>emoji 😀 test</p>", "emoji 😀 test"},

		// Whitespace within a block is collapsed and trimmed.
		{"collapse_spaces", "<p>multiple   spaces</p>", "multiple spaces"},
		{"trim_block", "<p>  spaces  </p>", "spaces"},

		// Script/style content must be dropped entirely.
		{"script_dropped", "<script>alert(1)</script>", ""},
		{"style_dropped", "<style>.x{}</style>", ""},

		// <br> inside a paragraph becomes a paragraph break.
		{"br_separator", "<p>line1<br>line2</p>", "line1\n\nline2"},

		// Embedded newlines within a block are preserved.
		{"embedded_newlines", "<p>line\nwith\nnewlines</p>", "line\nwith\nnewlines"},

		// Malformed / non-HTML input is tolerated, never errored.
		{"unclosed_tag", "<b>", ""},
		{"lonely_close_tag", "</b>", ""},
		{"not_html", "<<>>", "<<>>"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := convertHTMLToMarkdownString(c.in)
			// Contract: a strings.Reader never errors, so ConvertString must
			// never return an error for any string input. Pin that invariant
			// so a future change introducing spurious errors is caught.
			if err != nil {
				t.Fatalf("convertHTMLToMarkdownString(%q) returned unexpected error: %v", c.in, err)
			}
			if got != c.want {
				t.Errorf("convertHTMLToMarkdownString(%q)\n  got:  %q\n  want: %q", c.in, got, c.want)
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

// TestNewMarkdownConverterAppliesGitHubFlavored proves the factory wires the
// GitHub Flavored Markdown plugin (the converter.Use(plugin.GitHubFlavored())
// line in newMarkdownConverter). A <table> only renders as a GFM table, a
// <pre><code> only becomes a fenced block, and <del> only becomes ~~...~~
// because that plugin is applied. Without it the output degrades to plain
// inline text, so these observables are the contract-level proof that the
// plugin has not been dropped.
func TestNewMarkdownConverterAppliesGitHubFlavored(t *testing.T) {
	c := newMarkdownConverter()
	if c == nil {
		t.Fatal("newMarkdownConverter() returned nil")
	}

	t.Run("table_is_gfm", func(t *testing.T) {
		got, err := c.ConvertString("<table><tr><td>1</td></tr></table>")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// The GFM separator row is the signature that a table was rendered.
		if !strings.Contains(got, "| --- |") {
			t.Errorf("table did not render as GFM (no separator row): %q", got)
		}
	})

	t.Run("code_block_is_fenced", func(t *testing.T) {
		got, err := c.ConvertString("<pre><code>x</code></pre>")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Fenced code is delimited by triple backticks on their own lines.
		if !strings.HasPrefix(got, "```\n") || !strings.HasSuffix(got, "\n```") {
			t.Errorf("code block did not render as fenced: %q", got)
		}
	})

	t.Run("strikethrough", func(t *testing.T) {
		got, err := c.ConvertString("<del>x</del>")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "~~x~~" {
			t.Errorf("strikethrough not GFM-rendered: got %q, want ~~x~~", got)
		}
	})
}
