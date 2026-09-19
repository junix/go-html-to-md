# go-html-to-md

Go 实现的 HTML→Markdown 转换器，提供 stdin/stdout CLI 与 CGO C API；行为契约以 `spec/`、`convert.go` 和测试为准。

- CLI 只负责 I/O，转换规则集中在共享 Go 代码，CLI 与 CGO 不得分叉语义。
- CGO 返回的 C 字符串必须提供并使用对应释放函数；错误通过明确契约返回，不 panic 或泄漏内存。
- `read-url.py` 只是输入获取辅助工具，不属于转换核心，也不能把网络策略带入库。
- 不在说明文件复制依赖版本、生成文件或手工 build 命令；以 `go.mod` 和 `justfile` 为准。

验证：`just build && just test`。
