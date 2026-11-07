# OI-Desk — Project Specification

> **文档版本：** 2.0.0  
> **项目代号：** OI-Desk  
> **项目类型：** 跨平台桌面应用 / OI 训练工作台  
> **目标平台：** Windows、macOS、Linux  
> **核心技术栈：** Go + Wails + Vue 3 + TypeScript + SQLite  
> **文档用途：** 产品需求、系统架构、插件规范、AI 接入规范、数据库设计及 AI 编程 Agent 开发约束

---

# 0. 项目概述

## 0.1 项目定位

OI-Desk 是面向信息学竞赛选手的本地个人训练数据中心和可扩展桌面工作台。

它不是 OJ，不是代码编辑器，也不是单纯的工具集合。它希望把选手日常训练中分散的事情集中起来：

- 管理多个 OJ 的账号、提交记录和题目。
- 汇总不同 OJ 的做题历史和算法掌握情况。
- 保存、搜索和复用代码模板。
- 提供本地对拍、计算器、画图、记事本等工具。
- 管理比赛日期、报名截止时间和训练计划。
- 通过插件系统扩展工作台。
- 通过统一 AI 服务为核心和插件提供可配置的智能能力。

## 0.2 核心设计原则

1. **本地优先**：用户数据默认保存在本地，不依赖云端服务。
2. **简单实用**：优先解决高频需求，不为复杂架构牺牲使用体验。
3. **插件化**：工具应当能够独立安装、更新、禁用和卸载。
4. **数据安全**：插件之间默认隔离，敏感凭据由核心管理。
5. **AI 可配置**：用户自行决定 AI 服务商、模型、密钥和调用权限。
6. **平台独立**：不同 OJ 的原始数据和统计独立保存，不强制合并。
7. **可扩展**：核心功能、插件 API 和数据库结构需要支持后续扩展。
8. **跨平台**：Windows、macOS、Linux 均应有合理的实现方案。
9. **可恢复**：重要数据必须支持备份、导入和恢复。

## 0.3 第一阶段目标

第一阶段必须交付：

- 桌面应用基础框架。
- SQLite 数据库和迁移机制。
- OJ 账号管理及至少一个真实 OJ 的同步实现。
- 题目、提交记录和算法标签管理。
- 插件管理系统。
- 插件权限与安全隔离。
- 自定义 AI 服务配置及统一 AI Gateway。
- Notepad、Calculator、Snippets、Countdown、Stress Tester 等基础插件。
- Dashboard、全局搜索、设置和备份功能。

Plotter 可以作为第一阶段后半部分或第二阶段插件。

---

# 1. 技术选型

## 1.1 后端

| 技术 | 用途 |
|---|---|
| Go | 核心业务逻辑 |
| Wails | 桌面应用框架 |
| SQLite | 本地数据库 |
| GORM 或 sqlc | 数据访问层，二选一 |
| slog | 日志 |
| OS Keychain | API Key、OJ 凭据存储 |
| Go testing | 后端测试 |

Wails 必须在开发初期确定使用 v2 或 v3，并统一使用对应版本的 API。禁止混用两个版本的代码。

## 1.2 前端

| 技术 | 用途 |
|---|---|
| Vue 3 | UI 框架 |
| TypeScript | 类型系统 |
| Vite | 构建工具 |
| Tailwind CSS | 样式 |
| Lucide | 图标 |
| Pinia | 前端状态管理 |
| ECharts | 统计图表 |
| Monaco Editor | 代码编辑与高亮 |
| math.js | 数学表达式计算 |
| Plotly.js | 函数绘图 |

具体依赖版本必须在初始化时核实兼容性，并锁定版本。

## 1.3 对拍引擎

对拍工具使用独立进程运行 C++ 程序。

Go 负责：

- 创建任务工作目录。
- 管理编译器配置。
- 调用编译器。
- 启动和终止进程。
- 控制超时和输出大小。
- 收集 stdout、stderr 和退出码。
- 比较程序输出。
- 保存错误数据。

不得将用户提交的 C++ 代码直接编译进 Go 主进程。

## 1.4 开发工具

- Git。
- Go toolchain。
- Node.js。
- pnpm。
- Wails CLI。
- GCC / Clang。
- SQLite。
- Playwright。
- Vitest。

---

# 2. 系统总体架构

## 2.1 架构分层

系统分为五层：

1. **Core**：核心业务和本地数据管理。
2. **Adapters**：不同 OJ 的数据接入层。
3. **Plugin Runtime**：插件生命周期、运行环境和权限控制。
4. **AI Gateway**：统一 AI 服务接入与调用代理。
5. **Frontend**：主界面、工作区和插件 UI。

```text
OI-Desk
│
├── Frontend
│   ├── Dashboard
│   ├── OJ Management
│   ├── Algorithm Center
│   ├── Plugin Manager
│   ├── AI Settings
│   └── Plugin Workspace
│
├── Core
│   ├── Database
│   ├── OJ Accounts
│   ├── Submissions
│   ├── Problems
│   ├── Algorithms
│   ├── Contest Manager
│   ├── Backup
│   └── Search
│
├── OJ Adapters
│   ├── Luogu
│   ├── QOJ
│   ├── Codeforces
│   ├── UOJ
│   └── AtCoder
│
├── Plugin Runtime
│   ├── Manifest Loader
│   ├── Plugin Registry
│   ├── Permission Manager
│   ├── Sandbox
│   ├── Plugin Storage
│   └── Plugin API
│
├── AI Gateway
│   ├── Provider Manager
│   ├── Model Manager
│   ├── Credential Manager
│   ├── Permission Checker
│   ├── Request Router
│   ├── Streaming Handler
│   └── Usage Manager
│
└── Local Storage
    ├── SQLite
    ├── Plugin Data
    ├── Plugin Files
    ├── Cache
    ├── Logs
    └── Backups
```

## 2.2 模块依赖原则

- Frontend 不得直接操作 SQLite。
- 插件不得直接访问核心数据库。
- OJ Adapter 不得直接操作 UI。
- AI Gateway 不得绕过权限系统。
- 插件不得直接读取其他插件的私有数据。
- 插件不得直接读取核心保存的 OJ Cookie 或 AI API Key。
- 核心模块不得依赖某个具体插件。
- 所有跨模块调用必须经过明确的接口。

