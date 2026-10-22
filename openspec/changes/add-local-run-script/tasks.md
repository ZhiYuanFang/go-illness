## 1. 本地启动脚本

- [x] 1.1 新增 `hack/run-local.ps1`：默认加载 `manifest/docker/env/.env.prod`，支持 `-EnvFile`；解析 `#`/空行/首个 `=`；未设置时默认补 `ILLNESS_OSS_ENDPOINT`（`-SkipOssEndpoint` 可关）
- [x] 1.2 脚本切到仓库根后执行 `gf run main.go`；env 文件不存在时给出明确错误退出

## 2. README

- [x] 2.1 更新根 `README.MD`：本地运行步骤（env 准备、默认 `.env.prod`、`.\hack\run-local.ps1`、`:9686` 验收）；可附极简其它命令提示；不含真实密钥
- [x] 2.2 在 README 中醒目提示：默认连接 `.env.prod` 所指向的库，确认环境后再启动

## 3. 校验

- [x] 3.1 `openspec validate add-local-run-script --strict`
