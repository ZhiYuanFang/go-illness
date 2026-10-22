## REMOVED Requirements

### Requirement: 本阶段不改造上传接口

**Reason**：已被本变更 `illness-upload-oss` 取代；部署基线阶段的冻结已到期，上传改为服务端直传 OSS 并返回 CDN URL。

**Migration**：实现与验收以 `illness-upload-oss` 为准；配置基线（同桶同 CDN、前缀 `illness/`、密钥走 env）仍由本 capability 的「OSS 与 CDN 配置基线」要求约束（见未收版的 `revive-illness-deploy-baseline` 或收版后基线）。

## ADDED Requirements

### Requirement: 上传实现复用 OSS 配置基线

上传实现 MUST 复用仓库已约定的 OSS 配置（bucket、region、endpoint、`objectKeyPrefix`、`cdnBaseUrl`）及环境变量注入的 AccessKey；MUST NOT 在代码或 yaml 中硬编码 AccessKey。可选的 `ILLNESS_OSS_ENDPOINT` MUST 能覆盖默认 endpoint，以便非 VPC 环境使用公网 endpoint。

#### Scenario: 配置与密钥分离

- **WHEN** 审查上传相关实现与默认配置
- **THEN** objectKey 使用前缀 `illness/`，CDN 根与同桶约定一致，密钥仅来自环境变量

#### Scenario: 本地公网 endpoint 覆盖

- **WHEN** 设置 `ILLNESS_OSS_ENDPOINT` 为公网 endpoint 并启动服务后执行上传
- **THEN** 客户端连接使用该覆盖值，而非仅内网 endpoint