## 2.3 核心服务接口

建议使用依赖注入组织服务：

```go
type AppServices struct {
    Database       DatabaseService
    OJ             OJService
    Algorithms     AlgorithmService
    Plugins        PluginService
    AI             AIService
    Contests       ContestService
    Backup         BackupService
    Notifications  NotificationService
}
```

具体接口根据模块拆分，不要求所有业务都塞进一个巨大的 Service。

---

# 3. OJ 数据管理

## 3.1 支持的平台

第一批支持：

- 洛谷 Luogu
- QOJ
- Codeforces
- UOJ
- AtCoder

不同 OJ 必须有独立 Adapter。

## 3.2 OJ 账号管理

每个平台可以添加多个账号。

账号信息：

- 平台。
- 用户名。
- 用户 ID。
- 显示名称。
- 头像。
- 添加时间。
- 最近同步时间。
- 当前认证状态。

Cookie、Session、Token 等敏感信息不能明文写入 SQLite、日志或插件存储。

敏感凭据优先保存到：

- Windows Credential Manager。
- macOS Keychain。
- Linux Secret Service / 系统密钥环。

如果系统密钥服务不可用，应明确告知用户，并提供安全的替代方案；不得静默退化成明文保存。

## 3.3 同步功能

支持：

- 手动同步。
- 启动时同步。
- 定时同步。
- 增量同步。
- 失败后重试。
- 同步状态展示。
- 认证失效提醒。

同步必须满足：

1. 有稳定的去重机制。
2. 支持中断后继续。
3. 支持分页。
4. 遵守 OJ 的访问限制。
5. 不能把网络错误当成没有新提交。
6. 不能假设所有 OJ 都有公开 API。
7. 不得伪造 OJ 接口、字段或认证方式。

Adapter 接口：

```go
type Platform interface {
    Name() string

    CheckAuth(ctx context.Context) error

    GetUserInfo(
        ctx context.Context,
    ) (*UserInfo, error)

    GetSubmissions(
        ctx context.Context,
        cursor string,
    ) (*SubmissionPage, error)
}
```

具体平台可扩展题目详情、标签、比赛等接口。

## 3.4 提交记录

保存：

- OJ。
- 远程提交 ID。
- 题目 ID。
- 提交时间。
- 语言。
- 状态。
- 运行时间。
- 内存使用。
- 代码长度。
- 其他平台提供的原始信息。

唯一约束：

```sql
UNIQUE(platform, remote_id)
```

平台未提供某些字段时允许为空，不得伪造数据。

## 3.5 题目管理

每道题目独立归属于一个 OJ。

唯一约束：

```sql
UNIQUE(platform, remote_id)
```

不同 OJ 的同名题目默认不合并。

题目字段：

- 平台。
- 远程 ID。
- 标题。
- URL。
- 难度。
- 题目类型。
- 标签。
- 创建时间。
- 最近提交时间。

---

# 4. 算法知识库

## 4.1 算法体系

维护一棵全局算法分类树。

例如：

```text
算法
├── 数据结构
│   ├── 线段树
│   ├── 树状数组
│   ├── 可持久化数据结构
│   └── 平衡树
│
├── 图论
│   ├── 最短路
│   ├── 网络流
│   ├── 强连通分量
│   └── 树上问题
│
├── 动态规划
│   ├── 背包
│   ├── 区间 DP
│   ├── 状压 DP
│   └── 树形 DP
│
├── 数学
│   ├── 数论
│   ├── 组合数学
│   ├── 线性代数
│   └── 概率论
│
└── 字符串
    ├── KMP
    ├── Trie
    ├── 后缀数组
    └── 哈希
```

算法分类允许用户自行新增、修改、删除和调整层级。

## 4.2 题目与算法关系

一题可以对应多个算法。

算法与题目通过多对多关系连接。

标签来源：

- manual：用户手动标注。
- oj_tag：OJ 提供的标签。
- ai：AI 推测。
- imported：导入数据。

AI 标签不得覆盖手动标签。

## 4.3 平台独立统计

算法分类体系可以全局共用，但题目关联和统计必须按 OJ 独立。

例如：

```text
洛谷
  线段树：30 题
  树形 DP：15 题

QOJ
  线段树：12 题
  树形 DP：8 题

Codeforces
  线段树：25 题
  树形 DP：9 题
```

允许未来增加跨平台总览，但必须保留各平台独立统计。

## 4.4 掌握状态

每个题目和算法的学习状态可以是：

- 未接触。
- 尝试中。
- 已通过。
- 已掌握。
- 待复习。

状态变化应当可以由用户主动修改。

---

# 5. 插件系统

## 5.1 设计目标

插件是 OI-Desk 的重要组成部分。

记事本、计算器、代码片段、画图、比赛倒计时、对拍工具等都应当是独立插件，而不是全部写死在主程序中。

插件必须支持：

- 安装。
- 卸载。
- 启用。
- 禁用。
- 更新。
- 配置。
- 权限管理。
- 独立数据存储。
- 插件错误隔离。
- 插件版本管理。

## 5.2 插件技术方案

第一版采用 Web 插件。

插件由 HTML、CSS、JavaScript 或 TypeScript 编译产物组成。

每个插件使用独立目录：

```text
plugins/
├── notepad/
│   ├── manifest.json
│   ├── index.html
│   ├── assets/
│   └── icon.png
│
├── calculator/
│   ├── manifest.json
│   ├── index.html
│   └── assets/
│
└── snippets/
    ├── manifest.json
    ├── index.html
    └── assets/
```

插件 UI 在独立 iframe 中运行。

推荐采用独立来源或严格隔离的 iframe 沙箱，并使用 `MessageChannel` 建立通信。

插件不得直接访问：

- 核心 Go 对象。
- 核心 SQLite。
- 其他插件的存储目录。
- OJ 凭据。
- AI API Key。
- 任意本地文件。

插件需要的能力必须通过核心提供的 Plugin API 获取。

## 5.3 Manifest 规范

每个插件必须有 `manifest.json`。

示例：

