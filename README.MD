# go-illness（病变服务）

GoFrame 单服务：HTTP + WebSocket，默认端口 `:9686`。业务库名 `illness`；上传走阿里云 OSS（前缀 `illness/`），返回 CDN URL。

## 本地运行

> **注意**：启动脚本**默认**加载 `manifest/docker/env/.env.prod`，会连该文件里配置的 MySQL / OSS。请确认目标环境后再启动，避免误连生产。

### 1. 准备环境文件

```powershell
# 若尚无 env 目录下的真实文件：
Copy-Item manifest\docker\.env.example manifest\docker\env\.env.prod
# 编辑 .env.prod：至少填写 ILLNESS_DB_LINK、MYSQL_TCP_HOST；上传还需 ILLNESS_OSS_*
```

- 模板（可入库）：`manifest/docker/.env.example`
- 真实密钥目录（已 gitignore）：`manifest/docker/env/`
- 本机非 VPC 访问 OSS 时，脚本会在未设置时自动补 `ILLNESS_OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com`

### 2. 一键启动

需已安装 [GoFrame CLI](https://goframe.org)（`gf`）与 Go。

```powershell
cd d:\work\go-illness

# 若执行策略限制脚本，可先：
# Set-ExecutionPolicy -Scope Process Bypass

.\hack\run-local.ps1
```

指定其它 env：

```powershell
.\hack\run-local.ps1 -EnvFile manifest\docker\env\.env.local
```

跳过自动补 OSS 公网 endpoint：

```powershell
.\hack\run-local.ps1 -SkipOssEndpoint
```

### 3. 验收

- 浏览器：`http://127.0.0.1:9686/` → 应显示「病变」
- OpenAPI：`http://127.0.0.1:9686/api.json`

微信小游戏联调：将客户端 `API_BASE` / `WS_BASE` 指向 `127.0.0.1:9686`（见 wx-game 仓 `js/illness/config.js`）。

### 4. 可选：教学迁移

若需教学完成标记字段，在目标库执行：

`hack/migrations/add_user_tutorial_completed.sql`

## 其它常用命令

- `make cli`：安装/更新 GoFrame CLI
- `make dao` / `make service`：按 `hack` 配置生成 DAO / service 接口

## 镜像构建

`manifest/docker/Dockerfile` 基础镜像为官方名（`golang:1.19-alpine` / `alpine:3.19`）。国内拉取慢时请在 Docker Desktop / daemon 配置 `registry-mirrors`，不要在 Dockerfile 里写死第三方代理前缀。`GOPROXY` 仅加速 `go mod`，与 registry 无关。

部署与 Compose 见 `manifest/docker/` 与 OpenSpec `openspec/project.md`。
