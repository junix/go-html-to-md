set shell := ["bash", "-uo", "pipefail", "-c"]

arch_suffix := if arch() == "aarch64" { "arm64" } else { "x86" }
install_bin := home_directory() / "sync" / ("bin_" + arch_suffix)

default: build

build:
    go build -o html-to-markdown html-to-markdown.go

test:
    echo "html-to-markdown" | go run html-to-markdown.go | head -1

install: build
    mkdir -p "{{install_bin}}"
    cp html-to-markdown "{{install_bin}}/html-to-markdown"
    echo "Installed html-to-markdown to {{install_bin}}"
