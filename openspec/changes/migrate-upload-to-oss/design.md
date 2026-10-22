## Context

部署基线已写入 `oss` 配置段与 `ILLNESS_OSS_*` env，但 `POST /upload` 仍写 Mac 路径并挂静态 `/images`。小游戏首版不传图；本变更清服务端债务并兑现「服务端直传 OSS → CDN URL」。约束：单服务、密钥不入库、前缀 `illness/`、同桶同 CDN。

## Goals / Non-Goals

**Goals:**

- `POST /upload` 将文件写入 OSS（`illness/` 前缀），响应返回 CDN 绝对 URL
- 删除本地 `FileDir`、静态 `/images` 挂载
- 机器人等硬编码旧 `/images/...` URL 改为 CDN
- 启动时从 yaml + env 组装可用的 OSS 客户端（本地可用公网 endpoint 覆盖）

**Non-Goals:**

- 不改 wx-game 客户端
- 不做客户端预签名直传（STS / uploadUrl）——首版服务端收 multipart 即可
- 不做图片处理（缩略图、鉴黄）
- 不引入 Redis/MQ/Nginx，不声明本地图片 volume
- 不迁移历史已落在某 Mac 盘上的文件

## Decisions

### D1：保留路径 `POST /upload`，破坏性改响应

- **选择**：路由仍为 `POST /upload`；成功时响应体为 CDN 绝对 URL 字符串（或与现网一致的纯文本写出方式，但内容为完整 URL）。相对路径 `images/...` 废弃。
- **理由**：无小游戏调用方，破坏成本低；保留路径降低「另找路由」噪音。
- **备选**：新路径 `/v2/upload` 并 410 旧路由——多余。

### D2：服务端直传，非预签名

- **选择**：handler 读 multipart（字段名可继续兼容 `upload-file`）→ PutObject → 拼 `cdnBaseUrl` + objectKey 返回。
- **理由**：实现短；与「小游戏暂不传图、偶发服务端/工具上传」匹配。预签名留给未来客户端直传。
- **备选**：对齐 go_ai_talk UCG presign——对本仓当前调用面过重。

### D3：objectKey 与 md5

- **选择**：优先使用请求 header `md5`（若非空）生成 `illness/<md5>.png`（或保留扩展名策略与现逻辑一致）；缺省则服务端生成唯一 key（如 uuid）。同 key 覆盖或跳过均可，design 实现选「已存在则直接返回 CDN URL」以贴近旧「文件已存在」语义。
- **理由**：旧接口已有 md5 header；利于幂等。

### D4：OSS 客户端封装位置

- **选择**：薄封装于 `internal/platform/oss`（或等价），从 `g.Cfg()` 读 bucket/endpoint/prefix/cdn，从 env 注入 AK/SK（及可选 `ILLNESS_OSS_ENDPOINT` 覆盖 endpoint）。`cmd` 的 Upload 只调封装。
- **理由**：与现有 `dbcfg` 手法一致；避免把 SDK 细节堆在 `cmd.go`。

### D5：去掉静态 `/images`；机器人头像改 CDN

- **选择**：删除 `AddStaticPath` / `FileDir` / `Mkdir`。机器人头像改为配置或常量形式的 CDN URL（例如已上传到 `illness/robot_avatar.webp` 的完整 CDN 地址；若对象尚未存在，tasks 中注明需运维上传一次或临时用占位 CDN 路径并文档说明）。
- **理由**：静态目录是本地盘契约的一部分，必须随上传一并移除。

### D6：依赖

- **选择**：引入官方/常用阿里云 OSS Go SDK（与 go_ai_talk 同类即可），版本钉在 go.mod。
- **理由**：自写签名成本高。

## Risks / Trade-offs

- [本地/CI 无 AK 时上传失败] → Mitigation：未配置密钥时 Upload 返回明确错误；健康检查与教学/房间路径不依赖上传
- [内网 endpoint 在开发机不可达] → Mitigation：`.env.local` 设 `ILLNESS_OSS_ENDPOINT` 为公网 endpoint（已在 example 注释）
- [BREAKING 响应格式] → Mitigation：已确认 wx-game 无调用；文档/spec 标明
- [机器人头像对象不存在 404] → Mitigation：tasks 含「确认或上传固定 objectKey」；或先改 URL 并在运维清单列出
- [与未收版 revive 的「不改造上传」条文冲突] → Mitigation：本 change 对 `oss-cdn-illness-prefix` REMOVED 该要求；收版时合并顺序以较新 change 为准

## Migration Plan

1. 合入本变更 → 服务器/本地 env 填好 `ILLNESS_OSS_*`（开发机加公网 endpoint）
2. （可选）将 `robot_avatar.webp` 上传到 `illness/` 前缀并核对 CDN 可访问
3. 验证：`POST /upload` 返回 `https://resorce.cuplay.top/...`；`/` 与房间 API 不受影响
4. 回滚：恢复旧 handler 需回退镜像；OSS 上已写对象可保留

## Open Questions

- （无阻塞）机器人资源是沿用现成 CDN 文件还是本 change apply 时手工上传一次——实现时选「常量指向约定 CDN URL + tasks 运维核对」即可
