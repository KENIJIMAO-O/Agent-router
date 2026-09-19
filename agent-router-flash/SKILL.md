---
name: agent-router-flash
description: 在处理 agent-router 代码库时使用，涵盖架构、核心数据结构、控制流和修改位置指南（基于 main 分支 f2e23e7a 的快速阅读）
---

# Agent Router 快速指南

## 来源

- 源码：`/Users/kross/Projects/agent-router`
- 引用：`main`
- 提交：`f2e23e7a`
- 生成日期：`2026-09-15`
- 由 flash-code-read 生成（单次扫描，热点抽样）。请参阅“已知缺口”。

## 项目简介

Agent Router 是一个基于 Envoy/Envoy Gateway 的 AI 与智能体流量控制平面。应用客户端可使用兼容 OpenAI、Anthropic、特定服务商或 MCP 的协议；平台运维人员则通过 Kubernetes CRD 或独立运行的 `aigw run` 命令，配置路由、服务商 Schema、凭据、配额、请求变更、故障转移、指标和追踪（`README.md:10`、`README.md:18`、`cmd/aigw/main.go:169`）。

## 架构

```text
AIGatewayRoute / AIServiceBackend / GatewayConfig / MCPRoute
                            |
                            v
                 internal/controller
            HTTPRoute、Secret、状态、工作负载配置
                  |                         |
                  v                         v
      internal/extensionserver        filterapi.Config
             xDS 补丁                       |
                  |               分片 Secret 配置包
                  v                         |
                Envoy                       v
                  |              filterapi 配置包监视器
          +-------+--------+                |
          v                v                +--> RuntimeConfig
  路由级 ext_proc      上游级 ext_proc       +--> MCP 编译后配置
          |              |
          +--- 请求状态 ---+
                               |
             endpointspec -> translator -> 服务商
                               |
                       认证 / 指标 / 追踪
```

- `api/v1beta1` 是当前公开的 Kubernetes API；`api/v1alpha1` 仍用于配额策略，并非无用的兼容代码（`api/v1beta1/registry.go:23`、`internal/controller/quota_policy.go:62`）。
- `internal/controller` 将 CRD 转换为 Gateway API 资源、过滤器配置 Secret、工作负载变更和状态（`internal/controller/controller.go:120`、`internal/controller/ai_gateway_route.go:87`、`internal/controller/gateway.go:89`）。
- `internal/extensionserver` 修改 Envoy Gateway 转换后的 xDS，加入 ext-proc、后端元数据、MCP、配额和 InferencePool 资源（`internal/extensionserver/post_translate_modify.go:113`）。
- `internal/filterapi` 是控制平面与数据平面之间的配置契约。控制器写入分片 YAML 配置包，extproc/MCP 监视器负责验证并加载配置（`internal/controller/filter_config_bundle.go:46`、`internal/filterapi/watcher.go:38`、`internal/filterapi/runtime.go:73`）。
- `internal/extproc` 是 LLM 请求状态机。端点行为由 `internal/endpointspec` 选择、由 `internal/translator` 实现，并绑定到 `internal/apischema` 中的服务商协议类型（`internal/extproc/server.go:128`、`internal/endpointspec/endpointspec.go:62`、`internal/translator/translator.go:47`）。
- MCP 使用自己的 HTTP/JSON-RPC 聚合代理，而非通用 extproc `Processor` 接口（`internal/mcpproxy/mcpproxy.go:61`、`cmd/extproc/mainlib/main.go:353`）。
- `internal/internalapi` 是一个高扇入包，包含跨组件的 Header、元数据、命名和端点前缀约定（`internal/internalapi/internalapi.go:91`、`internal/internalapi/internalapi.go:146`、`internal/internalapi/internalapi.go:256`）。

依赖方向大致为 `cmd -> controller/extensionserver/extproc -> filterapi/endpointspec -> translator/auth/metrics/tracing -> apischema/internalapi`。Go 的 `internal` 规则意味着 `internal/` 下导出的名称属于项目内部扩展接口，而非稳定的第三方 API（`go.mod:1`）。

## 核心数据结构

### Kubernetes 与控制平面类型

