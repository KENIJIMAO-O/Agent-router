<p align="center">
  <a href="https://theagentrouter.ai">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="site/static/img/brand/ar-horizontal-on-dark.svg">
      <img alt="Agent Router" src="site/static/img/brand/ar-horizontal-primary.svg" width="420">
    </picture>
  </a>
</p>

<p align="center">
  <b>由 Envoy 驱动的 AI 与智能体流量开源控制平面。</b><br>
  Agentic AI Foundation 旗下项目，前身为 <i>Envoy AI Gateway</i>。
</p>

<p align="center">
  <a href="https://theagentrouter.ai/docs">文档</a> ·
  <a href="https://theagentrouter.ai/docs/getting-started/">快速开始</a> ·
  <a href="https://theagentrouter.ai/blog">博客</a> ·
  <a href="https://theagentrouter.ai/release-notes/">发行说明</a> ·
  <a href="https://theagentrouter.ai/talks">演讲</a>
</p>

# Agent Router

Agent Router 为应用团队提供统一且兼容 OpenAI 的 API，可连接从托管服务商到自托管推理服务及 MCP 服务器的各种模型和工具。平台团队可以集中管理凭据、路由、配额、故障转移和用量归因，并由 [Envoy](https://envoyproxy.io) 与 [Envoy Gateway](https://gateway.envoyproxy.io) 负责实施。

**Agent Router 负责控制，Envoy 负责承载。**

## 前身：Envoy AI Gateway

Envoy AI Gateway 现已更名为 **Agent Router**，并成为 Agentic AI Foundation 旗下项目。代码、维护者、发布节奏和 Apache 2.0 许可证均保持不变。你部署的任何资源都不需要重命名：

- **CRD 和 API 组保持不变**：`AIGatewayRoute`、`AIServiceBackend`、`BackendSecurityPolicy` 以及 `aigateway.envoyproxy.io` API 组均维持原样。
- **CLI 仍为 `aigw`**，`envoy-ai-gateway-system` 命名空间也没有变化。
- **容器镜像、Helm Chart 和 Go 模块路径保持不变**（`docker.io/envoyproxy/ai-gateway-*`、`github.com/envoyproxy/ai-gateway`）。
- **代码仓库已迁移**至 [`theagentrouter/agent-router`](https://github.com/theagentrouter/agent-router)；旧的 `envoyproxy/ai-gateway` 链接会自动跳转。
- **网站已迁移**至 [theagentrouter.ai](https://theagentrouter.ai)；所有 `aigateway.envoyproxy.io` 链接都会跳转到对应的新页面。

现有的清单文件可以继续使用。

## 快速开始

只需一条命令，即可在本机以独立模式运行兼容 OpenAI 的 Agent Router：

```shell
OPENAI_API_KEY=sk-your-key aigw run
```

然后将任意兼容 OpenAI 的客户端指向 `http://localhost:1975/v1`。有关安装和服务商自动配置，请参阅 [CLI 指南](https://theagentrouter.ai/docs/cli/)；如需配合 Envoy Gateway 部署到 Kubernetes，请参阅[入门指南](https://theagentrouter.ai/docs/getting-started/)。

## 使用方式

Agent Router 采用两级网关模式：**一级网关（Tier One Gateway）**充当集中式入口，**二级网关（Tier Two Gateway）**负责处理进入自托管模型服务集群的流量。

- **一级网关**负责身份验证、顶层路由和全局限流。
- **二级网关**对自托管模型访问提供精细控制，并通过端点选择器优化 LLM 推理。

![](site/blog/images/aigw-ref.drawio.png)

## 支持的 AI 服务商

Agent Router 支持多种 AI 服务商，便于集成你所使用的 LLM 服务：

<div align="center">
  <table>
    <tr>
      <td align="center" width="120">
        <img src="site/static/img/providers/openai.svg" width="60" height="60" alt="OpenAI"/>
        <br><sub><b>OpenAI</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/azure-openai.svg" width="60" height="60" alt="Azure OpenAI"/>
        <br><sub><b>Azure OpenAI</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/google-gemini.svg" width="60" height="60" alt="Google Gemini"/>
        <br><sub><b>Google Gemini</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/vertex-ai.svg" width="60" height="60" alt="Vertex AI"/>
        <br><sub><b>Vertex AI</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/aws-bedrock.svg" width="60" height="60" alt="AWS Bedrock"/>
        <br><sub><b>AWS Bedrock</b></sub>
      </td>
    </tr>
    <tr>
      <td align="center" width="120">
        <img src="site/static/img/providers/mistral.svg" width="60" height="60" alt="Mistral"/>
        <br><sub><b>Mistral</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/cohere.svg" width="60" height="60" alt="Cohere"/>
        <br><sub><b>Cohere</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/groq.svg" width="60" height="60" alt="Groq"/>
        <br><sub><b>Groq</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/together-ai.svg" width="60" height="60" alt="Together AI"/>
        <br><sub><b>Together AI</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/deepinfra.svg" width="60" height="60" alt="DeepInfra"/>
        <br><sub><b>DeepInfra</b></sub>
      </td>
    </tr>
    <tr>
      <td align="center" width="120">
        <img src="site/static/img/providers/deepseek.svg" width="60" height="60" alt="DeepSeek"/>
        <br><sub><b>DeepSeek</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/hunyuan.svg" width="60" height="60" alt="Hunyuan"/>
        <br><sub><b>Hunyuan</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/sambanova.svg" width="60" height="60" alt="SambaNova"/>
        <br><sub><b>SambaNova</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/grok.svg" width="60" height="60" alt="Grok"/>
        <br><sub><b>Grok</b></sub>
      </td>
      <td align="center" width="120">
        <img src="site/static/img/providers/tars.svg" width="60" height="60" alt="Tetrate Agent Router Service"/>
        <br><sub><b>Tetrate Agent Router Service</b></sub>
      </td>
    </tr>
    <tr>
      <td align="center" width="120">
        <img src="site/static/img/providers/anthropic.svg" width="60" height="60" alt="Anthropic"/>
        <br><sub><b>Anthropic</b></sub>
      </td>
    </tr>
  </table>
</div>

## 文档

- [项目文档](https://theagentrouter.ai/docs)：Agent Router 的完整文档。
- [快速开始](https://theagentrouter.ai/docs/getting-started/)：通过几个简单步骤开始使用 Agent Router。
- [CLI 指南](https://theagentrouter.ai/docs/cli/)：使用 `aigw run` 独立运行 Agent Router。
- [核心概念](https://theagentrouter.ai/docs/concepts/)：了解 Agent Router 的架构与资源。
- [博客](https://theagentrouter.ai/blog)：查看公告、功能介绍和参考架构。
- [演讲与分享](https://theagentrouter.ai/talks)：查看 Agent Router 相关演讲。

## 联系我们

- 在线交流：加入 [Agent Router Discord][Discord]，向社区提问并参与协作。
- 每周例会：维护者和社区成员每周一举行会议。你可以通过[公开会议文档][meeting]报名、添加议题并参与讨论。

## 参与贡献

作为 Agentic AI Foundation 旗下项目，我们遵循 LF Projects 的[行为准则][Code of conduct]。

如需通过 Pull Request 为项目贡献代码，请阅读 [CONTRIBUTING.md](CONTRIBUTING.md)，其中包含项目构建与测试说明。

## 背景

使用 Envoy Gateway 构建[云原生 LLM 网关][Cloud Native LLM Gateway]的提案促成了本项目的诞生。

---

Agent Router 是由 Envoy 驱动的 Agentic AI Foundation 项目。版权所有 © Agent Router a Series of LF Projects, LLC。有关网站使用条款、[商标政策](https://theagentrouter.ai/trademark-policy)及其他项目政策，请访问 <https://lfprojects.org>。

[meeting]: https://docs.google.com/document/d/10e1sfsF-3G3Du5nBHGmLjXw5GVMqqCvFDqp_O65B0_w/edit?tab=t.0
[Discord]: https://discord.gg/xuxtPq43gZ
[Code of conduct]: CODE_OF_CONDUCT.md
[Cloud Native LLM Gateway]: https://docs.google.com/document/d/1FQN_hGhTNeoTgV5Jj16ialzaSiAxC0ozxH1D9ngCVew/edit?tab=t.0#heading=h.uuu99yemq4eo