```json
{
  "id": "org.oidesk.notepad",
  "name": "Notepad",
  "version": "1.0.0",
  "description": "Markdown note editor",
  "author": "OI-Desk",
  "entry": "index.html",
  "icon": "icon.png",
  "minCoreVersion": "1.0.0",
  "permissions": [
    "storage",
    "notification",
    "clipboard"
  ],
  "ai": {
    "enabled": true,
    "capabilities": [
      "chat",
      "summarize"
    ],
    "defaultModel": null
  }
}
```

Manifest 必须经过严格校验。

要求：

- 插件 ID 全局唯一。
- 版本号符合语义化版本。
- 入口文件必须存在。
- 权限只能使用已注册权限。
- 不允许插件自行声明任意系统权限。
- 不允许通过路径穿越访问插件目录外文件。
- 未知字段应按兼容性策略处理，不得直接执行。

## 5.4 插件权限系统

权限采用最小权限原则。

第一批权限：

| 权限 | 能力 |
|---|---|
| storage | 读写自身插件数据 |
| clipboard | 访问剪贴板 |
| notification | 发送系统通知 |
| theme | 读取当前主题 |
| contest.read | 读取比赛信息 |
| problem.read | 读取题目信息 |
| submission.read | 读取提交记录 |
| algorithm.read | 读取算法分类 |
| ai.chat | 调用 AI 对话 |
| ai.summarize | 调用 AI 摘要能力 |
| ai.embedding | 调用向量嵌入能力 |
| network | 访问受控网络接口 |
| file.export | 导出用户主动选择的数据 |
| file.import | 导入用户主动选择的文件 |

权限处理规则：

1. 安装插件时展示权限说明。
2. 高风险权限必须单独确认。
3. 用户可以随时撤销权限。
4. 撤销权限后立即拒绝新的调用。
5. 权限不足时返回明确错误。
6. 插件不得自行修改权限。
7. 核心后端必须再次检查权限，不能只依赖前端判断。

`network` 不等于无限制网络访问。后续应考虑域名白名单、请求类型限制和网络访问审计。

## 5.5 插件生命周期

标准生命周期：

```text
发现插件
  ↓
解析 Manifest
  ↓
校验版本和文件
  ↓
检查权限
  ↓
安装
  ↓
注册插件
  ↓
创建运行环境
  ↓
加载 UI
  ↓
建立通信
  ↓
运行
  ↓
禁用 / 卸载
```

插件应支持：

- `onInstall`
- `onEnable`
- `onDisable`
- `onUninstall`

生命周期钩子应有超时和错误处理。

插件崩溃不得导致主程序退出。

## 5.6 插件通信协议

推荐使用 `MessageChannel`。

所有消息采用统一结构：

```ts
interface PluginMessage {
  id: string;
  pluginId: string;
  type: string;
  method: string;
  payload?: unknown;
}

interface PluginResponse {
  id: string;
  success: boolean;
  data?: unknown;
  error?: {
    code: string;
    message: string;
  };
}
```

通信要求：

- 消息必须有唯一 ID。
- 请求必须支持超时。
- 响应必须能关联原请求。
- 必须校验插件身份。
- 必须检查权限。
- 必须限制消息大小。
- 必须限制并发请求数。
- 必须对输入进行 schema 校验。
- 不允许任意执行插件传入的函数名。

插件 API 版本应独立管理，保证旧插件在一定范围内兼容。

## 5.7 插件存储

每个插件有自己的数据目录：

```text
data/
├── core/
│   └── oidesk.db
│
├── plugins/
│   ├── org.oidesk.notepad/
│   ├── org.oidesk.snippets/
│   └── org.oidesk.plotter/
│
├── cache/
└── backups/
```

规则：

- 插件只能访问自身存储。
- 核心不能随意清理插件数据。
- 禁用插件不删除数据。
- 卸载插件默认保留数据。
- 用户主动确认后才能清理插件数据。
- 插件升级必须保留原数据。
- 插件数据迁移必须有版本号和备份策略。

## 5.8 插件安装和更新

第一版支持本地 ZIP 安装。

安装过程：

1. 读取 ZIP。
2. 检查文件路径。
3. 防止路径穿越。
4. 校验 Manifest。
5. 检查入口文件。
6. 检查版本兼容性。
7. 展示权限。
8. 用户确认。
9. 安装到插件目录。
10. 注册插件。

第一版不要求在线插件市场。

未来插件市场需要支持：

- 插件列表。
- 版本信息。
- 更新检查。
- 作者信息。
- 权限说明。
- 下载来源校验。
- 签名验证。
- 下架和撤回机制。

---

# 6. 自定义 AI 接入系统

## 6.1 设计目标

OI-Desk 必须支持用户自定义 AI 服务，而不是将系统绑定到单一 AI 服务商。

用户可以：

- 添加多个 AI Provider。
- 配置不同 API Endpoint。
- 配置 API Key。
- 配置模型名称。
- 设置默认模型。
- 为不同用途指定不同模型。
- 测试连接。
- 查看调用状态和用量。
- 为不同插件分配 AI 权限。
- 随时禁用某个 AI Provider。

AI 接入必须由核心统一管理。

**插件只能请求 AI Gateway，不能自行获取 API Key 或绕过核心直接调用用户配置的 AI 服务。**

## 6.2 支持的 Provider 类型

第一版至少支持以下类型：

### A. OpenAI-compatible

支持常见 OpenAI 兼容接口。

配置项：

- Provider 名称。
- Base URL。
- API Key。
- Chat Completions 或 Responses API 的适配类型。
- 模型 ID。
- 超时时间。
- 最大重试次数。

不得假设所有标记为 OpenAI-compatible 的服务都完全兼容。需要允许用户选择接口类型，并在测试连接时验证。

### B. OpenAI 官方接口

作为独立 Provider 实现。

根据所选接口适配请求格式、流式响应和错误结构。

### C. Anthropic

作为独立 Provider 实现。

需要适配其消息格式、工具调用和流式协议。

### D. 本地模型

为本地模型保留接口。

可考虑兼容：

- Ollama。
- LM Studio。
- 其他提供本地 HTTP API 的推理服务。

本地模型不需要 API Key 时允许留空。

### E. 自定义 Provider

允许用户添加其他服务商。

用户填写：

