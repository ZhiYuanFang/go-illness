## Context

部署约定要求真实密钥不进库；当前 `.gitignore` 有部分 `.env` 规则，但未显式锁定 `manifest/docker/env/`，且缺少对 `__debug_bin.exe` / 全局 `*.exe` 的覆盖。用户决定：**`manifest/docker/env/` 整目录忽略定死**。

## Goals / Non-Goals

**Goals:**

- `manifest/docker/env/` 下任意文件默认不可被 `git add` 纳入版本库
- 调试二进制与常见构建垃圾被忽略
- 根级 / docker 级 `.env.example` 仍可提交

**Non-Goals:**

- 不在 `env/` 下保留 `.gitkeep`（整目录忽略即空目录可不入库；本地自建即可）
- 不改 `.env.example` 内容

## Decisions

### D1：整目录忽略 `manifest/docker/env/`（已锁定）

```gitignore
manifest/docker/env/
```

本地真实 env 只放此目录；仓库模板仅 `manifest/docker/.env.example`。Compose 文档/注释可继续写 `--env-file manifest/docker/env/.env.prod`（路径本地存在即可）。

### D2：叠加通用产物规则

- `*.exe`、`__debug_bin*`、`**/__debug_bin*`
- `*.log`、`logs/`
- 可选：`coverage.*`、`*.prof`

保留既有 `main` / `vendor/` / IDE 忽略。

### D3：env 模板例外

保留：

```gitignore
.env
.env.*
!.env.example
manifest/docker/.env
manifest/docker/.env.*
!manifest/docker/.env.example
manifest/docker/env/
```

整目录规则在后，确保 `env/` 内即使名为 `.env.example` 亦不提交（避免误把密钥当 example）。

### D4：已跟踪文件清理

若 `__debug_bin.exe` 等已在 index：`git rm --cached`；工作区文件保留。

## Risks / Trade-offs

- [本地忘记建 env 目录] → Mitigation：`.env.example` 注释写明复制到 `manifest/docker/env/`
- [曾暴露密钥] → Mitigation：tasks 提醒轮换；gitignore 只防未来

## Open Questions

- （无）
