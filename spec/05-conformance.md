# 05 Conformance

## Examples

### E.1 Pipe Mode — Basic Conversion

Given HTML input on stdin:

```html
<p>This is <strong>bold text</strong> and <em>italic text</em>.</p>
```

The tool writes to stdout:

```markdown
This is **bold text** and _italic text_.
```

### E.2 Pipe Mode — Help Flag

Invoking with `--help` produces:

```
HTML to Markdown Converter
Usage: cat input.html | html-to-markdown > output.md
Reads HTML from stdin and outputs Markdown to stdout
```

The process exits with status 0.

### E.3 File Mode — Conversion with Output File

Given an HTML file at `input.html` containing:

```html
<h1>Title</h1><p>Hello world</p>
```

Invoking `html-to-markdown-lib input.html output.md` writes to `output.md`:

```markdown
# Title

Hello world
```

### E.4 File Mode — No Pipe, No Args

Invoking `html-to-markdown-lib` with no arguments and stdin connected to a terminal produces usage text on stdout and exits with status 0.

### E.5 C FFI — String Conversion Failure

Calling `ConvertHTMLToMarkdown` with input that triggers a conversion engine error returns the C string `"Error: Failed to convert HTML to Markdown"`.

### E.6 C FFI — File Not Found

Calling `ConvertHTMLFileToMarkdown` with a nonexistent file path returns the C string `"Error: Failed to read file"`.

### E.7 GFM Table Conversion

Given HTML input:

```html
<table><thead><tr><th>A</th><th>B</th></tr></thead>
<tbody><tr><td>1</td><td>2</td></tr></tbody></table>
```

The output MUST contain a GFM table:

```markdown
| A | B |
| --- | --- |
| 1 | 2 |
```

### E.8 Nested List Conversion

Given HTML input:

```html
<ul><li>Item 1</li><li>Item 2<ul><li>Nested</li></ul></li></ul>
```

The output MUST contain nested list syntax:

```markdown
- Item 1
- Item 2
  - Nested
```

### E.9 Blockquote Conversion

Given HTML input:

```html
<blockquote><p>Quote paragraph one.</p><p>Quote paragraph two.</p></blockquote>
```

The output MUST contain blockquote syntax with blank-line paragraph separation:

```markdown
> Quote paragraph one.
>
> Quote paragraph two.
```

### E.10 Image Tag Conversion

Given HTML input:

```html
<img src="https://example.com/img.png" alt="Example">
```

The output MUST contain Markdown image syntax:

```markdown
![Example](https://example.com/img.png)
```

Note: HTML attributes not representable in standard Markdown image syntax (e.g., `width`) are not preserved in the output.

## Definition of Done

### A. Pipe-Based CLI

| ID | Behavior | Observable Result | Status |
|----|----------|-------------------|--------|
| A.1 | Stdin HTML is read until EOF | Conversion output appears on stdout only after stdin is closed | [U] Currently untested |
| A.2 | Stdin read error | Diagnostic written to stderr; process exits with non-zero status | [U] Currently untested |
| A.3 | Conversion error | Diagnostic written to stderr; process exits with non-zero status | [U] Currently untested |
| A.4 | Successful conversion | Markdown output written to stdout | [T] Verified by `just test` smoke test |
| A.5 | Help flag (`-h` or `--help`) | Usage text printed to stdout; exit status 0 | [U] Currently untested |

### B. File-Argument CLI

| ID | Behavior | Observable Result | Status |
|----|----------|-------------------|--------|
| B.1 | Input file + output file | Markdown written to the specified output file | [U] Currently untested |
| B.2 | Input file, no output file | Markdown written to stdout | [U] Currently untested |
| B.3 | Input file not found | Diagnostic to stderr; exit status 1 | [U] Currently untested |
| B.4 | Output file write failure | Diagnostic to stderr; exit status 1 | [U] Currently untested |
| B.5 | No args + piped stdin | Reads stdin, converts, writes stdout | [U] Currently untested |
| B.6 | No args + terminal stdin | Usage text printed to stdout; exit status 0 | [U] Currently untested |
| B.7 | Conversion error | Diagnostic to stderr; exit status 1 | [U] Currently untested |

### C. C FFI

| ID | Behavior | Observable Result | Status |
|----|----------|-------------------|--------|
| C.1 | `ConvertHTMLToMarkdown` with valid HTML | Returns C string containing Markdown | [U] Currently untested |
| C.2 | `ConvertHTMLToMarkdown` with conversion error | Returns C string starting with `"Error:"` | [U] Currently untested |
| C.3 | `ConvertHTMLFileToMarkdown` with valid file | Returns C string containing Markdown | [U] Currently untested |
| C.4 | `ConvertHTMLFileToMarkdown` with nonexistent file | Returns C string `"Error: Failed to read file"` | [U] Currently untested |
| C.5 | C FFI does not abort on any input | Process continues after error return | [U] Currently untested |
| C.6 | Returned C string memory is caller-owned | Caller can free the returned pointer | [U] Currently untested |

### D. Conversion Correctness

| ID | Behavior | Observable Result | Status |
|----|----------|-------------------|--------|
| D.1 | Bold (`<strong>`, `<b>`) | Output contains `**text**` | [T] Verified by sample HTML conversion |
| D.2 | Italic (`<em>`, `<i>`) | Output contains `_text_` or `*text*` | [T] Verified by sample HTML conversion |
| D.3 | Inline code (`<code>`) | Output contains `` `text` `` | [T] Verified by sample HTML conversion |
| D.4 | Headings (`<h1>` through `<h6>`) | Output contains `#`-prefixed heading syntax | [T] Verified by sample HTML conversion |
| D.5 | Unordered lists (`<ul>`, `<li>`) | Output contains `- item` syntax with nesting | [T] Verified by sample HTML conversion |
| D.6 | Ordered lists (`<ol>`, `<li>`) | Output contains `1. item` syntax | [T] Verified by sample HTML conversion |
| D.7 | Links (`<a href="...">`) | Output contains `[text](url)` | [T] Verified by sample HTML conversion |
| D.8 | Images (`<img src="..." alt="...">`) | Output contains `![alt](src)` | [T] Verified by sample HTML conversion |
| D.9 | Blockquotes (`<blockquote>`) | Output contains `>`-prefixed lines | [T] Verified by sample HTML conversion |
| D.10 | Code blocks (`<pre><code>`) | Output contains fenced code blocks | [T] Verified by sample HTML conversion |
| D.11 | Tables (`<table>`) | Output contains GFM pipe table syntax | [T] Verified by sample HTML conversion |
| D.12 | Horizontal rules (`<hr>`) | Output contains horizontal rule syntax | [T] Verified by sample HTML conversion |

### E. Converter Configuration Invariant

| ID | Behavior | Observable Result | Status |
|----|----------|-------------------|--------|
| E.1 | All entry points produce identical output for identical input | Pipe CLI, file CLI, and C FFI return the same Markdown for the same HTML | [U] Currently untested |
| E.2 | Smart symbol conversion is enabled | Typographic entities in HTML are converted to ASCII equivalents | [U] Currently untested |
| E.3 | No base URL is set | Relative links in HTML are not resolved to absolute URLs | [U] Currently untested |