- Endpoint。
- API Key。
- API 协议类型。
- 模型名称。
- 可选请求参数。

第一版不要求支持所有厂商的特殊功能，但必须有清晰的适配器接口。

## 6.3 Provider 配置示例

```json
{
  "id": "provider-001",
  "name": "My AI Provider",
  "type": "openai-compatible",
  "baseUrl": "https://example.com/v1",
  "apiKeyRef": "credential:provider-001",
  "apiMode": "chat-completions",
  "enabled": true,
  "models": [
    {
      "id": "model-a",
      "displayName": "Model A",
      "capabilities": [
        "chat",
        "streaming"
      ]
    }
  ]
}
```

注意：

- JSON 配置中只能保存凭据引用。
- 真实 API Key 必须放在系统密钥存储中。
- API Key 不能写入普通日志。
- API Key 不能返回给前端插件。
- API Key 不能进入插件存储。
- 导出配置时必须自动排除真实密钥。

## 6.4 Provider Adapter

Go 后端统一定义接口：

```go
type AIProvider interface {
    ID() string

    TestConnection(
        ctx context.Context,
    ) error

    ListModels(
        ctx context.Context,
    ) ([]AIModel, error)

    Chat(
        ctx context.Context,
        req *AIRequest,
    ) (*AIResponse, error)

    ChatStream(
        ctx context.Context,
        req *AIRequest,
        handler func(AIStreamEvent) error,
    ) error
}
```

不同服务商分别实现：

```text
ai/
├── provider/
│   ├── openai/
│   ├── anthropic/
│   ├── compatible/
│   └── local/
│
├── gateway/
├── models/
├── credentials/
├── permissions/
├── usage/
└── tests/
```

Provider Adapter 负责协议差异，不负责插件权限判断。

## 6.5 AI Gateway

所有 AI 调用统一经过 AI Gateway。

```text
Plugin / Core
     │
     ▼
AI Gateway
     │
     ├── 检查调用者身份
     ├── 检查 AI 权限
     ├── 检查 Provider 状态
     ├── 检查模型能力
     ├── 检查请求大小
     ├── 检查调用限制
     ├── 获取密钥引用
     ├── 选择 Provider
     ├── 执行请求
     ├── 处理流式响应
     └── 记录必要的用量信息
              │
              ▼
        AI Provider
```

AI Gateway 需要支持：

- 普通对话。
- 流式对话。
- 模型选择。
- 默认模型。
- 按插件指定模型。
- 超时控制。
- 取消请求。
- 并发限制。
- 请求大小限制。
- 失败重试。
- 用量记录。
- 错误归一化。

重试必须区分可重试错误和不可重试错误，不能对所有请求无限重试。

## 6.6 AI 模型能力

每个模型记录其能力：

```ts
interface AIModel {
  id: string;
  providerId: string;
  name: string;

  capabilities: {
    chat: boolean;
    streaming: boolean;
    vision: boolean;
    toolCalling: boolean;
    jsonOutput: boolean;
    embeddings: boolean;
  };

  contextWindow?: number;
}
```

模型能力必须来自实际配置或经过验证的 Provider 信息，不能仅凭模型名称推断。

插件在调用前可以查询模型能力。

如果某个模型不支持某项能力，应返回明确错误。

## 6.7 插件如何调用 AI

插件不能直接调用 Provider。

插件通过受控的 Plugin API 调用 AI Gateway。

例如：

```ts
const result = await oidesk.ai.chat({
  model: "default",
  messages: [
    {
      role: "user",
      content: "解释这段代码的时间复杂度"
    }
  ]
});
```

流式调用：

```ts
const controller = new AbortController();

await oidesk.ai.chatStream(
  {
    model: "default",
    messages: [
      {
        role: "user",
        content: "分析这道题的解法"
      }
    ]
  },
  {
    signal: controller.signal,

    onChunk(chunk) {
      console.log(chunk);
    },

    onFinish() {
      console.log("finished");
    },

    onError(error) {
      console.error(error);
    }
  }
);
```

以上为插件侧 API 设计示例，实际实现需要根据插件通信协议定义可序列化的取消和流式事件。

插件 API 不允许暴露：

```ts
oidesk.getApiKey()
oidesk.getProviderCredential()
oidesk.rawFetchWithCredential()
```

这类接口必须不存在。

## 6.8 AI 权限和用户授权

Manifest 中声明 AI 能力：

```json
{
  "ai": {
    "enabled": true,
    "capabilities": [
      "chat",
      "summarize"
    ],
    "defaultModel": null
  }
}
```

权限规则：

1. 插件声明不代表用户已经授权。
2. 安装时展示 AI 使用说明。
3. 首次调用时可进一步确认。
4. 用户可以关闭插件的 AI 权限。
5. 插件不能更改自己的授权状态。
6. 用户可以限制插件使用的 Provider。
7. 用户可以限制插件使用的模型。
8. 插件不能调用未授权的模型。
9. AI Gateway 必须在后端再次检查权限。

建议允许用户设置：

```text
插件 AI 权限
├── 是否允许 AI 调用
├── 允许使用的 Provider
├── 允许使用的模型
├── 是否允许发送题目内容
├── 是否允许发送提交代码
├── 是否允许发送本地笔记
└── 是否允许使用联网搜索等外部能力
```

未明确授权的敏感数据不得自动发送给第三方 AI。

## 6.9 AI 数据安全

AI 调用可能把用户数据发送到第三方服务商，因此必须明确展示数据流向。

规则：

- 发送前明确告知用户。
- 不得默认上传整个 SQLite 数据库。
- 不得默认上传全部提交代码。
- 不得默认上传全部笔记。
- 插件只应发送完成任务所需的最少内容。
- API Key 不得进入 Prompt。
- 日志默认不记录完整 Prompt 和响应。
- 用户可以清理本地 AI 调用历史。
- 对于本地模型，应明确标注为本地服务，但仍需遵守用户配置的网络和权限规则。

对于用户主动选择发送的代码或题目内容，应在 UI 中展示发送范围。

## 6.10 AI 调用限制

建议提供全局和插件级限制：

