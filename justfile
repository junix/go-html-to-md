set shell := ["bash", "-uo", "pipefail", "-c"]

install_bin := home_directory() / "sync/bin"
arch_suffix := if arch() == "aarch64" { "arm64" } else { "x86" }

default: build

build:
  go build -o html-to-markdown html-to-markdown.go

test:
  echo "html-to-markdown" | go run html-to-markdown.go | head -1

install:
  #!/usr/bin/env bash
  set -euo pipefail
  mkdir -p "{{install_bin}}"
  go build -o "{{install_bin}}/html-to-markdown.{{arch_suffix}}" html-to-markdown.go
  echo "Installed html-to-markdown.{{arch_suffix}} to {{install_bin}}"
