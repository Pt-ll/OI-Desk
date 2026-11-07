# OI-Desk 技术调研与架构设计

本目录记录 `SPEC.md` 2.0.0 约束下的 Phase 0 调研结果。当前仓库只有规范文件，尚未初始化 Git，也没有 Go、Wails 或前端业务代码。本阶段只冻结设计决策、验证结论和后续实现边界。

文档索引：

- [TECHNICAL_VALIDATION.md](TECHNICAL_VALIDATION.md)：工具链、版本、下载验证和跨平台结论。
- [ARCHITECTURE.md](ARCHITECTURE.md)：系统架构、目录、模块依赖和技术选型。
- [DATABASE.md](DATABASE.md)：数据库缺陷、修正后的模型和 ER 图。
- [PLUGIN_API.md](PLUGIN_API.md)：iframe 隔离、Plugin API、生命周期、权限和通信流程。
- [OJ_API.md](OJ_API.md)：五个平台的适配器能力和同步接口。
- [ROADMAP.md](ROADMAP.md)：分阶段计划、验收门槛和风险清单。

## 当前状态

本次网络验证受到执行环境 DNS/连接超时影响，无法下载 Wails 模块、npm 元数据和官方网页。因此文档区分“本机已验证”“离线设计结论”和“Phase 0 必须重新在线验证”的内容。Wails 版本目标确定为 v2.10.2；该版本的下载、最小构建和官方文档复核仍是 Phase 0 的第一项门槛。

