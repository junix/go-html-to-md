set shell := ["bash", "-euo", "pipefail", "-c"]

os_suffix := if os() == "macos" { "macos" } else { "linux" }
arch_suffix := if arch() == "aarch64" { "arm64" } else { "x86" }
install_bin := env("SYNC_BIN_DIR", home_directory() / "sync" / (os_suffix + "-" + arch_suffix + "-bin"))

default: build

# 构建（Release 模式）
build:
    go build -o html-to-markdown html-to-markdown.go convert.go

# 运行测试
test:
    go test convert.go convert_test.go
    echo "html-to-markdown" | go run html-to-markdown.go convert.go | head -1

# 安装到 ~/sync/<os>-<arch>-bin/
install: build
    mkdir -p "{{install_bin}}"
    cp html-to-markdown "{{install_bin}}/html-to-markdown"
    echo "Installed html-to-markdown to {{install_bin}}"
