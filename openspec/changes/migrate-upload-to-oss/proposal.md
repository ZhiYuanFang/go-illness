## Why

`revive-illness-deploy-baseline` 已落地 OSS/CDN 配置约定（同桶同 CDN、前缀 `illness/`），但明确推迟了上传改造；服务端仍使用 Mac 硬编码 `FileDir` 落盘与静态 `/images`，容器与 Windows 本地均不可用。微信小游戏首版**不调用**上传，但服务端需兑现「直传 OSS、返回 CDN URL」，并清掉本地盘债务（含机器人头像硬编码旧机路径）。

**对照说明**：已读 `openspec/project.md` 与空壳基线 `v0.0.1`。`oss-cdn-illness-prefix` 尚在活跃 change `revive-illness-deploy-baseline` 中（未收版进基线）；本变更复用其配置约定，并**解除**其中「本阶段不改造上传接口」冻结。

## What Changes

- **BREAKING**：`POST /upload` 成功响应由相对路径 `images/<md5>.png` 改为 **CDN 绝对 URL**；不再落本地盘、不再挂载静态 `/images`
- 服务端使用已配置的 OSS 客户端（密钥经 `ILLNESS_OSS_*`）将文件写入 `illness/` 前缀 objectKey，返回 `cdnBaseUrl` + key
- 删除 `FileDir` / `AddStaticPath("/images")` 等本地上传与静态服务逻辑
- 机器人等硬编码的 `.../images/robot_avatar.webp` 改为 CDN URL（或固定 `illness/` 下公开资源 URL）
- 不改微信小游戏客户端（已确认无 `/upload` 调用）
- 不引入 Redis/MQ/Nginx；不新增本地图片 volume

## Capabilities

### New Capabilities

- `illness-upload-oss`：上传 HTTP 契约——服务端直传 OSS、返回 CDN URL；去掉本地盘与静态 `/images`

### Modified Capabilities

- `oss-cdn-illness-prefix`：移除「本阶段不改造上传接口」要求；配置基线仍有效，供上传实现复用（该 capability 目前仅存在于未收版的 `revive-illness-deploy-baseline` delta，本变更以 MODIFIED/REMOVED 对齐）

## Impact

- 代码：`internal/cmd/cmd.go`（Upload、静态路径）、可能新增薄 OSS 客户端/`internal/platform` 封装；`next_illness` 等硬编码头像 URL
- 配置：复用 `manifest/config` 的 `oss` 段与 Compose/env 中 `ILLNESS_OSS_*`；本地非 VPC 须配置公网 `ILLNESS_OSS_ENDPOINT`
- API：`POST /upload` 响应格式破坏性变更；请求侧可保留 multipart + `md5` header（若实现简化可调整，须在 design/specs 写明）
- 客户端：wx-game 首版无依赖；任何仍期望相对路径的旧调用方需改吃 CDN URL
- 依赖：阿里云 OSS Go SDK（或与 go_ai_talk 同类封装）；新增依赖须最小
- 收版后：`project.md` 中「上传另开 change」条文应变更为「上传已走 OSS」类表述（可在本 change tasks 或 archive 时改）
