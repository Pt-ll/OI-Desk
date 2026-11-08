# Phase 1 / Task 2：前端基础壳

## 范围

本任务建立主界面的最小 UI 骨架，不接入 OJ、插件、AI 或业务数据：

- Pinia shell store：主题、当前导航、诊断状态和统一错误状态。
- Dashboard、OJ、Algorithms、Plugins、AI、Settings 六项导航占位。
- 深色默认主题和浅色主题切换，主题状态保存在前端 shell 设置中。
- 统一 `AppError` 类型和 `normalizeError`，错误展示带稳定错误码。
- 保留 Phase 0 iframe/MessageChannel 诊断探针，避免基础壳回归隔离验证。

## 实现文件

- `frontend/src/stores/shell.ts`：Pinia 状态和导航定义。
- `frontend/src/types/errors.ts`：版本化错误模型和归一化函数。
- `frontend/src/components/SidebarNav.vue`：主导航。
- `frontend/src/components/ErrorNotice.vue`：统一错误提示。
- `frontend/src/App.vue`、`frontend/src/style.css`：基础工作区布局和主题样式。
- `frontend/tests/shell.spec.ts`：导航和主题 Playwright 测试。

## 验证结果

- `vue-tsc --noEmit && vite build` 通过。
- Playwright 的 Phase 0 iframe 测试和 Task 2 导航/主题测试均通过（2 tests passed）。
- `go test -count=1 ./...` 通过。
- Wails v2.17.0 macOS 原生构建通过。

主题目前由 shell store 管理，后续应通过 Phase 1 的 Wails 设置 DTO 接入后端 `app_settings`；在此之前不把 UI 偏好当作核心数据。导航项是稳定入口定义，业务页面将在对应阶段实现。

