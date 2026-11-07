# 分阶段开发计划与风险

## 1. Phase 0：技术验证与架构确认

本次交付为 Phase 0 的设计材料，不实现业务功能。用户确认后，Phase 0 的下一步是初始化 Git、锁定 Wails/Go/Node/npm 版本，创建最小 Wails 工程，完成通信、iframe、MessageChannel、Keychain、AI mock、流式和 C++ 进程验证，并记录三平台构建差异。

通过门槛：Wails v2.10.2 精确版本可复现；前端能 `install --frozen-lockfile` 和生产构建；SQLite 迁移/备份可恢复；插件越权请求被后端拒绝；流式请求可取消；超时不会残留 C++ 进程。

## 2. 后续阶段

| 阶段 | 单一目标 | 交付和验收 |
|---|---|---|
| Phase 1 | 基础框架 | Wails、Vue、SQLite、迁移、日志、设置、主题和导航；三平台构建方案明确 |
| Phase 2 | 插件系统 | Hello World ZIP 安装、iframe、Port、权限、独立存储、错误隔离 |
| Phase 3 | AI Gateway | Provider、密钥环、OpenAI-compatible、普通/流式/取消、用量和权限 |
| Phase 4 | OJ 数据中心 | 一个真实 OJ 完整同步闭环，再逐个平台增加 adapter |
| Phase 5 | 算法知识库 | 分类树、标签来源、掌握状态、按平台统计和搜索 |
| Phase 6 | 官方插件 | Snippets、Calculator、Notepad、Countdown、Stress Tester、Plotter，逐个验收 |
| Phase 7 | Dashboard 与完善 | 搜索、工作区、通知、快捷键、备份恢复、性能和稳定性 |

每个阶段只交付一个明确目标；格式化、相关测试、编译、安全检查和文档更新是阶段完成条件。未通过验收不得进入下一阶段。

## 3. 潜在技术风险

| 风险 | 影响 | 应对 |
|---|---|---|
| Wails v2/v3 API 混用 | 高 | 只锁一个主版本；绑定和构建脚本在 Phase 0 固化 |
| WebView/iframe 隔离差异 | 高 | sandbox、CSP、Port 会话和后端鉴权同时使用；三平台实测 |
| 插件恶意资源或路径穿越 | 高 | ZIP 规范化、临时目录、hash、大小上限、原子安装 |
| 系统密钥环不可用 | 高 | 明确报错和用户选择；不静默落盘明文 |
| OJ 接口变化/反爬/限流 | 高 | 能力声明、适配器隔离、错误分类、退避和 fixture；不伪造 API |
| SQLite 跨平台构建 | 中 | sqlc + modernc.org/sqlite；固定 WAL/外键/迁移；CI 验证 |
| AI Provider 协议差异 | 高 | Provider Adapter 和统一错误码；mock 测试流式/取消/限流 |
| AI 数据外泄 | 高 | 最小数据、逐插件权限、发送前展示、日志脱敏、无 key 暴露 |
| C++ 工具链或 SDK 不一致 | 高 | 当前 macOS 链接测试被 SDK `.tbd` 错误阻断；Phase 0 固定编译器/SDK 并在三平台 CI 验证 |
| C++ 子进程残留 | 高 | 每次任务独立进程组/Job Object，超时杀进程树并测试残留 |
| 大量提交阻塞 UI | 中 | 分页、后台任务、Wails 事件进度、索引和租约 |
| Node 25/非 LTS 差异 | 中 | Node 22 LTS 作为 CI 基线，锁 pnpm lockfile |
| 规范与实现演进不一致 | 中 | 每次阶段只改明确任务；结构变更先写原因、影响和迁移方案 |

## 4. 当前阻塞项

网络环境无法完成官方模块和 npm 下载，因此精确 Wails 补丁版本、依赖 lockfile 和真实平台 API 探针不能在本轮伪造完成。恢复网络后，优先执行 TECHNICAL_VALIDATION.md 的最小矩阵，再开始 Phase 0 实现。