- 单次最大输入长度。
- 单次最大输出长度。
- 最大并发请求数。
- 请求超时时间。
- 最大重试次数。
- 每日调用次数上限（可选）。
- 每日 Token 预算（可选）。
- 单个插件调用预算（可选）。

Provider 未提供真实 Token 用量时，不得伪造精确用量。可以显示估算值并明确标记为估算。

## 6.11 AI 请求数据结构

```ts
interface AIRequest {
  model?: string;

  messages: Array<{
    role: "system" | "user" | "assistant";
    content: string;
  }>;

  temperature?: number;
  maxTokens?: number;

  responseFormat?: {
    type: "text" | "json";
  };

  metadata?: {
    task?: string;
  };
}
```

插件不得通过任意字段注入核心权限或伪造调用者身份。

调用者身份由 Plugin Runtime 注入，不信任插件提交的 `pluginId`。

## 6.12 AI 错误类型

统一错误码：

```text
AI_PROVIDER_NOT_FOUND
AI_PROVIDER_DISABLED
AI_MODEL_NOT_FOUND
AI_MODEL_UNSUPPORTED
AI_PERMISSION_DENIED
AI_CREDENTIAL_MISSING
AI_AUTH_FAILED
AI_RATE_LIMITED
AI_TIMEOUT
AI_REQUEST_TOO_LARGE
AI_NETWORK_ERROR
AI_PROVIDER_ERROR
AI_REQUEST_CANCELLED
AI_UNKNOWN_ERROR
```

前端显示用户能理解的错误信息，日志保留必要的技术细节，但不能泄露密钥。

## 6.13 AI 设置页面

设置页面需要包含：

```text
Settings
└── AI
    ├── Providers
    │   ├── Add Provider
    │   ├── Edit Provider
    │   ├── Delete Provider
    │   └── Test Connection
    │
    ├── Models
    │   ├── Model List
    │   ├── Model Capabilities
    │   └── Default Model
    │
    ├── Permissions
    │   ├── Plugin Permissions
    │   └── Data Sharing
    │
    ├── Usage
    │   ├── Request Count
    │   ├── Token Usage
    │   └── Error Statistics
    │
    └── Advanced
        ├── Timeout
        ├── Retry
        └── Concurrency
```

## 6.14 AI 与插件的关系

AI 不应当成为某一个插件的专属能力。

正确设计：

- Core 提供 AI Gateway。
- 插件通过 Plugin API 使用 AI。
- 插件根据自身需求设计 AI 功能。
- 用户统一配置 Provider 和模型。
- 用户统一管理 AI 权限。

例如：

| 插件 | AI 能力 |
|---|---|
| Notepad | 摘要、润色、整理笔记 |
| Snippets | 解释代码、生成注释、代码转换 |
| Stress Tester | 分析错误输出、解释可能的 WA 原因 |
| Plotter | 解释函数表达式、生成绘图表达式 |
| Algorithm Center | 生成算法总结、辅助整理知识点 |
| OJ Manager | 辅助分析提交记录 |
| Contest Planner | 生成训练计划 |

AI 生成的结果必须由用户确认后才能写入核心数据。

不得允许 AI 自动修改用户手动标注的算法标签、删除笔记或覆盖代码模板。

---

# 7. 第一批官方插件

## 7.1 Notepad

功能：

- Markdown 编辑。
- 多文档。
- 自动保存。
- 标签。
- 搜索。
- 代码块高亮。
- LaTeX。
- 导出 Markdown。
- 未保存内容恢复。
- AI 摘要和润色（可选）。

数据保存在：

```text
plugins/org.oidesk.notepad/
```

AI 权限：

- `ai.chat`
- `ai.summarize`

AI 修改笔记必须先展示差异并由用户确认。

## 7.2 Calculator

面向 OI 的计算器。

功能：

- 基本四则运算。
- 大整数。
- 模运算。
- GCD / LCM。
- 扩展欧几里得。
- 快速幂。
- 组合数。
- 进制转换。
- 位运算。
- 溢出辅助。
- 表达式计算。

不能使用不受限制的 `eval` 执行用户表达式。

应使用安全的表达式解析器。

## 7.3 Snippets

管理常用代码模板。

功能：

- C++ 模板。
- Go 模板。
- Shell 模板。
- 自定义语言。
- 分类。
- 标签。
- 搜索。
- 收藏。
- 一键复制。
- 导入导出。
- 代码高亮。
- 多版本管理（后续）。

用户模板必须保存在独立插件数据目录。

卸载插件时默认保留模板数据。

AI 功能：

- 解释模板。
- 生成注释。
- 整理代码。
- 根据用户输入生成模板草稿。

AI 不得未经确认覆盖原模板。

## 7.4 Plotter

功能：

- 函数绘制。
- 多函数绘制。
- 坐标范围设置。
- 散点图。
- 折线图。
- 柱状图。
- CSV 导入。
- PNG / SVG 导出。

表达式必须安全解析，不能执行任意 JavaScript。

## 7.5 Countdown

管理 OI 赛事和重要日期。

支持：

- CSP-J。
- CSP-S。
- NOIP。
- NOI。
- APIO。
- CTSC。
- IOI。
- Codeforces 比赛。
- AtCoder 比赛。
- 用户自定义日期。

每个事件可以记录：

- 开始时间。
- 结束时间。
- 报名截止。
- 时区。
- 提醒时间。
- 备注。

官方比赛日期必须经过可信来源验证，不能假设每年固定。

## 7.6 Stress Tester

对拍工具必须支持：

- main.cpp。
- std.cpp。
- generator.cpp。
- 自定义文件名。
- 自定义编译参数。
- GCC / Clang。
- 随机数据生成。
- 自动循环。
- 输出比较。
- 第一个差异定位。
- 保存错误数据。
- 超时终止。
- 输出大小限制。
- 自定义 checker。
- 运行日志。

目录示例：

```text
workspace/
└── stress-test/
    ├── main.cpp
    ├── std.cpp
    ├── generator.cpp
    ├── data.in
    ├── main.out
    └── std.out
```

必须支持安全地终止超时进程。

不能只杀死父进程而留下编译器或子进程。

不同操作系统的进程管理方式需要分别实现。

AI 功能：

- 分析 WA 的可能原因。
- 比较错误输出和标准输出。
- 解释崩溃信息。
- 辅助构造边界数据。

