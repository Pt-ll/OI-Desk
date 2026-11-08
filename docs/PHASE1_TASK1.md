# Phase 1 / Task 1：后端持久化基础

## 范围

本任务建立应用启动所需的最小持久化基础，不包含 OJ、插件、AI 或业务表：

- 用户数据、缓存、日志和备份目录解析。
- 结构化 JSON 日志文件。
- SQLite 连接初始化：外键、WAL、busy timeout、单写连接。
- 迁移 runner 和 `schema_migrations`。
- `app_settings` 初始表和类型校验的设置存储。
- Wails `OnStartup`/`OnShutdown` 生命周期接入。

## 实现文件

- `backend/app/paths.go`：跨平台应用目录。
- `backend/app/runtime.go`：启动编排和资源关闭。
- `backend/app/logging.go`：`slog` JSON handler。
- `backend/database/db.go`：SQLite 连接策略。
- `backend/database/migrations.go`：嵌入式、有序、幂等迁移。
- `backend/database/migrations/0001_initial.sql`：`app_settings`。
- `backend/database/settings.go`：设置读写、JSON 和类型校验。
- `app.go`、`main.go`：Wails 生命周期接入。

## 数据约束

`app_settings.value_json` 必须是合法 JSON，`value_type` 仅允许 string、number、boolean、object、array、null。迁移重复执行不会重复应用；应用关闭时先关闭日志文件，再关闭数据库连接。测试和临时运行可以通过 `OIDESK_DATA_DIR` 指定目录，正式运行使用操作系统用户配置目录。

## 验证结果

- `go test -count=1 ./...` 通过。
- `pnpm run build` 通过。
- Wails v2.17.0 macOS 构建通过。
- 应用实际启动后创建临时数据库、日志和 `schema_migrations(1, initial)`，随后自动退出。
- Phase 0 的 SQLite、AI 流式、iframe、C++ 和现有 CI 探针保持通过。

## 后续边界

本任务没有把完整业务 schema 提前加入迁移；后续每个核心模块应增加独立迁移和测试。下一项建议是 **Phase 1 / Task 2：前端基础壳**，只实现主题、导航和版本化错误展示。

