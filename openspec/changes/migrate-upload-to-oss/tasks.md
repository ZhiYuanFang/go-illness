## 1. OSS 客户端与配置注入

- [x] 1.1 增加阿里云 OSS Go 依赖（`go.mod`），实现 `internal/platform/oss`（或等价）薄封装：读 yaml `oss` 段 + env AK/SK，支持 `ILLNESS_OSS_ENDPOINT` 覆盖
- [x] 1.2 提供 PutObject（或 UploadBytes）与「按 prefix + key 拼 CDN URL」辅助方法；缺密钥时明确报错

## 2. 上传接口改造

- [x] 2.1 重写 `POST /upload`：multipart 字段 `upload-file`，可选 header `md5` → objectKey `illness/<md5>.png`（或 design 约定）；上传 OSS；响应写 CDN 绝对 URL
- [x] 2.2 实现 md5 已存在则直接返回 CDN URL 的幂等路径
- [x] 2.3 删除 `FileDir`、`Mkdir`、`AddStaticPath("/images")` 及本地 `files.Save` 逻辑

## 3. 内置资源 URL

- [x] 3.1 将机器人等硬编码 `.../images/robot_avatar.webp` 改为 CDN 绝对 URL（约定 `illness/` 下资源）
- [x] 3.2 在 tasks/设计备注或 `.env.example` 简短注释中说明：若该对象尚未上传，需运维上传一次

## 4. 文档与约定收尾

- [x] 4.1 更新 `openspec/project.md`：将「上传 HTTP 接口改造另开 change」改为已由本能力覆盖的表述（上传走 OSS、返回 CDN）
- [x] 4.2 `openspec validate migrate-upload-to-oss --strict`
