# Agent Router Symbol Map

Generated from `main` at `f2e23e7a` on 2026-09-15. This is a curated lookup index, not a complete declaration dump. Tests, generated clients, provider schema tables, and trivial wrappers are excluded.

## Entry points at a glance

| Symbol | Location | Purpose |
|---|---|---|
| `func main()` | `cmd/aigw/main.go:169` | standalone CLI process entry |
| `func (*rootCmd) doMain(...)` | `cmd/aigw/main.go:180` | dispatches `run`, version, healthcheck, and download commands |
| `func main()` | `cmd/controller/main.go:411` | Kubernetes controller and extension-server entry |
| `func main()` | `cmd/extproc/main.go:20` | extproc signal/lifecycle entry |
| `func Main(ctx, args, stderr) error` | `cmd/extproc/mainlib/main.go:208` | assembles and serves extproc, MCP, health, metrics, and tracing |
| `func StartControllers(...) error` | `internal/controller/controller.go:120` | registers and starts all reconcilers |
| `func (*Server) PostTranslateModify(...)` | `internal/extensionserver/post_translate_modify.go:113` | main Envoy Gateway xDS patch hook |
| `func (*Server) Process(stream) error` | `internal/extproc/server.go:128` | Envoy ExternalProcessor stream handler |
| `func NewMCPProxy(...)` | `internal/mcpproxy/mcpproxy.go:61` | constructs MCP HTTP proxy and compiled-config receiver |

## Public Kubernetes API — `api/v1beta1`

| Symbol | Location | Purpose |
|---|---|---|
| `type AIGatewayRoute struct` | `api/v1beta1/ai_gateway_route.go:37` | top-level AI routing CRD |
| `type AIGatewayRouteSpec struct` | `api/v1beta1/ai_gateway_route.go:57` | parent, hostname, rules, and costs |
| `type AIGatewayRouteRule struct` | `api/v1beta1/ai_gateway_route.go:216` | one match/backend routing rule |
| `type AIGatewayRouteRuleBackendRef struct` | `api/v1beta1/ai_gateway_route.go:316` | AI backend or InferencePool reference plus overrides |
| `type AIServiceBackend struct` | `api/v1beta1/ai_service_backend.go:28` | provider-aware backend CRD |
| `type AIServiceBackendSpec struct` | `api/v1beta1/ai_service_backend.go:50` | schema, Envoy Backend ref, and mutations |
| `type BackendSecurityPolicy struct` | `api/v1beta1/backendsecurity_policy.go:37` | backend credential policy CRD |
| `type BackendSecurityPolicySpec struct` | `api/v1beta1/backendsecurity_policy.go:56` | credential source and target selection |
| `type GatewayConfig struct` | `api/v1beta1/gateway_config.go:40` | gateway-wide Agent Router configuration |
| `type GatewayConfigSpec struct` | `api/v1beta1/gateway_config.go:60` | extproc, global costs, and forward proxy |
| `type MCPRoute struct` | `api/v1beta1/mcp_route.go:49` | MCP routing/aggregation CRD |
| `type MCPRouteSpec struct` | `api/v1beta1/mcp_route.go:69` | parents, rules, and MCP backends |
| `type VersionedAPISchema struct` | `api/v1beta1/shared_types.go:19` | provider schema name/version/prefix |
| `func AddKnownTypes(*runtime.Scheme) error` | `api/v1beta1/registry.go:42` | registers v1beta1 resources with Kubernetes scheme |

## Controller — `internal/controller`