AI 不应自动修改用户代码或自动认定某个程序一定正确。

---

# 8. Dashboard 与工作区

## 8.1 Dashboard

首页包含：

- 下一场重要比赛倒计时。
- 最近 OJ 提交。
- OJ 同步状态。
- 今日训练统计。
- 最近使用插件。
- 快捷工具。
- 最近笔记。
- 算法掌握概况。

Dashboard 支持：

- 隐藏组件。
- 拖动组件。
- 调整组件大小。
- 保存布局。
- 恢复默认布局。

## 8.2 工作区模式

提供三种模式：

### 普通模式

展示常用功能。

### 训练模式

突出：

- 题目。
- 提交记录。
- 算法分类。
- 代码片段。
- 对拍工具。
- 训练统计。

### 比赛模式

突出：

- 比赛倒计时。
- 计时器。
- 计算器。
- 常用模板。
- 本地笔记。
- 快捷工具。

工作区模式只调整 UI 布局，不得自动改变用户数据或插件权限。

## 8.3 全局交互

支持：

- 全局搜索。
- 快捷键。
- 命令面板。
- 通知中心。
- 右键菜单。
- 最近使用。
- 插件快速启动。

---

# 9. UI / UX 规范

## 9.1 视觉风格

采用现代开发者工具风格。

要求：

- 默认深色主题。
- 支持浅色主题。
- 低噪音。
- 清晰层级。
- 合理留白。
- 统一组件。
- 避免过度动画。
- 保证长时间使用的舒适性。

## 9.2 主导航

```text
OI-Desk
│
├── Dashboard
│
├── OJ
│   ├── Accounts
│   ├── Submissions
│   └── Problems
│
├── Algorithms
│   ├── Categories
│   ├── Problem List
│   └── Statistics
│
├── Plugins
│   ├── Installed
│   ├── Local Install
│   └── Settings
│
├── AI
│   ├── Providers
│   ├── Models
│   └── Usage
│
└── Settings
    ├── General
    ├── Appearance
    ├── Data
    ├── Shortcuts
    └── Advanced
```

## 9.3 插件 UI

插件 UI 应融入主程序，但必须保留隔离。

插件可以读取主题信息，并使用统一 UI 规范。

插件不能修改主程序全局样式。

插件可以使用独立布局，但不能覆盖主程序导航或其他插件内容。

---

# 10. 数据库设计

## 10.1 核心原则

- SQLite。
- Migration 管理结构版本。
- 事务写入。
- UTC 保存时间。
- 显示时转换本地时区。
- 重要数据定期备份。
- 关键字段建立索引。
- 使用外键约束。
- 插件数据与核心数据库隔离。

## 10.2 核心表

### oj_accounts

```sql
CREATE TABLE oj_accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL,
    username TEXT NOT NULL,
    remote_user_id TEXT,
    display_name TEXT,
    avatar_url TEXT,
    auth_status TEXT NOT NULL DEFAULT 'unknown',
    last_sync_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

不保存明文 Cookie 或 API Key。

### problems

```sql
CREATE TABLE problems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL,
    remote_id TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT,
    difficulty TEXT,
    metadata_json TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE(platform, remote_id)
);
```

### submissions

```sql
CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL,
    remote_id TEXT NOT NULL,
    problem_id INTEGER,
    account_id INTEGER,
    submit_time DATETIME,
    language TEXT,
    status TEXT,
    runtime INTEGER,
    memory INTEGER,
    code_length INTEGER,
    metadata_json TEXT,
    created_at DATETIME NOT NULL,
    FOREIGN KEY(problem_id) REFERENCES problems(id),
    FOREIGN KEY(account_id) REFERENCES oj_accounts(id),
    UNIQUE(platform, remote_id)
);
```

### algorithms

```sql
CREATE TABLE algorithms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER,
    name TEXT NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY(parent_id) REFERENCES algorithms(id)
);
```

### problem_algorithms

```sql
CREATE TABLE problem_algorithms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL,
    problem_id INTEGER NOT NULL,
    algorithm_id INTEGER NOT NULL,
    source TEXT NOT NULL,
    confidence REAL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY(problem_id) REFERENCES problems(id),
    FOREIGN KEY(algorithm_id) REFERENCES algorithms(id),
    UNIQUE(platform, problem_id, algorithm_id, source)
);
```

### sync_states

```sql
CREATE TABLE sync_states (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account_id INTEGER NOT NULL,
    sync_type TEXT NOT NULL,
    cursor TEXT,
    last_success_at DATETIME,
    last_error TEXT,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY(account_id) REFERENCES oj_accounts(id)
);
```

### contests

```sql
CREATE TABLE contests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source TEXT,
    remote_id TEXT,
    name TEXT NOT NULL,
    start_time DATETIME,
    end_time DATETIME,
    registration_deadline DATETIME,
    timezone TEXT,
    url TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

### plugin_registry

```sql
CREATE TABLE plugin_registry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plugin_id TEXT NOT NULL UNIQUE,
    version TEXT NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 0,
    manifest_json TEXT NOT NULL,
    installed_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

### plugin_permissions

```sql
CREATE TABLE plugin_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plugin_id TEXT NOT NULL,
    permission TEXT NOT NULL,
    granted INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME NOT NULL,
    UNIQUE(plugin_id, permission)
);
```

### ai_providers

```sql
CREATE TABLE ai_providers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    provider_type TEXT NOT NULL,
    base_url TEXT NOT NULL,
    api_mode TEXT NOT NULL,
    credential_ref TEXT,
    enabled INTEGER NOT NULL DEFAULT 1,
    config_json TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

### ai_models

```sql
CREATE TABLE ai_models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    provider_id TEXT NOT NULL,
    model_id TEXT NOT NULL,
    display_name TEXT,
    capabilities_json TEXT,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    FOREIGN KEY(provider_id) REFERENCES ai_providers(id),
    UNIQUE(provider_id, model_id)
);
```

### ai_usage

