# Phase 1 / Task 3：设置绑定与 Phase 1 收尾

## 范围

本任务把前端基础壳和后端持久化连接起来，完成 Phase 1 的设置闭环：

- Wails 暴露版本化 `SettingDTO`。
- `GetSetting`、`SetSetting`、`DeleteSetting` 通过后端 `app_settings` 工作。
- 前端主题启动时优先读取后端设置，浏览器/首次启动失败时使用本地 fallback。
- 主题切换立即更新 UI，并异步写回后端。
- 保留统一错误处理和 Phase 0 诊断探针。

## 验证结果

- `go test -count=1 ./...` 通过，包含 `SettingDTO` 往返测试。
- `vue-tsc --noEmit && vite build` 通过。
- Playwright 导航、主题和 iframe 隔离测试 2 项通过。
- Wails v2.17.0 macOS 原生构建通过。
- 前端依赖使用 pnpm frozen lockfile 安装通过。

## Phase 1 验收状态

`SPEC.md` Phase 1 的基础框架目标已完成：Wails 应用、Vue 主界面、Go 后端、SQLite、迁移、日志、配置路径、主题和基础导航均已具备。OJ、插件、AI、算法和 Dashboard 业务仍保持在后续阶段边界内。

下一阶段建议进入 **Phase 2：插件系统**，从 Hello World ZIP 安装、manifest 校验和后端权限拒绝开始。

