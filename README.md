构建 go-html-to-md 库，请运行以下命令：

```bash
cd apps/api/src/lib/go-html-to-md
go build -o html-to-markdown.so -buildmode=c-shared html-to-markdown.go
go build -o html-to-markdown html-to-markdown.go
chmod +x html-to-markdown.so
```