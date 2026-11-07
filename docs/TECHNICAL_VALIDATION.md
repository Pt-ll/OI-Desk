# 技术验证报告

## 1. 验证范围

本报告对应 `SPEC.md` 的 Phase 0 要求：确认 Wails、Go、Vue 3、TypeScript、SQLite、插件隔离、系统密钥存储、AI Provider、流式响应和 C++ 独立进程的可行性。当前只做环境和设计验证，不创建应用代码。

## 2. 仓库状态

执行目录为 `/Users/pt.ll/Desktop/OIDesk`。目录中只有 `SPEC.md`，`git status` 返回 `fatal: not a git repository`，说明仓库尚未初始化。没有发现 `AGENTS.md`、`go.mod`、`package.json`、Wails 工程或现有迁移文件。

## 3. 本机工具链

| 工具 | 本机结果 | 结论 |
|---|---:|---|
| OS | macOS 27.2，Darwin ARM64 | 当前验证主机是 Apple Silicon；不能代替 Windows/Linux 构建验证 |
| Git | 2.53.0 | 已安装，仓库尚未初始化 |
| Go | 1.26.2 | 可作为当前开发机工具；最终版本写入 `go.mod` 和 CI |
| Node.js | 25.6.1 | 可运行 JavaScript 工具，但不是推荐的 LTS 基线 |
| pnpm | 11.26.0 | 可作为包管理器，提交锁文件 |
| SQLite CLI | 3.54.0 | 本机 SQLite 可用 |
| C/C++ 编译器 | Apple Clang 21 | 可验证对拍编译流程；Windows/Linux 编译器仍需 CI 验证 |
| Wails CLI | 未安装 | Phase 0 需要安装并锁定 |

本机 Go 代理环境变量等价于错误的 `https://goproxy.cn,direc`，导致 `go list -m` 直接失败。修正为 `https://proxy.golang.org,direct` 后，官方代理和 `wails.io` 仍发生网络连接超时。npm registry 查询也未能完成。没有把失败的网络结果当成版本或 API 证据。

SQLite CLI 的外键和 WAL 最小测试通过（SQLite 3.54.0）。Node 25.6.1 提供实验性的 `node:sqlite`，但项目不依赖它；缓存中可见 TypeScript 5.9.3。Apple Clang 能生成 C++17 object file，但链接一个最小可执行文件时遇到本机 macOS 27 SDK `.tbd` 的 `unknown architecture` 错误，因此对拍器的完整编译运行必须在 Phase 0 修复 SDK 或转入 CI 后再验收。

## 4. 版本选择

### Wails

选择 **Wails v2.10.2**。这是 v2 稳定线路的精确目标版本；Phase 0 网络恢复后必须用模块下载、最小工程构建和官方文档再次确认该 tag 可用，再写入 `go.mod` 和工具链锁定文件。选择理由：

1. SPEC 要求在 v2/v3 中二选一；v2 的文档、模板和生态更成熟。
2. v2 的绑定模型、运行时事件和 asset server 足够承载 Vue 前端及受控插件 iframe。
3. v3 在设计阶段存在较高 API 迁移风险，不能与 v2 API 混用。

官方文档基线应为：

- https://wails.io/docs/v2.10/
- https://wails.io/docs/v2.10/gettingstarted/installation
- https://wails.io/docs/v2.10/howdoesitwork/
- https://wails.io/docs/v2.10/reference/runtime/

由于本次无法访问官方站点，官方文档阅读和 v2.10.2 下载验证无法在本轮完成，必须在 Phase 0 的最小 Wails 工程中复核。`go.mod` 只允许出现一个 Wails 主版本。

### Go、前端和 SQLite

最终锁定策略如下：