| Symbol | Location | Purpose |
|---|---|---|
| `type Options struct` | `internal/controller/controller.go:65` | reconciler and deployment configuration |
| `func StartControllers(...) error` | `internal/controller/controller.go:120` | composition root for controllers and watches |
| `func ApplyIndexing(...) error` | `internal/controller/controller.go:336` | installs reverse-lookup field indexes |
| `type AIGatewayRouteController struct` | `internal/controller/ai_gateway_route.go:58` | reconciles AI route resources |
| `func (*AIGatewayRouteController) Reconcile(...)` | `internal/controller/ai_gateway_route.go:87` | synchronizes HTTPRoute/filters/status |
| `func (*AIGatewayRouteController) syncAIGatewayRoute(...)` | `internal/controller/ai_gateway_route.go:161` | creates derived routing resources |
| `func (*AIGatewayRouteController) newHTTPRoute(...)` | `internal/controller/ai_gateway_route.go:244` | translates route rules and backend refs |
| `type GatewayController struct` | `internal/controller/gateway.go:75` | reconciles Gateway-level configuration |
| `func (*GatewayController) Reconcile(...)` | `internal/controller/gateway.go:89` | coordinates workloads and config snapshot |
| `func (*GatewayController) reconcileFilterConfigSecret(...)` | `internal/controller/gateway.go:389` | builds and publishes filterapi bundle |
| `type AIBackendController struct` | `internal/controller/ai_service_backend.go:27` | watches AIServiceBackend changes |
| `func (*AIBackendController) Reconcile(...)` | `internal/controller/ai_service_backend.go:45` | triggers affected AI route reconciliation |
| `type MCPRouteController struct` | `internal/controller/mcp_route.go:42` | reconciles MCPRoute resources |
| `func (*MCPRouteController) Reconcile(...)` | `internal/controller/mcp_route.go:64` | builds MCP proxy/backend routes and status |
| `type BackendSecurityPolicyController struct` | `internal/controller/backend_security_policy.go:57` | reconciles backend authentication policy |
| `type QuotaPolicyController struct` | `internal/controller/quota_policy.go:31` | reconciles live v1alpha1 quota resources |
| `func FilterConfigBundleIndexSecretName(...) string` | `internal/controller/secret_name.go:48` | derives bundle index Secret name |

## Envoy Gateway extension — `internal/extensionserver`

| Symbol | Location | Purpose |
|---|---|---|
| `type Server struct` | `internal/extensionserver/extensionserver.go:28` | Envoy Gateway extension implementation |
| `func New(...) (*Server, error)` | `internal/extensionserver/extensionserver.go:51` | constructs extension server |
| `func (*Server) PostTranslateModify(...)` | `internal/extensionserver/post_translate_modify.go:113` | orchestrates xDS post-translation patches |
| `func (*Server) maybeModifyListenerAndRoutes(...) error` | `internal/extensionserver/post_translate_modify.go:678` | patches AI/InferencePool listeners and routes |
| `func (*Server) insertRouterLevelAIGatewayExtProc(...) error` | `internal/extensionserver/post_translate_modify.go:904` | installs router-level extproc filter |
| `func (*Server) maybeInjectQuotaRateLimiting(...)` | `internal/extensionserver/quota_ratelimit.go:59` | adds quota cluster, descriptors, and filters |
| `func (*Server) insertRequestHeaderToMetadataFilters(...) error` | `internal/extensionserver/header_to_metadata.go:23` | protects internal routing metadata |

## Shared configuration — `internal/filterapi`

