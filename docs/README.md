# OI-Desk 技术调研与架构设计

本目录记录 `SPEC.md` 2.0.0 约束下的 Phase 0 调研结果。当前仓库已完成 Git 初始化和 Wails v2.17.0 最小工程初始化；本阶段只验证基础技术链路，不实现业务功能。

文档索引：

- [TECHNICAL_VALIDATION.md](TECHNICAL_VALIDATION.md)：工具链、版本、下载验证和跨平台结论。
- [ARCHITECTURE.md](ARCHITECTURE.md)：系统架构、目录、模块依赖和技术选型。
- [DATABASE.md](DATABASE.md)：数据库缺陷、修正后的模型和 ER 图。
- [PLUGIN_API.md](PLUGIN_API.md)：iframe 隔离、Plugin API、生命周期、权限和通信流程。
- [OJ_API.md](OJ_API.md)：五个平台的适配器能力和同步接口。
- [AI_GATEWAY.md](AI_GATEWAY.md)：AI 权限、凭据、流式与取消设计。
- [DEVELOPMENT.md](DEVELOPMENT.md)：开发环境和验证命令。
- [PHASE1_TASK1.md](PHASE1_TASK1.md)：Phase 1 / Task 1 持久化基础交付。
- [PHASE1_TASK2.md](PHASE1_TASK2.md)：Phase 1 / Task 2 前端基础壳交付。
- [ROADMAP.md](ROADMAP.md)：分阶段计划、验收门槛和风险清单。

## 当前状态

Wails 固定为 v2.17.0。Phase 0、Phase 1 / Task 1 和 Task 2 已完成；本机 Go/Vue/SQLite、Keychain、C++ 进程组、Chromium 和原生 WKWebView 验证通过；Windows、macOS、Ubuntu 24.04 三平台 CI 已全部通过。交互式 WebView、密钥环和安装包分发仍需专项验收。