- Go：选择 Go 1.26.2 作为当前开发基线；Phase 0 用 Wails v2.10.2 的真实构建确认兼容性，并在 CI 固定同一补丁版本。
- Node.js：使用 Node 22 LTS 作为开发和 CI 基线，避免把 Node 25 的非 LTS 行为带入跨平台构建；`engines` 约束为 `>=22.12.0`。
- Vue：Vue 3.5.x。
- TypeScript：TypeScript 5.9.x。
- Vite：Vite 6.x，使用与 Node 22 LTS 相容的补丁版本；Wails 前端构建通过 Vite 输出静态资源。
- SQLite：SQLite 3.x；应用层使用 `database/sql + sqlc + modernc.org/sqlite`。这样数据库驱动不额外依赖 CGO，Wails 本身所需的 WebView 构建依赖仍按平台准备。
- 包管理：pnpm 11，提交 `pnpm-lock.yaml`。

Vue、TypeScript 和 Vite 的精确补丁版本必须在联网的 Phase 0 中通过 `pnpm add`、`pnpm install --frozen-lockfile` 和生产构建确认后写入锁文件。不要使用浮动的 `latest`。

### 为什么不使用 GORM 或 mattn/go-sqlite3

SPEC 允许 GORM 或 sqlc。这里选 sqlc，因为查询和迁移可以保持显式、可审查，并能把平台一致性约束放在 SQL 中。`mattn/go-sqlite3` 会把 CGO 编译链带入数据库层，增加 Windows 和交叉编译的变数；Wails 的 WebView 依赖已经需要按平台准备，因此没有必要再扩大数据库驱动的 CGO 面。

## 5. 最小验证矩阵

| 项目 | 最小验证 | 通过标准 |
|---|---|---|
| Wails | 创建空 v2 工程、`wails dev`、`wails build` | Go 与 Vue 绑定可调用，三平台构建命令明确 |
| Vue/TS/Vite | `pnpm install --frozen-lockfile`、类型检查、生产构建 | 无类型错误，静态资源能被 Wails asset server 加载 |
| SQLite | 打开数据库、启用外键、迁移、事务、WAL、备份 API | 重启后数据存在，迁移可重复执行，备份可恢复 |
| iframe | `sandbox="allow-scripts"`、CSP、禁止父页面访问 | 插件只能通过 Port 调用允许 API |
| MessageChannel | 请求/响应、超时、取消、流式事件 | 关联 ID 正确，断开可回收，超大消息被拒绝 |
| Keychain | macOS Keychain；Windows Credential Manager；Linux Secret Service | 密钥只返回给 Go Provider，插件和日志不可见 |
| AI | 一个 OpenAI-compatible mock server | 普通、流式、取消、认证/限流/超时错误均可分类 |
| C++ | 编译、启动、超时、输出上限、进程树终止 | 主程序不阻塞，不残留子进程 |

## 6. 跨平台结论

Wails 的前端代码可保持跨平台，但系统能力必须分层：

- Windows：WebView2 Runtime、Credential Manager、Job Object 进程组、MSVC 或 LLVM/MinGW。
- macOS：WKWebView、Keychain Services、进程组/信号、Apple Clang；需要处理应用签名和 hardened runtime。
- Linux：WebKitGTK 4.1、Secret Service（如可用）、process group/cgroup；发行版需要明确 WebKitGTK 开发包。

统一接口由 Go 定义，平台实现放在 `backend/platform/{windows,darwin,linux}`。没有系统密钥服务时必须报错并让用户选择退出或显式配置的安全外部凭据方案，不能静默写入 SQLite。

## 7. 需要在 Phase 0 重新确认的事项

1. Wails v2.10.x 的最新补丁版本及其支持的 Go 版本。
2. Wails v2 asset server 下 iframe 的实际来源、CSP 和沙箱行为。
3. Vue 3.5.x、TypeScript 5.9.x、Vite 6.x 的 npm peer dependency。
4. `modernc.org/sqlite` 与选定 Go 补丁版本的构建结果。
5. 三个平台 WebView、密钥环和 C++ 进程组的 CI 构建。

