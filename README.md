# OI-Desk

Go + Wails + Vue + SQLite 的 OI 训练工作台。产品和开发约束以 [SPEC.md](SPEC.md) 为准。

当前已完成 **Phase 0 技术验证**和 **Phase 1 基础框架**（Task 1–3）；尚未实现 OJ、插件或 AI 业务功能。Wails 固定为 **v2.17.0**。

- [架构设计](docs/ARCHITECTURE.md)
- [技术验证结果](docs/TECHNICAL_VALIDATION.md)
- [Phase 1 / Task 1](docs/PHASE1_TASK1.md)
- [Phase 1 / Task 2](docs/PHASE1_TASK2.md)
- [Phase 1 / Task 3](docs/PHASE1_TASK3.md)
- [开发与验证命令](docs/DEVELOPMENT.md)
- [数据库审查](docs/DATABASE.md)
- [插件设计](docs/PLUGIN_API.md)
- [开发路线](docs/ROADMAP.md)

macOS 构建：`./scripts/phase0-macos-build.sh`。需要 Go 1.26.2、Node 22.23.4、pnpm 11.26.0 和 Xcode 工具链。