- `v1beta1.AIGatewayRoute` 是主要的用户路由资源。其 Spec 包含父级引用、主机名、路由规则和路由级成本定义（`api/v1beta1/ai_gateway_route.go:37`、`api/v1beta1/ai_gateway_route.go:57`）。
- `v1beta1.AIGatewayRouteRule` 组合后端引用、匹配条件、请求/流式超时和规则元数据。同一条规则不能混用 `AIServiceBackend` 和 `InferencePool`，且最多只能包含一个 `InferencePool`（`api/v1beta1/ai_gateway_route.go:209`、`api/v1beta1/ai_gateway_route.go:216`）。
- `v1beta1.AIGatewayRouteRuleBackendRef` 指向 `AIServiceBackend` 或 `InferencePool`，并可覆盖模型、Header、Body、权重和优先级（`api/v1beta1/ai_gateway_route.go:316`）。
- `v1beta1.AIServiceBackendSpec` 必须包含 `schema` 和 `backendRef`。被引用对象必须是 `gateway.envoyproxy.io` 组中的 Envoy Gateway `Backend`；Header 和 JSON Body 变更为可选项（`api/v1beta1/ai_service_backend.go:50`、`api/v1beta1/ai_service_backend.go:68`、`api/v1beta1/ai_service_backend.go:71`）。
- `v1beta1.GatewayConfigSpec` 包含可选的 extproc 设置、全局 LLM 请求成本和正向代理配置。Gateway 通过同一命名空间内的 `aigateway.envoyproxy.io/gateway-config` 注解选择它（`api/v1beta1/gateway_config.go:16`、`api/v1beta1/gateway_config.go:59`、`internal/controller/gateway.go:1485`）。

### 发布配置与运行时配置

- `filterapi.Config` 是序列化后的控制平面/数据平面 DTO，包含版本、UUID、成本、扁平化后端、模型列表、按主机限定的模型和 MCP 配置（`internal/filterapi/filterconfig.go:25`）。
- `filterapi.Backend` 是扁平化的路由/规则/后端引用记录，包含唯一名称、模型覆盖、服务商 Schema、认证以及 Header/Body 变更（`internal/filterapi/filterconfig.go:193`）。其名称由 `PerRouteRuleRefBackendName` 构建（`internal/internalapi/internalapi.go:146`）。
- `filterapi.RuntimeConfig` 将后端切片转换为按名称索引的 Map，构造认证处理器，并在接受配置快照前编译 CEL 成本程序（`internal/filterapi/runtime.go:29`、`internal/filterapi/runtime.go:73`）。
- `filterapi.RuntimeBackend` 将不可变的后端配置绑定到运行时 `BackendAuthHandler`，后者用于变更出站 Header（`internal/filterapi/runtime.go:19`、`internal/filterapi/runtime.go:51`）。

### LLM 请求状态

- `extproc.Server` 持有当前的 `RuntimeConfig`、精确路径处理器注册表，以及内部请求 ID 到路由处理器的关联 Map。关联 Map 由互斥锁保护，但配置指针没有（`internal/extproc/server.go:55`、`internal/extproc/server.go:81`、`internal/extproc/server.go:148`）。
- `routerProcessor` 保存单个下游请求的持久状态：原始类型化/原始 Body、原始模型、流式标志、追踪 Span、重试次数，以及最后绑定的上游处理器（`internal/extproc/processor_impl.go:73`）。
- `upstreamProcessor` 保存单次上游尝试的状态：选中的后端与路由、请求级 Translator、认证和变更处理器、压缩/流式响应状态、Token 用量与指标（`internal/extproc/processor_impl.go:112`、`internal/extproc/processor_impl.go:683`）。
- `translator.Translator[ReqT, SpanT]` 是服务商协议边界，负责转换请求 Body、响应 Header、成功响应 Body 和错误 Body。实例按请求创建，并非线程安全（`internal/translator/translator.go:40`、`internal/translator/translator.go:47`）。
- `metrics.TokenUsage` 通过私有设置位区分“明确上报为零”和“服务商未上报”；流式分块仅覆盖显式提供的字段（`internal/metrics/metrics.go:143`、`internal/metrics/metrics.go:256`）。

### MCP 状态

- `filterapi.MCPConfig`、`MCPRoute` 和 `MCPBackend` 表示 MCP 路由、前缀策略、Header 转发、授权以及工具/提示词选择（`internal/filterapi/mcpconfig.go:9`、`internal/filterapi/mcpconfig.go:37`、`internal/filterapi/mcpconfig.go:62`）。
- `mcpproxy.ProxyConfig` 在验证后整体替换已编译的路由/后端 Map，并在客户端可见的工具表发生变化时发出信号（`internal/mcpproxy/config.go:23`、`internal/mcpproxy/config.go:253`、`internal/mcpproxy/config.go:348`）。
- `mcpproxy.session` 聚合多个后端会话。后端 ID、能力、主体和事件 ID 被编码进加密的复合客户端会话 ID，从而避免使用进程级全局会话存储（`internal/mcpproxy/session.go:50`、`internal/mcpproxy/session.go:507`、`internal/mcpproxy/session.go:706`）。

