# Phase 0 技术验证报告

验证日期：2026-11-08（Asia/Shanghai）。依据：根目录 `SPEC.md` 2.0.0。Wails 版本变更已获用户确认。本阶段仅建立技术探针和最小工程，尚未开发业务功能，也未开始 Phase 1。

## 1. 仓库与初始化

- Git 已初始化，分支 `main`，远程 `https://github.com/Pt-ll/OI-Desk.git`。
- 最小工程包含 `main.go`、`app.go`、`go.mod`、`go.sum`、`wails.json`、`frontend/` 和 `build/`。
- Phase 0 工程、探针和文档已提交到 `main`；当前提交已通过远程三平台 CI。
- `SPEC.md` 保持原样。该文件是实际规范入口；规范内提及的 `PROJECT_SPEC.md` 与 `AGENTS.md` 当前不存在。

## 2. 已锁定版本

| 项目 | 版本 / 方式 | 证据 |
|---|---|---|
| Wails | **v2.17.0**，仅 v2 API | CLI、`go.mod`、官方 tag、完整桌面构建 |
| Go | **1.26.2** | `go.mod`、本机测试、CI 版本文件 |
| Node | **22.23.4 LTS** | `.node-version`；临时下载后实际构建和浏览器测试通过 |
| pnpm | **11.26.0** | `frontend/package.json` 的 packageManager、lockfile |
| Vue | **3.5.45** | package.json 与 pnpm-lock.yaml |
| TypeScript | **5.9.3** | 类型检查通过 |
| Vite | **7.3.7** | 生产构建通过 |
| @vitejs/plugin-vue | **6.0.10** | 锁文件及实际生产构建 |
| vue-tsc | **2.2.12** | 类型检查通过 |
| Playwright | **1.65.0** | 固定开发依赖、Chromium 自动化测试 |
| SQLite Go 驱动 | **modernc.org/sqlite v1.62.0** | 真实 Go 驱动事务/外键/备份测试 |
| 密钥环库 | **github.com/zalando/go-keyring v0.2.8** | macOS Keychain 临时凭据往返 |
| 数据访问规划 | `database/sql + sqlc` | sqlc 查询生成与业务迁移在 Phase 1 实施 |

本机全局 Node 仍为 25.6.1；验证使用单独下载到临时目录的 Node 22.23.4，不修改用户系统 Node。Node 的依赖下限设置为 `>=22.13.0`，满足 Vite 7 和 pnpm 11 的共同要求。UI 的 Pinia/Tailwind 等依赖在各自实施阶段安装，不提前引入所有业务依赖。

## 3. 官方文档阅读与源码复核

Wails 官网直接访问触发 Cloudflare 403；已改用官方 GitHub API 获取 **v2.17.0 tag** 下的原始文档，解码并完整阅读：

