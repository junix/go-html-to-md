# 03 Runtime Model

## Execution Semantics

### Algorithm 1 StdinPipeConversion

```
Algorithm 1  StdinPipeConversion
Require: stdin is open for reading
Ensure:  stdout receives the Markdown output of converting the HTML input,
         or the process exits with non-zero status
 1: html ← ReadAll(stdin)                          ▷ read until EOF
 2: if read error occurs then
 3:   write diagnostic to stderr
 4:   exit(1)
 5: end if
 6: converter ← CreateConverter()
 7: md, err ← converter.Convert(html)
 8: if err ≠ nil then
 9:   write diagnostic to stderr
10:   exit(1)
11: end if
12: write(md) to stdout                            ▷ no trailing newline appended beyond converter output
```

### Algorithm 2 FileArgumentConversion

```
Algorithm 2  FileArgumentConversion
Require: input_file path is provided as first positional argument
         output_file path MAY be provided as second positional argument
Ensure:  Markdown output is written to output_file if provided,
         or to stdout otherwise;
         process exits with non-zero status on any I/O or conversion error
 1: content ← ReadFile(input_file)
 2: if read error occurs then
 3:   write diagnostic to stderr
 4:   exit(1)
 5: end if
 6: converter ← CreateConverter()
 7: md, err ← converter.Convert(content)
 8: if err ≠ nil then
 9:   write diagnostic to stderr
10:   exit(1)
11: end if
12: if output_file is provided then
13:   WriteFile(output_file, md)
14:   if write error occurs then
15:     write diagnostic to stderr
16:     exit(1)
17:   end if
18: else
19:   write(md) to stdout
20: end if
```

### Algorithm 3 CFFIStringConversion

```
Algorithm 3  CFFIStringConversion
Require: html is a null-terminated C string
Ensure:  returns a null-terminated C string containing Markdown output,
         or a null-terminated C string prefixed with "Error:" on conversion failure
 1: html_go ← Unmarshal(html) from C string to native string
 2: converter ← CreateConverter()
 3: md, err ← converter.Convert(html_go)
 4: if err ≠ nil then
 5:   return Marshal("Error: Failed to convert HTML to Markdown") to C string
 6: end if
 7: return Marshal(md) to C string
```

### Algorithm 4 CFFIFileConversion

```
Algorithm 4  CFFIFileConversion
Require: filepath is a null-terminated C string
Ensure:  returns a null-terminated C string containing Markdown output,
         or a null-terminated C string prefixed with "Error:" on failure
 1: content, err ← ReadFile(filepath)
 2: if err ≠ nil then
 3:   return Marshal("Error: Failed to read file") to C string
 4: end if
 5: return CFFIStringConversion(content)            ▷ delegates to Algorithm 3
```

## Cross-Operation Invariants

1. **Determinism**: For the same HTML input, all four algorithms MUST produce identical Markdown output (same converter configuration).
2. **Monotonicity**: `length(md_output) ≥ 0` for any input. Empty or whitespace-only HTML input produces output determined by the conversion engine (MAY be empty).
3. **Error containment (FFI)**: Algorithms 3 and 4 MUST NOT cause the host process to abort. Errors are returned as human-readable C strings.

## CLI Entry Dispatch

When invoked without file arguments and with piped stdin, Algorithm 1 executes. When invoked with at least one file argument, Algorithm 2 executes. The help flag (`-h` or `--help`) is checked before any dispatch and prints usage text to stdout, then exits with status 0.

```
(event,                  condition)               → (action,                     exit code)
(no args + piped stdin,  true)                    → (Algorithm 1,                0 on success)
(no args + no pipe,      true)                    → (print usage to stdout,      0)
(-h or --help,           first arg matches)       → (print help to stdout,       0)
(input_file,             first arg is a path)     → (Algorithm 2,                0 on success)
(input_file + out_file,  two args provided)       → (Algorithm 2 with out_file,  0 on success)
(any,                    I/O or conversion error) → (diagnostic to stderr,       1)
```

## Validation and Diagnostics

### Error Model

| Trigger | Severity | Target | Human-Readable Message | Exit / Return |
|---------|----------|--------|------------------------|----------------|
| stdin read failure | error | stderr | "Error reading from stdin: \<detail\>" | exit(1) |
| file read failure (CLI) | error | stderr | "Error: error reading input file: \<detail\>" | exit(1) |
| file write failure (CLI) | error | stderr | "Error: error writing to output file: \<detail\>" | exit(1) |
| conversion failure (CLI) | error | stderr | "Error converting HTML to Markdown: \<detail\>" or "Error converting HTML: \<detail\>" | exit(1) |
| conversion failure (FFI) | error | return value | "Error: Failed to convert HTML to Markdown" | returned as C string |
| file read failure (FFI) | error | return value | "Error: Failed to read file" | returned as C string |

No machine-parseable error codes are defined. The `\<detail\>` portion contains the underlying error message from the conversion engine or operating system.

### Diagnostic Reachability

All diagnostics listed above are reachable from their respective public entry points. The FFI error messages are only reachable via the C FFI surface, not from the CLI. CLI error messages are only reachable from CLI invocation, not from the FFI.