| Symbol | Location | Purpose |
|---|---|---|
| `type Config struct` | `internal/filterapi/filterconfig.go:25` | serialized control/data configuration |
| `type Model struct` | `internal/filterapi/filterconfig.go:59` | declared model metadata |
| `type GlobalLLMRequestCost struct` | `internal/filterapi/filterconfig.go:71` | gateway-level cost default |
| `type LLMRequestCost struct` | `internal/filterapi/filterconfig.go:88` | route-scoped cost rule |
| `type VersionedAPISchema struct` | `internal/filterapi/filterconfig.go:133` | runtime provider schema selection |
| `type Backend struct` | `internal/filterapi/filterconfig.go:193` | flattened backend configuration |
| `type BackendAuth struct` | `internal/filterapi/filterconfig.go:211` | runtime authentication union |
| `func UnmarshalConfigYaml(path string) (*Config, error)` | `internal/filterapi/filterconfig.go:398` | decodes standalone YAML configuration |
| `type RuntimeConfig struct` | `internal/filterapi/runtime.go:29` | indexed and compiled config snapshot |
| `type RuntimeBackend struct` | `internal/filterapi/runtime.go:51` | backend plus auth handler |
| `func NewRuntimeConfig(...) (*RuntimeConfig, error)` | `internal/filterapi/runtime.go:73` | validates and compiles a snapshot |
| `type ConfigReceiver interface` | `internal/filterapi/watcher.go:24` | target of bundle hot reload |
| `func StartConfigBundleWatcher(...) error` | `internal/filterapi/watcher.go:38` | initial load and periodic bundle reload |
| `type ConfigBundleIndex struct` | `internal/filterapi/config_bundle.go:24` | version/checksum/part manifest |
| `func ReassembleBundleConfig(...) (*Config, error)` | `internal/filterapi/config_bundle.go:66` | verifies and decodes ordered parts |

## External processor — `internal/extproc`

| Symbol | Location | Purpose |
|---|---|---|
| `type ProcessorFactory func(...)` | `internal/extproc/processor.go:19` | creates a request processor for an endpoint |
| `type Processor interface` | `internal/extproc/processor.go:25` | four-phase extproc contract |
| `type Server struct` | `internal/extproc/server.go:55` | extproc registry, config, and request association state |
| `func NewServer(...) (*Server, error)` | `internal/extproc/server.go:67` | constructs extproc server |
| `func (*Server) LoadConfig(...) error` | `internal/extproc/server.go:81` | compiles and swaps runtime configuration |
| `func (*Server) Register(path, factory)` | `internal/extproc/server.go:91` | registers one exact endpoint path |
| `func (*Server) Process(stream) error` | `internal/extproc/server.go:128` | drives one Envoy extproc stream |
| `func NewFactory[...] (...) ProcessorFactory` | `internal/extproc/processor_impl.go:59` | binds metrics/tracing/endpoint spec |
| `func (*routerProcessor) ProcessRequestBody(...)` | `internal/extproc/processor_impl.go:232` | parses and saves original downstream request |
| `func (*upstreamProcessor) ProcessRequestHeaders(...)` | `internal/extproc/processor_impl.go:344` | translates, mutates, and authenticates upstream request |
| `func (*upstreamProcessor) ProcessResponseHeaders(...)` | `internal/extproc/processor_impl.go:480` | configures provider response processing |
| `func (*upstreamProcessor) ProcessResponseBody(...)` | `internal/extproc/processor_impl.go:530` | translates response and records usage/cost |
| `func (*upstreamProcessor) SetBackend(...)` | `internal/extproc/processor_impl.go:683` | binds selected backend and request-scoped translator |

## Endpoint and translation — `internal/endpointspec`, `internal/translator`

| Symbol | Location | Purpose |
|---|---|---|
| `type Spec[...] interface` | `internal/endpointspec/endpointspec.go:62` | endpoint parsing, translation, and streaming contract |
| `type ChatCompletionsEndpointSpec struct` | `internal/endpointspec/endpointspec.go:119` | OpenAI chat endpoint behavior |
| `func (ChatCompletionsEndpointSpec) ParseBody(...)` | `internal/endpointspec/endpointspec.go:130` | parses chat body and forces usage for costed streams |
| `func (ChatCompletionsEndpointSpec) GetTranslator(...)` | `internal/endpointspec/endpointspec.go:163` | selects chat translator by backend schema |
| `type Translator[ReqT, SpanT] interface` | `internal/translator/translator.go:47` | request/response/error protocol conversion |
| `type ContentTypeSetter interface` | `internal/translator/translator.go:88` | optional multipart content-type hook |
| `type RequestHeadersSetter interface` | `internal/translator/translator.go:94` | optional provider request-header hook |
| `type HeaderValueFilterSetter interface` | `internal/translator/translator.go:104` | optional multi-value header filter hook |
| `type ResponseRedactor interface` | `internal/translator/translator.go:110` | optional safe response logging hook |

