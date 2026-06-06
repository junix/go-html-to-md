# Spec Changelog

## 2026-06-04T21:59 — CREATED

Initial specification for go-html-to-md: CLI tool and C FFI shared library for HTML-to-Markdown conversion with GFM support.

- Files: 00, 01, 02, 03, 04, 05
- Code basis: b3ad413
- Small project (< 5 source files); feature matrix omitted (single-purpose converter with 2 entry surfaces, < 15 features)
- Glossary created (5 core nouns identified)
- Most DoD items marked [U]; only conversion correctness for common HTML elements verified via sample HTML
- No automated test suite exists; `just test` is a single stdin smoke test
