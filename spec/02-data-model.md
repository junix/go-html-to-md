# 02 Data Model

## Core Data Shapes

### HTML Input

An opaque UTF-8 text string containing zero or more HTML elements. The input MAY be empty. The input MAY be malformed HTML; the conversion engine's HTML parser applies best-effort recovery. The specification does not mandate specific behavior for malformed input beyond "the conversion engine does not abort."

### Markdown Output

A UTF-8 text string representing the converted Markdown. The output format is GitHub Flavored Markdown. The output MAY be empty (e.g., for HTML input containing only whitespace or only non-renderable elements).

### File Path

A UTF-8 string referencing a filesystem location. Used as input file path (pointing to an HTML file) or output file path (where Markdown output will be written). No validation of file extension is performed.

### Converter Configuration

The converter is instantiated once per conversion invocation with the following fixed parameters:

| Parameter | Value | Effect |
|-----------|-------|--------|
| Base URL | empty | Relative URLs in the HTML are not resolved |
| Smart symbol conversion | enabled | Common typographic symbols (e.g., em-dash, curly quotes) are converted to their ASCII equivalents in the Markdown output |
| Plugin set | GFM rules | Tables, strikethrough, task lists, and autolinks are handled by GFM-specific rendering rules |

No user-configurable options exist. The converter configuration is not stateful across invocations; a fresh converter is created for each conversion call.

## Lifecycle

```
(HTML Input)
     |
     v
[Create Converter] --(fixed config)--> [Converter Instance]
     |
     v
[Convert] --(HTML string)--> (Markdown Output)
     |
     v
[Destroy] (implicit; converter is not reused)
```

The converter has no persistent state. Each invocation creates, uses, and discards a converter independently.

## Relationships

- One HTML input produces exactly one Markdown output (deterministic for a given conversion engine version).
- A file path resolves to exactly one HTML input (the file contents at read time).
- The shared library's C FFI wraps the same conversion pipeline; the C string boundary is a marshalling concern, not a semantic distinction.
