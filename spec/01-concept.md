# 01 Concept

## Scope

This specification covers a command-line tool and a C-compatible shared library for converting HTML text to GitHub Flavored Markdown. The scope includes three entry surfaces: a pipe-based CLI, a file-argument CLI, and a C FFI for shared-library callers. An auxiliary Python script for fetching HTML from URLs is documented as a companion utility but is not part of the core conversion contract.

## Problem Statement

Consumers need a reliable, reproducible way to transform arbitrary HTML documents into equivalent Markdown text, preserving structure (headings, lists, tables, code blocks, links, images, blockquotes, horizontal rules) and inline formatting (bold, italic, inline code). The tool MUST produce output compatible with GitHub Flavored Markdown rendering.

## Goals

1. Convert well-formed HTML to equivalent Markdown via a pipe-based CLI (stdin to stdout).
2. Convert an HTML file to Markdown via file-argument CLI, writing to a specified output file or stdout.
3. Expose conversion as C-callable functions via a shared library, returning Markdown as a C string.
4. Enable GitHub Flavored Markdown extensions (tables, task lists, strikethrough, autolinks) in all conversion paths.

## Non-Goals

1. This specification does not guarantee round-trip fidelity: `markdown_to_html(html_to_markdown(X)) = X` is not required.
2. URL fetching is not a core conversion responsibility; the companion Python script is informative, not normative.
3. There is no streaming or partial-conversion mode. The entire HTML input MUST be consumed before conversion begins.
4. There is no configuration or plugin system exposed to end users beyond the built-in GFM support.
5. Memory management of C strings returned by the shared library is the caller's responsibility; this specification does not define a free function.

## Design Principles

1. **Single conversion engine**: All entry surfaces (CLI, file-args, C FFI) MUST use the same conversion configuration and plugin set, ensuring identical output for identical input regardless of entry point.
2. **Fail loudly in CLI, fail softly in FFI**: CLI entry points MUST exit with a non-zero status on conversion failure. C FFI functions MUST return a human-readable error string (prefixed with "Error:") rather than aborting the process.
3. **No base URL**: The converter is initialized without a base URL, meaning relative links in the input HTML are not resolved.

## Dependencies

| Capability | Abstract description |
|------------|---------------------|
| HTML-to-Markdown conversion engine | Accepts an HTML string; returns a Markdown string. Supports extensible rendering rules for specific HTML elements. |
| GitHub Flavored Markdown plugin | A rule set that extends the conversion engine to handle GFM-specific constructs (tables, strikethrough, task lists, autolinks). |
| HTML DOM parser | Parses raw HTML into a traversable document tree for the conversion engine to walk. |
| HTTP client library | Used only by the companion Python script for fetching HTML from URLs. Not a dependency of the core conversion tool. |

## Notational Conventions

### Normative Keywords

This specification uses RFC 2119 / BCP 14 keywords: **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, **MAY**. All other modifiers ("usually", "typically", "recommended") are prohibited in normative statements and MUST be rewritten as one of the five keywords or downgraded to (informative).

### ABNF

Input syntax for CLI arguments is expressed in ABNF (RFC 5234 + RFC 7405). Non-terminals not defined here reference RFC 5234 Appendix B.1 core rules.

### Algorithm Blocks

Key operations are expressed as paper-style pseudocode with numbered lines, `Require:` / `Ensure:` pre/post conditions, assignment via `←`, and comments via `▷`.

### State Transitions

Where state-bearing behavior exists, it is expressed as `(state, event) -> (state', action, output)` tables.