- [Installation](https://github.com/wailsapp/wails/blob/v2.17.0/website/docs/gettingstarted/installation.mdx)，Git blob `6bdf220dee28b6701cdfed0df18f1a1d3c76f161`。
- [Creating a Project](https://github.com/wailsapp/wails/blob/v2.17.0/website/docs/gettingstarted/firstproject.mdx)，Git blob `973b03c58d14831833f7169519e4dde1a5aee8bf`。
- [Building](https://github.com/wailsapp/wails/blob/v2.17.0/website/docs/gettingstarted/building.mdx)，Git blob `c4c16ea48442d52b8f953e370d36d55335047e11`。

对应版本确认：Vue-TS 模板为 `wails init -t vue-ts`；`wails build` 输出到 `build/bin/`；Windows 使用 WebView2；macOS 使用 Xcode 工具链；Ubuntu 24.04 使用 WebKitGTK 4.1 并附加 `-tags webkit2_41`。

安装说明仍提及较早的 Go 最低版本，而 v2.17.0 模块实际声明 `go 1.25.0`。本项目依照真实模块和 SQLite 驱动要求固定 Go 1.26.2，不沿用旧页面的最低版本。

同时阅读本地同版本 `pkg/options/assetserver/options.go`、`pkg/options/options.go`、macOS binding dispatcher 与 Linux pkg-config 声明，确认：

- 静态资源找不到时可回退至 `AssetServer.Handler`，用于加载 Phase 0 iframe fixture。
- macOS dispatcher 在后端拒绝非主 frame 的绑定调用。该结论来自对应版本源码；其他平台仍需各自原生实测。
- 禁止将 `BindingsAllowedOrigins` 放开为任意来源。插件的生产权限检查仍需在 Phase 2 的 Go Runtime 实施。

## 4. 本机测试结果

| 验证项 | 结果 | 实际覆盖 |
|---|---|---|
| 依赖完整性 | 通过 | Go 模块校验；pnpm frozen lockfile 安装 |
| Go 测试 | 通过 | `go test -count=1 -v -timeout 90s ./...`，包含下述探针 |
| 前端类型检查/生产构建 | 通过 | Node 22.23.4 与本机 Node 25 均构建成功 |
| Wails 完整构建 | 通过 | 使用 Xcode SDK 后编译、打包、自签名成功 |
| 原生 WKWebView | 通过 | 应用启动后打印 `PHASE0_WEBVIEW_RESULT binding=true isolation=true` 并自动退出 |
| SQLite Go 驱动 | 通过 | 外键拒绝、事务回滚、迁移版本表最小创建、备份后重开及恢复 |
| SQLite CLI | 通过 | SQLite 3.54.0、WAL、事务、在线备份与恢复 |
| Chromium iframe | 通过 | 8 项隔离/通信断言，1 个 Playwright 测试 |
| 系统密钥环 | 通过（macOS） | 随机临时条目的写入/读取/删除；清理后确认不存在 |
| AI HTTP mock | 通过 | 本机普通 JSON 请求、SSE chunk、流中取消、401/429/502 错误传播 |
| C++ 独立进程 | 通过（macOS） | 编译/启动；进程组超时取消；子进程不能在取消后写标记文件 |
| Windows/Linux 测试包编译 | 通过 | CGO_ENABLED=0，Windows amd64 与 Linux amd64 的 Phase 0 测试包交叉编译成功；不代表原生桌面程序通过 |
| Windows/Linux 原生运行 | 通过 CI | 远程 Windows、macOS、Ubuntu 24.04 workflow 全部通过；本机仍只直接运行 macOS |

AI 测试只验证受控本机 HTTP 传输契约，不使用真实 API Key，不产生外部调用费用，**不代表已经验证 OpenAI/Anthropic/任一兼容服务的协议或账号**。真实 Provider 接入和错误归一化属于 Phase 3。

插件探针中的 `theme.get` 白名单、伪造身份拒绝和消息大小限制由 fixture 宿主执行，证明传输方案可行；不是已交付的 Go Plugin Runtime。生产环境的后端二次鉴权、安装生命周期、速率限制与持久化权限仍由 Phase 2 实现。

## 5. SDK 阻断与修复

初次 Wails/C++ 最终链接失败，错误为 SDK `.tbd` 含 `arm64e.x1` 架构。实际原因是 Xcode 21 工具链使用了较新的 Command Line Tools macOS 27.2 SDK。项目内 `scripts/phase0-macos-build.sh` 为该次构建选择 Xcode 自带 SDK，随后完整构建和 C++ 进程测试通过。

该脚本不修改系统 SDK、不编辑 `.tbd` 文件、不更改全局 xcode-select。`wails doctor` 在当前环境仍无法取得系统信息；它的失败没有被视为成功，但独立的真实构建和原生启动证据已经取得。

## 6. 跨平台验证边界

| 平台 | 桌面依赖 | 凭据服务 | 进程终止方案 | 当前状态 |
|---|---|---|---|---|
| Windows | WebView2 | Credential Manager | Job Object，挂起创建后加入 Job 再恢复 | CI 通过；真实交互式 WebView/密钥环仍需专门桌面测试 |
| macOS ARM64 | WKWebView、配套 Xcode SDK | Keychain | 独立 process group | 本机编译、运行、密钥环、进程组通过 |
| Linux | GTK3、WebKitGTK 4.1 | D-Bus Secret Service | process group；强限制后续评估 cgroup | CI 通过；真实交互式 WebView/密钥环仍需专门桌面测试 |

CI 默认不自动写系统密钥环，因为无图形会话时该服务可能不存在。没有密钥服务必须显式失败，不落盘明文。插件 iframe 不是 OS 级 CPU/内存隔离，也不能承诺任意恶意脚本不会冻结共享 WebView；这是插件阶段的剩余风险。
