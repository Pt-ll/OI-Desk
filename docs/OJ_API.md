# OJ 适配器与同步设计

## 1. 对现有接口的审查

SPEC 中的 `CheckAuth`、`GetUserInfo`、`GetSubmissions(cursor)` 是正确的最小方向，但不足以稳定覆盖 Luogu、QOJ、Codeforces、UOJ、AtCoder：不同平台的认证方式、分页形式、提交详情、题目标签、限流和错误语义不同；有的平台没有统一公开 API。接口需要能力声明、账号上下文、可恢复游标、分页元数据和分类错误。

## 2. 目标接口

```go
type Platform interface {
    ID() string
    Name() string
    Capabilities() Capabilities
    CheckAuth(ctx context.Context, account AccountRef) (AuthState, error)
    GetUserInfo(ctx context.Context, account AccountRef) (UserInfo, error)
    ListSubmissions(ctx context.Context, account AccountRef, q SubmissionQuery) (SubmissionPage, error)
    GetProblem(ctx context.Context, account AccountRef, ref ProblemRef) (Problem, error)
    Close() error
}

type SubmissionQuery struct {
    Cursor       string
    Since        *time.Time
    PageSize     int
    RequestToken string
}

type SubmissionPage struct {
    Items         []Submission
    NextCursor    string
    Exhausted     bool
    HighWatermark string
    RetryAfter    time.Duration
}

type Capabilities struct {
    AuthModes       []string
    SubmissionList  bool
    ProblemDetails  bool
    ProblemTags     bool
    ContestList     bool
    MaxPageSize     int
    RequiresCookie  bool
}
```

适配器不得把秘密写入日志；`AccountRef` 只包含账号 ID，运行时通过 Credential Manager 取得短生命周期凭据。所有平台远程 ID 转成字符串保存，保留原始 JSON 供诊断但执行大小限制。

## 3. 同步状态机

```text
Idle -> Running -> Succeeded
             ├-> AuthRequired
             ├-> RateLimited (retry_at)
             ├-> TransientFailed (backoff)
             └-> PermanentFailed
```

每个 `(account_id, sync_type)` 只有一个状态记录和一个活动租约。写入顺序为：获取页、校验页完整性、在事务中 upsert 题目/提交、更新游标和高水位；网络错误发生时不推进游标，也不能被解释成空页。请求带幂等 token，重复执行不会重复插入。

## 4. 五个平台能力矩阵

| 平台 | 认证/数据事实 | 适配器策略 | 主要风险 |
|---|---|---|---|
| Luogu | 公开页面和接口形态可能变化，部分数据依赖登录会话 | 先做真实页面/API 探针，按官方/实际可用协议实现 | 反爬、字段变更、Cookie 失效 |
| QOJ | 站点部署和 API 能力可能因实例不同 | 把 base URL 和协议能力配置化；没有公开能力则标记 unsupported | 实例差异、认证方式差异 |
| Codeforces | 提交 ID、用户信息和竞赛接口相对结构化 | 优先官方 API；保留 API key/secret 的独立凭据引用 | API 限流、签名和时间偏差 |
| UOJ | 站点接口和登录会话依赖部署 | 先验证目标实例的公开接口，再实现 cookie/session adapter | 实例版本差异、登录流程变化 |
| AtCoder | 提交/题目页和比赛数据存在公开页面但不等同于统一 API | 以页面/已验证接口为准，分页和解析隔离在 adapter 内 | 语言/状态字段本地化、访问频率限制 |

不能凭空声明任何平台存在公开 API。Phase 4 先选择一个已完成真实探针的平台闭环，再扩展其他 adapter；每个平台都必须有 fixture、网络错误、限流、认证失效和字段缺失测试。

## 5. 归一化规则

- `platform + remote_id` 是题目和提交的幂等键。
- 平台未提供的字段保持 NULL；不以默认值伪造。
- 原始状态保存为受控枚举，未知值进入 `unknown` 并保留原文本。
- 时间转换为 UTC；无法解析的时间导致该条记录进入可审计的解析错误队列。
- OJ 标签写入 `source=oj_tag`；AI 标签写入 `source=ai`，不能覆盖 `manual`。

