## ADDED Requirements

### Requirement: 上传直传 OSS 并返回 CDN URL

系统 MUST 提供 `POST /upload`：接收 multipart 文件（兼容字段名 `upload-file`），将对象写入与配置一致的 OSS 桶，objectKey MUST 落在 `objectKeyPrefix`（`illness/`）之下，成功响应 MUST 为可通过 `cdnBaseUrl` 拼出的 **CDN 绝对 URL**（明文写出，与现 handler 写字符串风格一致即可）。系统 MUST NOT 将上传文件持久化到本地业务目录，MUST NOT 再挂载静态路径 `/images` 提供上传产物。

#### Scenario: 成功上传返回 CDN

- **WHEN** 客户端 `POST /upload` 提交合法图片且 OSS 凭证与配置可用
- **THEN** 对象出现在 `illness/` 前缀下，响应体为以配置的 `cdnBaseUrl` 为根的绝对 URL（非 `images/` 相对路径）

#### Scenario: 幂等 md5

- **WHEN** 请求携带可用的 `md5` header 且对应 objectKey 已存在于 OSS
- **THEN** 系统可不重复上传，直接返回同一 CDN URL

#### Scenario: 缺少 OSS 凭证

- **WHEN** 未配置 `ILLNESS_OSS_ACCESS_KEY_ID` / `ILLNESS_OSS_ACCESS_KEY_SECRET`（或等价注入失败）时调用上传
- **THEN** 请求失败并返回可诊断错误，且 MUST NOT 回退为本地盘保存

### Requirement: 移除本地盘上传与静态 images

系统 MUST 删除基于硬编码本地目录（如历史 `FileDir`）的保存逻辑，以及为上传产物服务的 `AddStaticPath("/images")`（或等价静态挂载）。进程启动 MUST NOT 依赖创建本机图片目录才能提供非上传类 API。

#### Scenario: 无本地图片目录依赖

- **WHEN** 在无历史 Mac/本地图片目录的环境中启动服务（仅配置 DB/OSS env）
- **THEN** HTTP 服务可启动；上传走 OSS，其它已有 API 不因缺少本地图片目录而失败

### Requirement: 内置资源 URL 使用 CDN

系统内硬编码或默认的头像等图片 URL（例如机器人头像）MUST 使用 CDN 绝对地址（`illness/` 前缀下资源或等价公开 CDN URL），MUST NOT 再指向服务主机的 `/images/...` 本地静态路径。

#### Scenario: 机器人头像

- **WHEN** 房间或相关逻辑返回机器人玩家信息中的头像字段
- **THEN** 该字段为 CDN 绝对 URL，而非 `http(s)://<host>:9686/images/...`
