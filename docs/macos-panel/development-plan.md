# Oasis macOS 开发实施计划

| 属性 | 内容 |
|------|------|
| 文档版本 | v0.1（草案） |
| 对应 PRD | [`prd.md`](./prd.md) v0.2 |
| 产品名称 | Oasis |
| 默认数据目录 | `~/Library/Application Support/Oasis` |
| 默认访问范围 | 允许局域网访问（LAN） |
| 首批暂缓 | 本机 shell 终端、商店更新源最终形态 |

---

## 1. 开发目标

首批开发目标是把 1Panel 2.x 架构改造成可在 macOS 上原生运行的 Oasis 基线版本：

1. `core`、`agent` 可以编译为 darwin arm64/amd64 原生二进制。
2. 使用 launchd 托管 Oasis 自身进程。
3. 默认数据、日志、运行时 socket 均落在 `~/Library/Application Support/Oasis` 下。
4. 面板默认允许局域网访问，并在首次启动和设置页展示安全提示。
5. Docker 引擎可自动探测 OrbStack / Docker Desktop，并支持手动配置 socket。
6. 容器、镜像、网络、卷、Compose 的常用管理能力可用。
7. 前端隐藏 Linux-only 模块，保留 macOS v1 范围内的功能。
8. 建立 Mac 应用商店的最小骨架，先跑通少量精选应用，再扩展到 PRD 要求的 15 个应用。

---

## 2. 分支与 PR 策略

macOS 产品线建议长期使用独立主干：

```text
platform/macos
  ├── feature/platform-foundation
  ├── feature/docker-engine
  ├── feature/frontend-platform
  ├── feature/mac-store
  └── release/macos-v1.0.x
```

在当前仓库执行时，使用 `cursor/` 前缀创建短生命周期分支。每个 PR 只承载一个可独立评审的主题，避免把平台抽象、Docker、前端裁剪和商店一次性混在一起。

推荐 PR 顺序：

| 顺序 | PR 主题 | 主要目的 | 合并前验收 |
|------|---------|----------|------------|
| 1 | 平台路径与配置地基 | 引入 Oasis 路径约定，去掉关键 Linux 路径硬编码 | Linux 现有构建不破坏；darwin 可编译到下一个阻塞点 |
| 2 | darwin 构建与运行时 socket | 新增 darwin 构建目标，迁移 agent socket 到 Oasis 运行目录 | core/agent 可生成 darwin 二进制；socket 路径可配置 |
| 3 | launchd 生命周期 | 用 launchd 托管 core/agent，并提供 `oasisctl` | `oasisctl status/restart/stop` 可用；卸载不残留服务 |
| 4 | Docker 引擎探测 | 支持 OrbStack / Docker Desktop socket 探测与断连提示 | 无 Docker 时面板不崩溃；有 Docker 时 `docker info` 成功 |
| 5 | Docker 管理闭环 | 验证容器、镜像、网络、卷、Compose 常用操作 | 容器列表、启停、日志、Compose up/down 可用 |
| 6 | 前端平台开关 | 隐藏 Linux-only 菜单和路由，加入 Oasis 文案 | macOS 构建只显示 v1 范围内功能 |
| 7 | Mac 商店骨架 | 新增 `resource/apps/mac` 规范和后端读取流程 | 至少一个应用可通过 Compose 安装和卸载 |
| 8 | 发布与诊断 | tar.gz 发布包、诊断包、升级/快照命名 | 发布物命名正确；诊断包可导出关键状态 |

---

## 3. 阶段拆分

### P0：平台地基

目标：让 Oasis 在 macOS 上具备正确的目录、配置、构建和进程管理模型。

#### 交付物

- 平台抽象接口与 darwin/linux 实现。
- Oasis 默认目录结构：

```text
~/Library/Application Support/Oasis/
  conf/
  db/
  log/
  run/
  apps/
  resource/apps/mac/
  backup/
  tmp/
  bin/
```

- `oasis-core`、`oasis-agent` darwin 构建目标。
- `oasisctl` CLI 命名与基础生命周期命令。
- launchd plist：

```text
~/Library/LaunchAgents/dev.oasis.core.plist
~/Library/LaunchAgents/dev.oasis.agent.plist
```

#### 关键代码入口

