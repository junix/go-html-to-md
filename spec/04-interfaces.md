# 04 Interfaces

## Architecture / Layering

```
                   ┌──────────────────────────────┐
                   │        Entry Surfaces         │
                   │  ┌─────────┐  ┌────────────┐ │
                   │  │  CLI    │  │  C FFI     │ │
                   │  │(pipe)   │  │(shared lib)│ │
                   │  └────┬────┘  └─────┬──────┘ │
                   │       │             │        │
                   │  ┌────┴────┐  ┌─────┴──────┐ │
                   │  │ CLI     │  │ CLI        │ │
                   │  │(file)   │  │(file+FFI)  │ │
                   │  └────┬────┘  └─────┬──────┘ │
                   └───────┼─────────────┼────────┘
                           │             │
                   ┌───────▼─────────────▼────────┐
                   │    Conversion Engine          │
                   │  (HTML-to-Markdown + GFM)     │
                   └──────────────────────────────┘
```

All entry surfaces converge on a single conversion engine with identical configuration. There is no user-accessible middleware, plugin registration, or configuration layer between the entry surface and the engine.

## CLI Interface

### Invocation Syntax (ABNF)

```abnf
cli-invocation  = pipe-mode / file-mode / help-mode

help-mode       = binary-name SP help-flag
help-flag       = "-h" / "--help"

pipe-mode       = binary-name  ; reads HTML from stdin, writes Markdown to stdout

file-mode       = binary-name SP input-file [SP output-file]

binary-name     = 1*CHAR       ; the executable name
input-file      = 1*CHAR       ; path to HTML file
output-file     = 1*CHAR       ; path to write Markdown output
```

### CLI Behavior

| Surface | Input Source | Output Destination | Exit 0 | Exit 1 |
|---------|-------------|-------------------|--------|--------|
| Pipe mode | stdin (until EOF) | stdout | Conversion succeeds | Read error or conversion error |
| File mode (input only) | `input-file` | stdout | Conversion succeeds | Read/convert error |
| File mode (input + output) | `input-file` | `output-file` | Conversion succeeds | Read/write/convert error |
| Help mode | none | stdout (usage text) | always | never |

### Help Output

When invoked with `-h` or `--help` as the first argument, the tool MUST print the following to stdout and exit with status 0:

```
HTML to Markdown Converter
Usage: cat input.html | html-to-markdown > output.md
Reads HTML from stdin and outputs Markdown to stdout
```

This help text is only produced by the pipe-mode binary. The file-mode binary produces a different usage string:

```
Usage: html-to-markdown-lib <input_file> [output_file]
       cat file.html | html-to-markdown-lib
```

### Stdin Detection (File-Mode Binary)

The file-mode binary distinguishes between "no arguments and piped input" vs "no arguments and no piped input":

```
(condition,                        action)
(no args + stdin is a pipe/device, read stdin, convert, write stdout)
(no args + stdin is a terminal,    print usage to stdout, exit 0)
```

## C FFI Interface

### Shared Library Exports

The shared library exports exactly two functions:

```abnf
ffi-export = ffi-string-conv / ffi-file-conv

ffi-string-conv = "ConvertHTMLToMarkdown" "(" html-param ")" return-type
ffi-file-conv   = "ConvertHTMLFileToMarkdown" "(" filepath-param ")" return-type

html-param     = "char* html"
filepath-param = "char* filepath"
return-type    = "char*"
```

### Function Contracts

#### ConvertHTMLToMarkdown

| Aspect | Contract |
|--------|----------|
| Input | A null-terminated C string containing HTML |
| Output on success | A null-terminated C string containing Markdown |
| Output on failure | A null-terminated C string: `"Error: Failed to convert HTML to Markdown"` |
| Memory ownership | The returned C string is allocated by the library. The caller MUST free it. |
| Process safety | MUST NOT abort or terminate the calling process on any input |

#### ConvertHTMLFileToMarkdown

| Aspect | Contract |
|--------|----------|
| Input | A null-terminated C string containing a file path |
| Output on success | A null-terminated C string containing Markdown (file contents converted) |
| Output on file-read failure | A null-terminated C string: `"Error: Failed to read file"` |
| Output on conversion failure | A null-terminated C string: `"Error: Failed to convert HTML to Markdown"` |
| Memory ownership | The returned C string is allocated by the library. The caller MUST free it. |
| Process safety | MUST NOT abort or terminate the calling process on any input |

### C Header Guard

The shared library's C header wraps all declarations in `extern "C"` blocks for C++ compatibility, and uses include guards.

## Extension Points

There are no user-accessible extension points. The conversion engine's plugin set is fixed to GFM rules at build time. Future versions MAY introduce configuration, but this specification defines no such mechanism.
