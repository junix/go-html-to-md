# Go-Html-To-Md (`jtools/go-html-to-md/`)

## 概述

HTML 到 Markdown 转换工具，提供 CLI 和 CGO 共享库两种使用方式。

## 核心文件

```
go-html-to-md/
├── html-to-markdown.go       # CLI 版本（stdin/stdout）
├── html-to-markdown-lib.go   # CGO 库版本（导出 C 接口）
└── read-url.py               # Python URL 读取工具
```

## 主要功能

### CLI 工具（html-to-markdown.go）
从 stdin 读取 HTML，输出 Markdown 到 stdout。

```bash
cat example.html | ./html-to-markdown > output.md
```

### CGO 库（html-to-markdown-lib.go）
导出 C 函数供其他语言调用：

```c
char* ConvertHTMLToMarkdown(const char* html);
char* ConvertHTMLFileToMarkdown(const char* filepath);
```

CLI 用法：
```bash
./html-to-markdown-lib input.html output.md
```

### Python 工具（read-url.py）
读取 URL 内容并输出 HTML：

```bash
python read-url.py -u "https://example.com" > input.html
```

## 构建

### CLI 版本
```bash
go build -o html-to-markdown html-to-markdown.go
```

### CGO 共享库
```bash
go build -o html-to-markdown.so -buildmode=c-shared html-to-markdown-lib.go
```

## 依赖

- **Go 1.19+**
- `github.com/tomkosm/html-to-markdown` — HTML 转 Markdown 核心库
- `github.com/PuerkitoBio/goquery` — HTML 解析
- `golang.org/x/net` — 网络支持
- **Python**: `requests`（仅 read-url.py 使用）

## 注意事项

- CGO 版本导出的函数返回的 C 字符串需要调用方释放内存
- 使用 GitHub Flavored Markdown 插件进行转换
- 错误处理：CGO 版本在出错时返回错误字符串，而非崩溃