## 控制流

### LLM 主请求路径

1. `AIGatewayRouteController.Reconcile` 获取 v1beta1 路由，调用 `syncAIGatewayRoute`，并写入 `Accepted` 或 `NotAccepted` 状态（`internal/controller/ai_gateway_route.go:87`、`internal/controller/ai_gateway_route.go:100`）。
2. `syncAIGatewayRoute` 创建两个静态 `HTTPRouteFilter` 对象，创建或更新同名、同命名空间的 `HTTPRoute`，随后触发所引用 Gateway 的同步（`internal/controller/ai_gateway_route.go:117`、`internal/controller/ai_gateway_route.go:186`、`internal/controller/ai_gateway_route.go:224`）。
3. `GatewayController.Reconcile` 收集路由、后端和工作负载状态，并调用 `reconcileFilterConfigSecret`；后者将 Schema、认证、模型、变更和 CEL 成本输出为带版本的配置包（`internal/controller/gateway.go:89`、`internal/controller/gateway.go:389`）。
4. `PostTranslateModify` 向转换后的 xDS 加入后端元数据、路由/上游 extproc 过滤器、UDS 集群，以及可选的 MCP、配额和 InferencePool 资源（`internal/extensionserver/post_translate_modify.go:113`、`internal/extensionserver/post_translate_modify.go:678`）。
5. `mainlib.Main` 启动 extproc、管理端和可选的 MCP 监听器；初始化追踪与指标；注册端点工厂；并启动配置包监视器（`cmd/extproc/mainlib/main.go:208`、`cmd/extproc/mainlib/main.go:318`、`cmd/extproc/mainlib/main.go:348`）。
6. 在路由级 extproc 中，`Server.Process` 按精确的 `:path` 选择处理器，创建 `原始 x-request-id + UUID`，并将路由处理器保存到共享关联 Map（`internal/extproc/server.go:128`、`internal/extproc/server.go:190`、`internal/extproc/server.go:228`）。
7. `routerProcessor.ProcessRequestBody` 通过端点 Spec 解析 JSON 或 Multipart，保留原始类型化/原始输入及模型，注入受保护的模型/路径 Header，启动 Span，并清除 Envoy 路由缓存（`internal/extproc/processor_impl.go:232`、`internal/extproc/processor_impl.go:275`、`internal/extproc/processor_impl.go:314`）。
8. Envoy 选择后端后，上游级 `Server.setBackend` 读取后端元数据，并使用内部 ID 找回路由处理器（`internal/extproc/server.go:375`、`internal/extproc/server.go:387`）。
9. `upstreamProcessor.SetBackend` 绑定后端、认证、变更、模型和新创建的请求级 Translator。每次 Envoy 重试都会创建新的上游处理器，同时保留最初的路由输入（`internal/extproc/processor_impl.go:683`、`internal/extproc/processor_impl.go:708`、`internal/extproc/processor_impl.go:716`）。
10. `upstreamProcessor.ProcessRequestHeaders` 从原始 Body 重新转换，应用 Header/Body 变更、执行认证，并返回 `CONTINUE` 或 `CONTINUE_AND_REPLACE`（`internal/extproc/processor_impl.go:344`、`internal/extproc/processor_impl.go:397`、`internal/extproc/processor_impl.go:419`）。
11. `ProcessResponseHeaders` 为成功的流式响应配置流式处理。`ProcessResponseBody` 负责解压、转换成功/错误 Body、累计 Token 用量、计算 CEL 成本元数据、记录指标，并在流结束时完成 Span（`internal/extproc/processor_impl.go:480`、`internal/extproc/processor_impl.go:530`、`internal/extproc/processor_impl.go:618`、`internal/extproc/processor_impl.go:654`）。
12. 路由级响应处理委托给最后绑定的上游处理器，因此重试时最终使用的后端决定响应转换方式（`internal/extproc/processor_impl.go:169`、`internal/extproc/processor_impl.go:179`）。

### 错误边界