## Authentication — `internal/backendauth`

| Symbol | Location | Purpose |
|---|---|---|
| `func NewHandler(...) (BackendAuthHandler, error)` | `internal/backendauth/auth.go:18` | selects concrete auth implementation |
| `var ErrCredentialMissing` | `internal/backendauth/credential_override.go:25` | signals absent required dynamic credentials |
| `func WithEnvoyMetadata(...) context.Context` | `internal/backendauth/credential_override.go:34` | exposes protected Envoy metadata to auth handlers |

## MCP proxy — `internal/mcpproxy`

| Symbol | Location | Purpose |
|---|---|---|
| `func NewMCPProxy(...)` | `internal/mcpproxy/mcpproxy.go:61` | creates MCP proxy config and HTTP mux |
| `type SessionCrypto interface` | `internal/mcpproxy/crypto.go:21` | encrypts/decrypts composite session IDs |
| `func NewPBKDF2AesGcmSessionCrypto(...) SessionCrypto` | `internal/mcpproxy/crypto.go:37` | default session crypto implementation |
| `type FallbackEnabledSessionCrypto struct` | `internal/mcpproxy/crypto.go:98` | supports encryption-key rotation fallback |
| `type MCPConfig struct` | `internal/filterapi/mcpconfig.go:9` | serialized MCP configuration root |
| `type MCPRoute struct` | `internal/filterapi/mcpconfig.go:37` | one compiled MCP route input |
| `type MCPBackend struct` | `internal/filterapi/mcpconfig.go:62` | one upstream MCP backend input |

## Observability — `internal/metrics`, `internal/tracing`

| Symbol | Location | Purpose |
|---|---|---|
| `func NewMeterFromEnv(...)` | `internal/metrics/metrics.go:35` | creates Prometheus/OTel meter and shutdown hook |
| `type Metrics interface` | `internal/metrics/metrics.go:97` | per-request LLM metrics contract |
| `type Factory interface` | `internal/metrics/metrics.go:130` | constructs request metrics instances |
| `type TokenUsage struct` | `internal/metrics/metrics.go:143` | token usage with explicit-set semantics |
| `type MCPMetrics interface` | `internal/metrics/mcp_metrics.go:112` | MCP metrics contract |
| `func NewMCP(...) MCPMetrics` | `internal/metrics/mcp_metrics.go:147` | constructs MCP metrics implementation |

## Internal conventions — `internal/internalapi`

| Symbol | Location | Purpose |
|---|---|---|
| `type Header [2]string` | `internal/internalapi/headers.go:9` | compact header mutation pair |
| `func PerRouteRuleRefBackendName(...) string` | `internal/internalapi/internalapi.go:146` | creates stable flattened backend identity |
| `func ParseRequestHeaderAttributeMapping(...)` | `internal/internalapi/internalapi.go:188` | parses attribute-to-header mapping |
| `func MergeRequestHeaderAttributeMappings(...)` | `internal/internalapi/internalapi.go:222` | merges base and override mappings |
| `type EndpointPrefixes struct` | `internal/internalapi/internalapi.go:256` | configurable public endpoint namespaces |
| `func ParseEndpointPrefixes(...)` | `internal/internalapi/internalapi.go:277` | parses prefixes and fills defaults |
| `func GetUserFacingError(error) error` | `internal/internalapi/user_facing_errors.go:32` | unwraps errors safe to expose to clients |

## Not covered

- Provider wire-schema declarations under `internal/apischema/` are intentionally omitted because several files are large flat schema tables.
- Generated Kubernetes clients, listers, informers, and deepcopy files under `api/*/client/` and `zz_generated.*` are excluded.
- Tests, fixtures, Helm templates, website code, and release automation are outside this symbol index.
- See `SKILL.md` Known Gaps before relying on this map for provider-specific, MCP authorization, compatibility, or packaging changes.
