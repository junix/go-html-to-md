# go-html-to-md

CLI tool and CGO shared library for converting HTML to Markdown, powered by [tomkosm/html-to-markdown](https://github.com/tomkosm/html-to-markdown) with GitHub Flavored Markdown support.

## Usage

**CLI (stdin/stdout):**

```bash
cat input.html | html-to-markdown > output.md
```

**CLI with file arguments (CGO library build):**

```bash
html-to-markdown input.html output.md
```

**As a shared library from C:**

```c
char* md = ConvertHTMLToMarkdown("<b>hello</b>");
char* md = ConvertHTMLFileToMarkdown("input.html");
```

## Build / Test / Install

```bash
just build     # go build -o html-to-markdown
just test      # quick smoke test via stdin
just install   # install binary to ~/sync/bin_<arch>/
```

## Dependencies

- Go 1.19+
- `github.com/tomkosm/html-to-markdown` -- core conversion engine with GFM plugin