- 未注册的路径会返回 HTTP 404 `ImmediateResponse`，并以 `NotFound` 结束 gRPC 流（`internal/extproc/server.go:204`）。
- 格式错误的客户端 Body 可转为本地 HTTP 400；不受支持但应向用户显示的转换错误会变成 422；缺少请求级凭据会变成 401；其他 Translator/认证错误则作为 extproc gRPC 错误传播（`internal/extproc/processor_impl.go:247`、`internal/extproc/processor_impl.go:365`、`internal/extproc/processor_impl.go:421`）。
- 服务商的非 2xx 响应会经过 `Translator.ResponseError`，将追踪标记为失败，并记录未成功完成的指标（`internal/extproc/processor_impl.go:567`）。
- 初始配置包损坏或缺少分片时会重试；其他初始化错误会阻止启动。运行时更新失败会被记录，之前的配置继续生效（`internal/filterapi/watcher.go:38`、`internal/filterapi/watcher.go:63`、`internal/filterapi/watcher.go:127`）。

### 并发与关闭

- 每条 extproc gRPC 流都有一个处理器。路由流与上游流仅通过互斥锁保护的内部 ID Map 共享请求状态（`internal/extproc/server.go:128`、`internal/extproc/server.go:148`、`internal/extproc/server.go:392`）。
- MCP 后端初始化和逐请求扇出会为每个后端使用一个 Goroutine，并通过 Channel/WaitGroup 汇总（`internal/mcpproxy/mcpproxy.go:216`、`internal/mcpproxy/session.go:298`）。
- extproc 关闭时使用 `GracefulStop`，并通过五秒 Context 关闭健康检查连接、管理服务器、追踪、指标和 MCP 服务器（`cmd/extproc/mainlib/main.go:406`）。
- 旧版 Kubernetes Sidecar 的关闭路径会刻意将取消操作延迟 15 秒（`cmd/extproc/main.go:28`）。

## 各类改动应修改哪里

### 新增 LLM HTTP 端点

1. 在 `internal/apischema/<provider>/` 下定义请求/响应协议类型；OpenAI 的主 Schema 位于 `internal/apischema/openai/openai.go:6`。
2. 基于 `Translator` 添加 Translator 别名或实现（`internal/translator/translator.go:47`、`internal/translator/translator.go:130`）。
3. 添加端点 `Spec` 及其 `GetTranslator` 分支（`internal/endpointspec/endpointspec.go:62`、`internal/endpointspec/endpointspec.go:163`）。
4. 在 `cmd/extproc/mainlib/main.go:292` 和 `cmd/extproc/mainlib/main.go:318` 注册精确路径及指标/追踪工厂。

### 为现有端点新增服务商

1. 在 `internal/translator/<input>_<provider>.go` 中实现该端点的 `Translator` 契约（`internal/translator/translator.go:47`）。
2. 在该端点的 `GetTranslator` 中添加服务商分支；Chat Completions 根据 `filterapi.VersionedAPISchema.Name` 分派（`internal/endpointspec/endpointspec.go:164`）。
3. 如果服务商引入新的 Schema 标识，则扩展 `VersionedAPISchema`/`APISchemaName`（`internal/filterapi/filterconfig.go:133`、`internal/filterapi/filterconfig.go:159`）。
4. 仅在服务商 Translator 需要时传入 `schema.Version` 或 `schema.Prefix`；Azure OpenAI 和 AWS Anthropic 当前会使用 `Version`（`internal/endpointspec/endpointspec.go:172`、`internal/endpointspec/endpointspec.go:174`）。

### 新增后端认证类型

1. 扩展公开的安全策略 CRD（`api/v1beta1/backendsecurity_policy.go:56`）。
2. 在 `internal/controller/gateway.go:861` 将其映射到运行时配置。
3. 扩展 `filterapi.BackendAuth`（`internal/filterapi/filterconfig.go:211`）。
4. 实现 `BackendAuthHandler.Do`（`internal/filterapi/runtime.go:19`）。
5. 在 `backendauth.NewHandler` 中注册构造逻辑（`internal/backendauth/auth.go:18`）。

### 新增后端引用类型

