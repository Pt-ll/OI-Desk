# 开发与 Phase 0 复现

规范入口为根目录 `SPEC.md`。当前仅有技术验证工程；每次只推进一个已授权阶段，测试和文档通过后汇报。Phase 1 尚未开始。

## 工具链

- Go 1.26.2；Wails CLI v2.17.0。
- Node 22.23.4（见 `.node-version`）；pnpm 11.26.0。
- macOS 使用配套 Xcode SDK；Linux 安装 GTK3/WebKitGTK 4.1；Windows 安装 WebView2。

```sh
go install github.com/wailsapp/wails/v2/cmd/wails@v2.17.0
pnpm --dir frontend install --frozen-lockfile
go mod verify
```

Go 官方代理不可达时可为单次命令设置 `GOPROXY=https://goproxy.cn,direct`。保留 Go checksum 校验，勿关闭 GOSUMDB，勿提交本地路径 replace。

## 测试

```sh
pnpm --dir frontend build
go test -v -timeout 90s ./...
pnpm --dir frontend exec playwright install chromium
cd frontend && pnpm test:phase0
```

默认 Go 测试跳过会触碰系统密钥环和编译器的 opt-in 探针。macOS 显式验证：

```sh
export SDKROOT="$(xcode-select -p)/Platforms/MacOSX.platform/Developer/SDKs/MacOSX.sdk"
OIDESK_KEYRING_PROBE=1 OIDESK_CPP_PROBE=1 go test -count=1 -v -timeout 90s ./phase0
./scripts/phase0-sqlite-check.sh
```

密钥环探针只创建随机命名的 `org.oidesk.phase0.probe` 临时条目并在结束时删除，使用虚拟值，不使用真实凭据。

## 构建和原生验证

```sh
# macOS（同时处理本机 CLT/Xcode SDK 不一致）
./scripts/phase0-macos-build.sh
OIDESK_PHASE0_AUTO_EXIT=1 ./build/bin/OI-Desk.app/Contents/MacOS/oidesk

# Linux with WebKitGTK 4.1
wails build -tags webkit2_41

# Windows
wails build
```

原生验证成功输出：`PHASE0_WEBVIEW_RESULT binding=true isolation=true`。该输出由页面收到 Go 绑定结果和 iframe 测试结果后回传；单纯构建成功不等于运行成功。15 秒没有完成将打印超时，不能视为通过。

## CI 和未验证范围

`.github/workflows/phase0.yml` 在 Windows、macOS、Ubuntu 24.04 执行冻结安装、Go 测试、浏览器测试与原生构建。当前只写入本地 workflow，尚未推送或取得远端运行结果。CI 的无头 Chromium 测试不替代各平台原生 WebView 测试。Windows Job Object、Linux Secret Service 和跨平台安装包分发仍需后续验收。