| 主题 | 路径 |
|------|------|
| 配置初始化 | `core/init/viper/viper.go`、`agent/init/viper/viper.go` |
| 目录初始化 | `agent/init/dir/dir.go` |
| core 到 agent 代理 | `core/init/proxy/proxy.go` |
| agent 监听 | `agent/server/server.go` |
| CLI 配置 | `core/utils/ctl_conf/ctl_conf.go`、`core/cmd/server/cmd/` |
| 服务管理 | `core/utils/controller/`、`agent/utils/controller/` |
| 构建 | `Makefile` |

#### 验收标准

- macOS 上能生成 `oasis-core` 和 `oasis-agent`。
- 首次启动自动创建 Oasis 数据目录。
- agent socket 不再依赖 `/etc/1panel`。
- `oasisctl status` 能看到 core/agent 状态。
- Linux 默认路径和现有启动方式不被破坏。

---

### P1：Docker 引擎和容器管理

目标：让 Oasis 能稳定连接本机 Docker 引擎，并完成 Docker 日常管理闭环。

#### 交付物

- Docker socket 探测顺序：

```text
1. 用户配置 DockerSockPath
2. DOCKER_HOST
3. unix:///var/run/docker.sock
4. unix://$HOME/.docker/run/docker.sock
5. docker context 解析（增强项）
```

- 引擎识别：
  - OrbStack
  - Docker Desktop
  - 通用 Docker
- Docker 未运行时的后端错误模型和前端引导。
- macOS 下禁用 systemd 管理 Docker 的操作。
- 容器、镜像、网络、卷、Compose 常用操作冒烟通过。

#### 关键代码入口

| 主题 | 路径 |
|------|------|
| Docker 客户端 | `agent/utils/docker/docker.go` |
| Docker 服务状态 | `agent/app/service/docker.go` |
| 容器管理 | `agent/app/service/container.go` |
| Compose 命令 | `agent/utils/compose/compose.go` |
| Compose 业务 | `agent/app/service/container_compose.go` |

#### 验收标准

- Docker Desktop 和 OrbStack 均可自动识别或手动配置成功。
- Docker 断连时 API 返回明确错误，前端不白屏。
- 容器列表、详情、启停、删除、日志可用。
- 镜像拉取、删除、清理可用。
- Compose 项目列表、启动、停止、拉取镜像、编辑 compose 文件可用。

---

### P2：前端平台裁剪与 Oasis 品牌

目标：让用户看到的是 Oasis 的 macOS 功能面，而不是裁剪不完整的 Linux 面板。

#### 交付物

- 新增平台配置：

```text
frontend/src/config/platform.ts
```

- 支持构建变量：

```text
VITE_PLATFORM=darwin
```

- 菜单和路由裁剪：
  - 保留：首页、容器、应用商店（Mac）、文件、终端、计划任务精简版、设置、日志。
  - 隐藏：网站、主机防火墙、SSH、FTP、fail2ban、ClamAV、宿主机数据库、宿主机 Supervisor。
- 关于页、标题栏、安装页统一使用 Oasis。
- 首页增加 Docker 引擎状态卡片。
- 默认允许 LAN 访问的安全提示。

#### 关键代码入口

| 主题 | 路径 |
|------|------|
| 路由 | `frontend/src/routers/router.ts`、`frontend/src/routers/modules/` |
| 菜单初始化 | `core/init/migration/helper/menu.go` |
| 菜单渲染 | `frontend/src/layout/components/Sidebar/` |
| 首页 | `frontend/src/views/home/` |
| 设置页 | `frontend/src/views/setting/` |

#### 验收标准

- macOS 构建中不出现 Linux-only 菜单入口。
- 直接访问隐藏路由时返回无权限、404 或安全降级页面。
- Oasis 命名在主界面、关于页、CLI 提示中一致。
- 首次启动或设置页明确说明当前允许局域网访问。

---

### P3：Mac 应用商店最小闭环

目标：建立 Oasis 自研 Mac 商店目录规范，并跑通应用安装、卸载、升级的最小链路。

#### 交付物

- 新增目录：

```text
resource/apps/mac/
  index.yaml
  <app-key>/
    manifest.yaml
    logo.png
    docker-compose.yml
    .env.example
```

- 新增或扩展 `resource` 类型：`mac`。
- 安装前检查：
  - 架构
  - Docker 可达性
  - 端口冲突
  - 安装目录文件共享提示
  - 内存和磁盘提示
- 先跑通 3 个试点应用：
  - Home Assistant
  - Open WebUI
  - PostgreSQL 或 Redis

#### 关键代码入口