1. 扩展 `AIGatewayRouteRuleBackendRef` 及其 CEL 校验（`api/v1beta1/ai_gateway_route.go:209`、`api/v1beta1/ai_gateway_route.go:316`）。
2. 扩展 `newHTTPRoute` 中的 InferencePool/AIServiceBackend 分支，并构造新的 Gateway API `BackendObjectReference`（`internal/controller/ai_gateway_route.go:257`）。
3. 在读取或输出目标之前，通过 `referenceGrantValidator` 验证跨命名空间使用（`internal/controller/ai_gateway_route.go:265`、`internal/controller/referencegrant.go:81`）。
4. 如果该目标参与数据平面转换，则更新过滤器配置的扁平化逻辑（`internal/controller/gateway.go:389`）。

### 修改未匹配路由的行为

生成的规则列表以 `route-not-found` 结尾；该规则匹配 `rootPrefix`，并引用返回 404 的直接响应过滤器（`internal/controller/ai_gateway_route.go:139`、`internal/controller/ai_gateway_route.go:331`）。如果要重命名或替换此行为，需要同时修改生成逻辑和 `api/v1beta1/ai_gateway_route.go:213` 中预留的 API 规则名称。

### 为端点添加 Multipart 支持

实现 `Spec.ParseMultipartBody(body, contentType, costConfigured)`，返回模型、解析后的请求、流式标志、独立的变更后 Body 和错误（`internal/endpointspec/endpointspec.go:87`）。不支持的端点返回 `errMultipartNotSupported`（`internal/endpointspec/endpointspec.go:126`）；可参考转录端点的具体实现（`internal/endpointspec/endpointspec.go:883`）。

### 修改控制平面/数据平面配置格式

1. 修改 `internal/filterapi/filterconfig.go:25` 中的 DTO。
2. 修改 `internal/controller/gateway.go:389` 中的控制器构造逻辑。
3. 如果序列化发生变化，更新配置包的写入器和读取器（`internal/controller/filter_config_bundle.go:46`、`internal/filterapi/config_bundle.go:66`）。
4. 更新 `internal/filterapi/runtime.go:73` 中的运行时派生逻辑。
5. 保留或有意迁移 `internal/filterapi/filterconfig.go:27` 和 `internal/filterapi/watcher.go:119` 中的版本兼容性检查。

### 修改成本计算

- Gateway 级默认值来自 `GatewayConfigSpec.GlobalLLMRequestCosts`（`api/v1beta1/gateway_config.go:64`）。
- 路由级定义由 Gateway 控制器连同路由名称一起发布（`internal/controller/gateway.go:550`）。
- 运行时计算会先应用匹配的路由级成本，再以全局值补齐缺失的元数据键；若两者都不存在，则不会为该键输出值（`internal/extproc/processor_impl.go:928`、`internal/extproc/processor_impl.go:943`、`internal/extproc/processor_impl.go:964`）。
- 只要配置了任意成本，流式 Chat 请求就会强制设置 `StreamOptions.IncludeUsage=true`，并重写原始 JSON 中的 `stream_options.include_usage`（`internal/endpointspec/endpointspec.go:138`、`internal/endpointspec/endpointspec.go:141`）。

### 新增或修改 MCP 行为

1. 在 `api/v1beta1/mcp_route.go:69` 修改公开资源结构。
2. 在 `internal/controller/mcp_route.go:209` 生成路由与凭据资源。
3. 更新 `internal/filterapi/mcpconfig.go:9` 中的数据平面 DTO。
4. 在 `internal/mcpproxy/mcpproxy.go:158` 和 `internal/mcpproxy/session.go:69` 修改请求/会话行为。
5. 必要时在 `internal/extensionserver/mcproute.go:40` 添加 Envoy 侧监听器/路由支持。

### 新增由控制器管理的 CRD

1. 在 `api/v1beta1` 下定义并注册资源；Scheme 装配从 `internal/controller/controller.go:62` 开始。
2. 实现 `Reconcile`；`internal/controller/ai_service_backend.go:45` 是一个简洁示例。
3. 在 `StartControllers` 中注册 Builder、Watch 和所有权（`internal/controller/controller.go:120`）。
4. 如需反向查找，通过 `ApplyIndexing` 添加字段索引（`internal/controller/controller.go:336`）。

## 注意事项