```sql
CREATE TABLE ai_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    provider_id TEXT,
    model_id TEXT,
    caller_type TEXT NOT NULL,
    caller_id TEXT,
    request_status TEXT NOT NULL,
    input_tokens INTEGER,
    output_tokens INTEGER,
    estimated_tokens INTEGER,
    duration_ms INTEGER,
    created_at DATETIME NOT NULL
);
```

AI 调用记录不得默认保存完整 Prompt、响应或密钥。

### app_settings

保存非敏感的应用设置，例如主题、默认模型引用、布局和同步间隔。

敏感凭据不进入该表。

## 10.3 数据关系

```text
OJ Account
    │
    ├── Submissions
    │       │
    │       └── Problems
    │               │
    │               └── Problem Algorithms
    │                         │
    │                         └── Algorithms
    │
    └── Sync States

Plugin Registry
    │
    └── Plugin Permissions

AI Providers
    │
    └── AI Models
          │
          └── AI Usage
```

---

# 11. 项目目录结构

```text
OI-Desk/
│
├── PROJECT_SPEC.md
├── AGENTS.md
├── README.md
├── CHANGELOG.md
│
├── docs/
│   ├── ARCHITECTURE.md
│   ├── DATABASE.md
│   ├── PLUGIN_API.md
│   ├── PLUGIN_SECURITY.md
│   ├── AI_GATEWAY.md
│   ├── AI_PROVIDER.md
│   ├── OJ_API.md
│   ├── DEVELOPMENT.md
│   └── ROADMAP.md
│
├── backend/
│   ├── app/
│   ├── core/
│   ├── database/
│   ├── migrations/
│   ├── oj/
│   ├── algorithms/
│   ├── plugins/
│   ├── ai/
│   │   ├── gateway/
│   │   ├── provider/
│   │   ├── credentials/
│   │   ├── permissions/
│   │   └── usage/
│   ├── contests/
│   ├── backup/
│   ├── notifications/
│   └── tests/
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/
│   │   ├── services/
│   │   ├── plugins/
│   │   ├── ai/
│   │   ├── router/
│   │   └── styles/
│   └── tests/
│
├── plugins/
│   ├── official/
│   │   ├── notepad/
│   │   ├── calculator/
│   │   ├── snippets/
│   │   ├── plotter/
│   │   ├── countdown/
│   │   └── stress-tester/
│   └── examples/
│
├── runner/
│   ├── cpp/
│   └── tests/
│
├── scripts/
│
└── build/
```

---

# 12. 开发阶段与实施流程

禁止一开始就同时开发所有功能。

必须分阶段完成，每阶段通过验收后再进入下一阶段。

## Phase 0：技术验证与架构确认

目标：

- 创建仓库。
- 确定 Wails 版本。
- 验证 Go、Vue、TypeScript、SQLite。
- 验证 Wails 前后端通信。
- 验证 iframe 插件隔离。
- 验证 MessageChannel。
- 验证系统密钥存储。
- 验证 AI Provider 的最小请求。
- 验证流式响应。
- 验证 C++ 独立进程运行。

交付：

- 技术选型报告。
- 架构图。
- 数据库 ER 图。
- 插件通信设计。
- AI Gateway 设计。
- 风险清单。

此阶段不实现完整业务功能。

## Phase 1：基础框架

目标：

- Wails 应用。
- Vue 主界面。
- Go 后端。
- SQLite。
- Migration。
- 日志。
- 配置管理。
- 主题系统。
- 基础导航。

验收：

- 应用可正常启动。
- 数据库可创建和迁移。
- 基础设置可保存。
- Windows、macOS、Linux 的基础构建方案明确。

## Phase 2：插件系统

目标：

- Manifest 解析。
- 插件安装。
- 插件卸载。
- 插件启用和禁用。
- iframe 隔离。
- MessageChannel 通信。
- 插件权限。
- 插件独立存储。
- 插件错误隔离。

先实现一个 Hello World 插件。

验收：

- 插件可以安装和启动。
- 插件不能直接访问核心数据库。
- 未授权 API 调用会被拒绝。
- 禁用插件不会删除数据。
- 插件崩溃不会导致主程序退出。

## Phase 3：AI Gateway

目标：

- Provider 管理。
- API Key 安全保存。
- OpenAI-compatible Provider。
- 模型配置。
- 默认模型。
- 测试连接。
- 普通对话。
- 流式对话。
- 请求取消。
- 插件 AI 权限。
- 调用日志和用量统计。

验收：

- 可以配置一个 OpenAI-compatible 服务。
- 可以配置多个 Provider。
- API Key 不会暴露给前端插件。
- 插件只能调用获授权的 AI 能力。
- 流式请求可以正常结束和取消。
- 网络错误、认证错误和限流错误能够正确区分。

## Phase 4：OJ 数据中心

目标：

- 账号管理。
- OJ Adapter。
- 提交同步。
- 题目同步。
- 增量同步。
- 失败恢复。
- 同步状态展示。

优先实现一个真实 OJ 的完整闭环，再逐步增加其他平台。

验收：

- 可以正确同步提交。
- 不重复插入数据。
- 失败后可继续同步。
- 不同 OJ 数据独立保存。

## Phase 5：算法知识库

目标：

- 算法分类树。
- 题目标签。
- 多对多关联。
- 独立平台统计。
- 掌握状态。
- 搜索和筛选。

验收：

- 一道题可对应多个算法。
- 同名题目不会自动跨 OJ 合并。
- 手动标签不会被 AI 标签覆盖。

## Phase 6：官方插件

按优先级开发：

1. Snippets。
2. Calculator。
3. Notepad。
4. Countdown。
5. Stress Tester。
6. Plotter。

每个插件都必须独立测试。

验收：

- 可安装和卸载。
- 数据能够保留。
- 权限符合 Manifest。
- AI 功能可选。
- 不依赖其他插件才能启动。

## Phase 7：Dashboard 与完善

目标：

- Dashboard。
- 全局搜索。
- 工作区模式。
- 通知中心。
- 快捷键。
- 数据备份。
- 数据恢复。
- UI 优化。
- 性能优化。

验收：

- 常用操作有清晰入口。
- Dashboard 布局可保存。
- 数据备份可恢复。
- 主程序长时间运行稳定。

---

# 13. 测试要求

## 13.1 后端测试

使用 Go testing。

必须测试：