| 主题 | 路径 |
|------|------|
| 应用列表和同步 | `agent/app/service/app.go` |
| 应用安装 | `agent/app/service/app_install.go` |
| 应用模型 | `agent/app/model/`、`agent/app/dto/` |
| 目录初始化 | `agent/init/dir/dir.go` |
| 前端应用商店 | `frontend/src/views/app-store/` |

#### 验收标准

- Mac 商店索引可被后端读取。
- 至少一个应用能从前端完成参数填写、安装、启动、卸载。
- 安装失败时能看到明确日志和回滚结果。
- 安装记录复用现有 `AppInstall` 能力，不引入重复模型。

---

### P4：发布、升级和诊断

目标：让 Oasis 具备外部试用版需要的发布物、升级约束和排障能力。

#### 交付物

- 发布物命名：

```text
oasis-macos-{version}-arm64.pkg
oasis-macos-{version}-amd64.pkg
oasis-macos-{version}.tar.gz
SHA256SUMS
```

- 快照命名：

```text
oasis-{scope}-{version}-darwin-{arch}-{time}
```

- 禁止跨 OS 恢复。
- 诊断包导出：
  - Oasis 版本
  - macOS 版本
  - CPU/内存/磁盘摘要
  - Docker info
  - docker compose 项目列表
  - 最近 core/agent 错误日志

#### 关键代码入口

| 主题 | 路径 |
|------|------|
| 升级 | `core/app/service/upgrade.go` |
| 快照 | `agent/app/service/snapshot_create.go` |
| 备份 | `agent/app/service/backup*.go`、`agent/app/service/cronjob*.go` |
| 日志 | `core/init/log/`、`agent/init/log/` |
| 发布脚本 | `Makefile`、`scripts/macos/` |

#### 验收标准

- 发布物命名不再使用 Linux 包格式。
- 不允许 Linux 快照恢复到 Oasis，也不允许 Oasis 快照恢复到 Linux 版。
- 诊断包不包含敏感 Token、密码或私钥。
- 用户能用诊断包定位 Docker 不可达、端口冲突、launchd 失败等常见问题。

---

## 4. 首批开发检查清单

### 开始编码前

- [ ] 确认 Go 模块是否允许仓库根部新增共享包；如不允许，先在 core/agent 内各自放薄封装。
- [ ] 梳理所有 `/opt/1panel`、`/etc/1panel`、`/usr/local/bin/1pctl` 的关键路径使用点。
- [ ] 确认 core 和 agent 的配置读取顺序，避免迁移路径后首次启动找不到配置。
- [ ] 确认现有 Linux CI 或构建命令，平台抽象不能破坏上游 Linux 行为。

### 每个 PR 的固定验收

- [ ] `git diff --check` 通过。
- [ ] 修改 Go 代码时运行对应模块的格式化和最小测试。
- [ ] 修改前端时运行类型检查或项目已有的 lint/build 命令。
- [ ] 修改文档时同步 README 或 PRD 链接。
- [ ] 涉及平台差异时说明 linux/darwin 行为差异。

---

## 5. 技术风险和处理原则

| 风险 | 处理原则 |
|------|----------|
| Linux 逻辑散落 | 先抽路径、服务、Docker 三个最小接口，不一开始做大而全平台层 |
| Go module 边界 | 若根部 `internal/platform` 无法被 core/agent 同时引用，则改用 `pkg/platform` 或双端薄封装 |
| Docker Desktop 和 OrbStack 差异 | 所有自动探测都允许用户在设置页手动覆盖 socket |
| launchd 调试复杂 | 安装脚本必须提供 `oasisctl logs` 和诊断命令 |
| 默认 LAN 访问风险 | 首次启动强制改密，并在设置页持续展示当前绑定地址 |
| 文件共享权限 | 应用安装 preflight 提前检查并给出 Docker Desktop/OrbStack 指引 |
| 上游 drift | Mac 分支只 cherry-pick 安全修复和通用 bugfix，避免整分支 merge |

---

## 6. 近期建议从哪里开始

第一批实际开发建议从 **P0 平台地基** 开始，按以下顺序落地：

1. 新增 `Paths` 抽象，只覆盖安装目录、配置目录、日志目录、运行目录、应用目录。
2. 将 core/agent 中启动必需的硬编码路径切到 `Paths`。
3. 把 agent socket 移到 `{InstallDir}/run/agent.sock`。
4. 增加 darwin 构建目标，先让编译错误显性化。
5. 新增 launchd plist 模板和 `oasisctl status`。

这组工作完成后，Oasis 的后续 Docker、前端裁剪和 Mac 商店都可以在稳定的 macOS 运行时基础上继续推进。
