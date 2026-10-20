## 1. 更新忽略规则

- [x] 1.1 更新根 `.gitignore`：定死 `manifest/docker/env/`；补强 `*.exe`、`__debug_bin*`、日志等；保留 `.env.example` 例外
- [x] 1.2 （可选）在 `manifest/docker/.env.example` 顶部注释写明：真实文件复制到 `manifest/docker/env/`，该目录不入库

## 2. 清理索引

- [x] 2.1 对已 staged 的应忽略文件执行 `git rm --cached`（至少 `internal/cmd/__debug_bin.exe`；若有 env 误加一并处理）
- [x] 2.2 用 `git check-ignore -v` 验证 `manifest/docker/env/.env.prod`（或占位路径）与 `__debug_bin.exe` 被忽略；确认 `.env.example` 不被忽略

## 3. 校验

- [x] 3.1 `openspec validate harden-gitignore-secrets --strict`
