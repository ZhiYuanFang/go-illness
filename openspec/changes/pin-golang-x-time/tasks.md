## 1. 钉版本

- [x] 1.1 `go get golang.org/x/time@v0.6.0` 并 `go mod tidy`，确认 `go.mod`/`go.sum` 为 v0.6.0
- [x] 1.2 若 tidy/build 仍因其它 `golang.org/x/*` 报过高 Go 版本，按同样原则钉到兼容版本（最小范围）

## 2. 验收

- [x] 2.1 本机 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o NUL .`（或等价）通过
- [x] 2.2 `openspec validate pin-golang-x-time --strict`