- 数据库 CRUD。
- Migration。
- OJ 去重。
- 增量同步。
- 算法关联。
- 插件权限。
- 插件存储隔离。
- AI Gateway。
- Provider 错误处理。
- 凭据读取。
- 备份恢复。

## 13.2 前端测试

使用 Vitest。

重点测试：

- 状态管理。
- 插件管理 UI。
- AI 配置 UI。
- 错误提示。
- 表单校验。
- 流式响应处理。

## 13.3 E2E 测试

使用 Playwright 或适合 Wails 的测试方式。

测试：

- 应用启动。
- 插件安装。
- 插件运行。
- 权限拒绝。
- AI 配置。
- 数据同步。
- 数据备份。

## 13.4 安全测试

必须覆盖：

- ZIP 路径穿越。
- 插件身份伪造。
- 越权 API 调用。
- 插件间数据访问。
- API Key 泄露。
- XSS。
- iframe 沙箱逃逸风险。
- 恶意 Prompt 和工具调用。
- 超大请求。
- 超时请求。
- 恶意 C++ 子进程。
- 超大输出。
- 进程残留。

---

# 14. 安全与非功能要求

## 14.1 插件安全

插件属于不可信代码。

即使是官方插件，也必须经过相同的权限检查。

不能仅依赖 iframe 沙箱实现安全隔离。需要结合：

- 独立来源或严格的沙箱策略。
- CSP。
- 消息来源校验。
- 后端权限校验。
- API 参数校验。
- 插件存储隔离。
- 网络请求限制。

## 14.2 AI 安全

AI 服务商返回的内容是不可信输入。

不得将模型输出直接作为系统命令执行。

不得让 AI 自动执行未经用户确认的文件删除、代码覆盖、网络请求或其他高风险操作。

AI 的工具调用必须经过核心权限系统。

## 14.3 性能

- 启动速度应尽量快。
- 大量提交记录不能阻塞 UI。
- 同步过程应异步执行。
- 插件加载失败不能影响核心。
- AI 请求不得阻塞主线程。
- 对拍任务不能阻塞主程序。
- 大数据查询必须有分页或限制。

## 14.4 备份

至少支持：

- 核心数据库备份。
- 插件数据备份。
- 用户设置备份。
- AI Provider 配置备份（不包含真实密钥）。
- 数据恢复。

恢复过程应先校验备份完整性。

## 14.5 日志

日志应包含：

- 时间。
- 模块。
- 级别。
- 错误码。
- 调用上下文。

日志不得包含：

- API Key。
- Cookie。
- Session。
- 完整认证头。
- 用户未授权记录的完整 Prompt 和响应。

---

# 15. 暂不实现的功能

第一版不应将以下内容作为核心开发目标：

- 在线 OJ。
- 在线代码编辑器。
- 在线评测。
- 云端账号系统。
- 云端同步。
- 复杂社交系统。
- 插件付费交易。
- 完整插件市场。
- 自动执行 AI 生成代码。
- AI 自动修改用户全部历史数据。
- 强制绑定某个 AI 服务商。

这些功能可以在后续版本评估。

---

# 16. AI 编程 Agent 工作规范

本项目将由 AI 编程 Agent 协助开发。

Agent 必须遵守以下规则。

## 16.1 开始开发前

1. 完整阅读 `PROJECT_SPEC.md`。
2. 阅读 `AGENTS.md`。
3. 检查仓库现状。
4. 检查现有代码。
5. 检查依赖版本。
6. 确认 Wails 版本。
7. 阅读对应版本官方文档。
8. 检查跨平台兼容性。
9. 不得凭空编造 API。

## 16.2 开发过程中

- 每次只实现一个明确模块。
- 先设计接口，再写实现。
- 尽量保持模块低耦合。
- 不得擅自改变数据库核心结构。
- 不得擅自改变插件通信协议。
- 不得擅自降低安全要求。
- 不得为了通过测试而删除测试。
- 不得用假数据冒充真实 OJ 数据。
- 不得将真实密钥写入代码或配置文件。
- 不得把 AI 调用绕过 AI Gateway。
- 不得让插件绕过权限系统。
- 不得在没有说明的情况下引入大型依赖。

## 16.3 每个模块完成后

必须：

1. 格式化代码。
2. 运行相关测试。
3. 检查编译错误。
4. 检查边界情况。
5. 检查安全风险。
6. 更新文档。
7. 汇报修改文件。
8. 汇报测试结果。
9. 汇报尚未解决的问题。

## 16.4 遇到不确定的问题

不得自行编造答案。

应当：

- 查阅官方文档。
- 检查真实接口。
- 编写最小验证程序。
- 记录不确定性。
- 向用户说明影响。
- 给出可验证的替代方案。

---

# 17. 最终交付标准

OI-Desk 第一阶段完成时，应当具备：

- 可运行的跨平台桌面应用。
- 可用的本地数据库。
- 可扩展的 OJ Adapter。
- 至少一个真实 OJ 的同步闭环。
- 独立的算法知识库。
- 可安装、卸载和管理的插件系统。
- 插件权限管理。
- 插件独立数据存储。
- 自定义 AI Provider。
- 统一 AI Gateway。
- 安全的 API Key 管理。
- 插件 AI 权限控制。
- 基础官方插件。
- 本地 C++ 对拍工具。
- Dashboard。
- 数据备份和恢复。
- 完整开发文档。
- 自动化测试。

---

# 18. 项目开发总原则

**先保证核心稳定，再扩展插件；先保证权限安全，再开放 AI；先保证真实数据可靠，再做智能分析。**

OI-Desk 的目标不是把所有功能堆进一个桌面程序，而是建立一个稳定、可扩展、适合 OI 选手长期使用的个人工作台。

任何新功能都应优先考虑：

1. 是否解决真实的高频问题？
2. 是否应该属于 Core？
3. 是否更适合做成插件？
4. 是否需要 AI？
5. 是否会增加不必要的复杂度？
6. 是否会影响现有数据安全？
7. 是否可以独立测试和维护？

如果一个功能不能清晰回答这些问题，就不应急于实现。

**本 SPEC 是产品和架构约束。实际开发中如需修改，必须先说明原因、影响范围和迁移方案，再由项目负责人确认。**