# Phase 2 / Task 1：插件包与 Manifest 安全校验

## 范围

本任务只处理插件进入运行时之前的静态校验：

- Manifest JSON 解析。
- 插件 ID、语义化版本、核心版本下限和入口文件校验。
- 已注册权限与 AI capability 白名单。
- 入口/icon 相对路径和插件目录边界校验。
- ZIP 路径穿越、反斜杠混淆、符号链接、重复文件、文件数量和解压大小限制。
- 权限去重和稳定排序。

本任务没有安装、卸载、启停、iframe 会话、插件存储或后端 Plugin API；这些能力保持为后续独立任务。

## 实现文件

- `backend/plugins/manifest.go`：Manifest 模型、版本比较、权限白名单和资源校验。
- `backend/plugins/package.go`：ZIP 检查、大小限制和内存文件系统视图。
- `backend/plugins/manifest_test.go`：有效 Manifest、非法权限、路径穿越和合法 ZIP 测试。

## 安全约束

- 未知 Manifest 字段被忽略，不会被执行。
- 插件不能声明未注册权限；AI capability 不能在 `ai.enabled=false` 时启用。
- ZIP 中禁止绝对路径、`..` 路径、Windows 反斜杠路径、符号链接和重复路径。
- 解析阶段只建立内存视图，不把不可信 ZIP 内容写入正式插件目录。
- 默认上限为 2048 个文件和 64 MiB 解压后内容；后续安装器仍需原子替换和 hash/签名策略。

## 验证结果

`go test -count=1 ./backend/plugins` 通过。下一项建议是 **Phase 2 / Task 2：插件注册表与安全安装事务**，把已通过静态检查的包安装到临时目录并写入注册表。

