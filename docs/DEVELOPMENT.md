# 开发与 Phase 0 复现

规范入口为根目录 `SPEC.md`。Phase 0、Phase 1 / Task 1 和 Task 2 已完成；每次只推进一个已授权阶段，测试和文档通过后汇报。

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

### pnpm 依赖构建授权

`frontend/pnpm-workspace.yaml` 显式允许锁文件中的 `esbuild@0.28.2` 执行安装脚本，并保留 `strictDepBuilds: true`。该文件需要和 package.json、pnpm-lock.yaml 一起提交，适用于本地、CI 和 Wails 的前端安装命令。

若 CI 报 `ERR_PNPM_IGNORED_BUILDS`，表示依赖下载后有未授权构建脚本，不是锁文件过期。升级 esbuild 时先审查版本和脚本，再同步更新该白名单；不要在 CI 中关闭检查或允许全部依赖执行脚本。

该修复已在 Node 22.23.4 / pnpm 11.26.0 下使用全新临时目录验证：没有配置时复现退出码 1；加入配置后 esbuild postinstall、冻结安装和前端生产构建均成功，锁文件内容不变。远端三平台结果以重新运行的 CI 为准。

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

`.github/workflows/phase0.yml` 在 Windows、macOS、Ubuntu 24.04 执行冻结安装、Go 测试、浏览器测试与原生构建，当前远程运行已全部通过。CI 的无头 Chromium 测试不替代各平台交互式 WebView 测试；Windows Job Object、Linux Secret Service 和跨平台安装包分发仍需后续专项验收。