- `Server.LoadConfig` 和 MCP `ProxyConfig.LoadConfig` 会在没有同步保护的情况下替换实时配置指针；源码注释明确指出这两处更新存在竞态（`internal/extproc/server.go:86`、`internal/mcpproxy/config.go:349`）。如要修复，优先考虑不可变快照配合 `atomic.Pointer`。
- `ReferenceGrant` 验证会比较目标 Group/Kind，但不会检查授权中可选的目标 `name`；因此，仅限特定名称的授权可能会被视为允许该 Kind 下的所有对象。这是基于 `internal/controller/referencegrant.go:162` 的推断，修改前应对照 Gateway API 语义和测试进行核实。
- 过滤器配置包最多包含八个 700 KiB 分片；计入索引开销前约为 5.47 MiB（`internal/controller/filter_config_bundle.go:24`、`internal/controller/filter_config_bundle.go:49`）。
- extproc 处理器会在去除查询参数后按精确路径进行选择，并非模式路由器（`internal/extproc/server.go:100`）。
- 转换过程不得原地修改原始请求 Body，因为重试转换必须保持幂等；共享的 `sjsonOptions` 设置了 `ReplaceInPlace=false`（`internal/translator/translator.go:159`）。
- 上游 `ProcessRequestBody` 被视为不可达并会触发 Panic，因为上游转换已经在请求 Header 阶段使用了 `CONTINUE_AND_REPLACE`（`internal/extproc/processor_impl.go:474`）。
- 后端变更默认值先于路由/后端引用变更应用，因此冲突时路由级值优先；删除集合则取并集（`internal/controller/gateway.go:281`、`internal/controller/gateway.go:328`）。
- 如果配置快照的版本与正在运行的 extproc 不同，快照会被拒绝，旧快照继续生效；滚动升级时需注意这一点（`internal/filterapi/watcher.go:119`）。
- `api/v1alpha1` 仍用于配额处理，不能误认为已被完全取代而删除（`internal/controller/quota_policy.go:62`、`internal/extensionserver/quota_ratelimit.go:125`）。
- MCP 扇出在代理层没有显式设置 `http.Client` 超时；完成时间取决于请求 Context 和 Envoy 层超时（`internal/mcpproxy/mcpproxy.go:68`、`internal/mcpproxy/session.go:305`）。

## 构建、运行与测试

- 独立模式快速启动：`OPENAI_API_KEY=sk-your-key aigw run`（`README.md:46`）。
- 构建二进制文件：`make build.aigw`、`make build.controller`、`make build.extproc`（`Makefile:291`）。
- 单元测试：`make test`（`Makefile:184`）。
- 数据平面集成测试：`make test-data-plane` 和 `make test-data-plane-mcp`（`Makefile:210`、`Makefile:222`）。
- 控制器集成测试：`make test-controller`（`Makefile:233`）。
- CLI 端到端测试：`make test-e2e-aigw`（`Makefile:269`）。
- 提交前必需检查：`make precommit test`（`CONTRIBUTING.md:17`、`Makefile:44`）。

这些命令通过静态分析确认。本次快速阅读期间未执行这些命令，因为 Go 1.27.1 工具链曾尝试在沙箱允许范围之外创建模块缓存下载锁（`go.mod:4`）。

## 已知缺口

以下领域未在本次快速阅读中得到充分确认。处理相关工作时，请直接阅读源码和测试。

- **服务商字段映射和流式边界情况**——已追踪 Translator 分派，但未逐一审查每个服务商的请求、响应和分块转换。请从 `internal/endpointspec/endpointspec.go:62` 及 `internal/translator/` 下选定的文件开始。
- **MCP 授权与重连细节**——已追踪会话/扇出结构，但未全面分析 CEL 授权、SSE 重放以及旧版/新版协议行为。请从 `internal/mcpproxy/authorization.go`、`internal/mcpproxy/legacy.go:186` 和 `internal/mcpproxy/session.go:194` 开始。
- **v1alpha1/v1beta1 兼容策略**——已确认配额功能仍在使用旧版 API，但未系统审查转换和弃用行为。请从 `api/v1alpha1/` 和 `internal/controller/quota_policy.go:62` 开始。
- **测试和生成的客户端**——按照 flash-code-read 规则未纳入分析。生产代码中的发现不代表已有测试覆盖。请从 `tests/`、`*_test.go` 和 `api/v1beta1/client/` 开始。
- **打包与交付**——Helm 模板、镜像/发布自动化和文档站点不在热点采样范围内。请从 `manifests/charts/`、`.github/workflows/` 和 `site/` 开始。
- **构建验证**——当前沙箱内未完成编译与测试；以上命令来自仓库声明，并非已观察到执行通过。

自检：`0/8`；一次针对性修订解决了全部八项考核文档缺口。根据 flash-code-read 规则，未重新执行考核。使用了七个可用子智能体槽位中的六个。
