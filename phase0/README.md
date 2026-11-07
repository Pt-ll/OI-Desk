# Phase 0 技术探针

这些文件只验证基础能力，不属于业务功能：

- `stream_test.go`：使用标准库 `httptest` 验证 SSE 流式响应、取消和错误传播。
- `sqlite_test.go`：真实 Go SQLite 驱动的事务、外键和备份恢复。
- `keyring_test.go`：显式启用后验证临时密钥环条目的读写和清理。
- `cpp_unix_test.go`：Unix C++ 独立进程与子进程超时清理。
- `plugin-host.html` 与 `plugin.html`：验证 sandboxed iframe、nonce 握手和 `MessageChannel` 请求/响应。
- `../scripts/phase0-sqlite-check.sh`：验证 SQLite 外键、WAL、事务和备份文件。
- `../scripts/phase0-cpp-check.sh`：运行 C++ 独立编译、超时和进程组清理测试；macOS 可设置 SDKROOT 选择配套 Xcode SDK。

浏览器探针通过 `frontend/tests/phase0.spec.ts` 自动执行，并集成到 Wails 验证页面。插件不获得 `allow-same-origin`，所有示例能力通过 Port 提供。8 项断言只证明最小隔离与传输可行，不代替 Phase 2 的 Go 后端鉴权。

