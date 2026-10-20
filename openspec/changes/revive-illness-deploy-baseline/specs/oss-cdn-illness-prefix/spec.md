## ADDED Requirements

### Requirement: OSS 与 CDN 配置基线

仓库 MUST 提供与兄弟仓 `go_ai_talk` **相同的 OSS 桶与 CDN 根地址**约定，并为 illness 指定 objectKey 前缀 **`illness/`**。AccessKey 等密钥 MUST 经环境变量注入；配置文件中对应字段 MUST 留空或仅含非密钥的公开项（bucket、region、endpoint、cdnBaseUrl、prefix）。

#### Scenario: 前缀与同源 CDN

- **WHEN** 查阅仓库 OSS 配置约定或 `.env.example`
- **THEN** 可见 `objectKeyPrefix`（或等价项）为 `illness/`，且 bucket / `cdnBaseUrl` 与 go_ai_talk 当前 UCG OSS 约定一致（同桶同 CDN）

#### Scenario: 密钥不入库

- **WHEN** 审查默认配置文件中的 OSS 段
- **THEN** accessKeyId / accessKeySecret 为空，真实值仅出现在未提交的环境文件或密钥管理系统中

### Requirement: 本阶段不改造上传接口

在本 capability 落地范围内，系统 MUST NOT 修改现有上传 HTTP 接口的路径、请求字段与响应格式。OSS/CDN 约定供后续变更实现「服务端直传 OSS 并返回 CDN URL」时复用。

#### Scenario: 上传路由保持原状

- **WHEN** 完成本变更的部署与配置落地
- **THEN** 既有上传接口行为与变更前一致（仍可为本地盘实现）；不因本变更强制切换为 OSS 上传
