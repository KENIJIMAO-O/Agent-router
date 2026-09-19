// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package extproc

import (
	"context"
	"log/slog"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	extprocv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

	"github.com/envoyproxy/ai-gateway/internal/filterapi"
)

// ProcessorFactory is the factory function used to create new instances of a processor.

// kross: 它不是普通函数，而是一个函数类型, 只要一个函数接收这些参数，并返回 (Processor, error)，它就是一个 ProcessorFactory。
// kross: 因为每个用户请求都需要创建自己的 Processor，保存请求级状态
type ProcessorFactory func(_ *filterapi.RuntimeConfig, _ map[string]string, _ *slog.Logger, isUpstreamFilter bool, enableRedaction bool) (Processor, error)

// Processor is the interface for the processor which corresponds to a single gRPC stream per the external processor filter.
// This decouples the processor implementation detail from the server implementation.
//
// This can be either a router filter level processor or an upstream filter level processor.
type Processor interface {
	// ProcessRequestHeaders processes the request headers message.
	ProcessRequestHeaders(context.Context, *corev3.HeaderMap) (*extprocv3.ProcessingResponse, error)
	// ProcessRequestBody processes the request body message.
	ProcessRequestBody(context.Context, *extprocv3.HttpBody) (*extprocv3.ProcessingResponse, error)
	// ProcessResponseHeaders processes the response headers message.
	ProcessResponseHeaders(context.Context, *corev3.HeaderMap) (*extprocv3.ProcessingResponse, error)
	// ProcessResponseBody processes the response body message.
	ProcessResponseBody(context.Context, *extprocv3.HttpBody) (*extprocv3.ProcessingResponse, error)
	// SetBackend instructs the processor to set the backend to use for the request. This is only called
	// when the processor is used in the upstream filter.
	//
	// routerProcessor is the processor that is the "parent" which was used to determine the route at the
	// router level. It holds the additional state that can be used to determine the backend to use.
	SetBackend(ctx context.Context, backend *filterapi.RuntimeBackend, routeName string, routerProcessor Processor) error
}

// passThroughProcessor implements the Processor interface.
type passThroughProcessor struct{}

// ProcessRequestHeaders implements [Processor.ProcessRequestHeaders].
// kross: 处理请求 Header
//
//	kross: &corev3.HeaderMap{
//	    Headers: []*corev3.HeaderValue{
//	        {Key: ":method", Value: "POST"},
//	        {Key: ":path", Value: "/v1/chat/completions"},
//	        {Key: "content-type", Value: "application/json"},
//	        {Key: "authorization", Value: "Bearer ..."},
//	    },
//	},这个是输入，这个是输出 *extprocv3.ProcessingResponse
func (p passThroughProcessor) ProcessRequestHeaders(context.Context, *corev3.HeaderMap) (*extprocv3.ProcessingResponse, error) {
	return &extprocv3.ProcessingResponse{Response: &extprocv3.ProcessingResponse_RequestHeaders{}}, nil
}

// kross: 处理用户请求 Body
//
//	输入：&extprocv3.HttpBody{
//	    Body: []byte(`{
//	        "model": "gpt-5.6-sol",
//	        "messages": [
//	            {
//	                "role": "user",
//	                "content": "hello"
//	            }
//	        ]
//	    }`),
//	    EndOfStream: true,
//	}
//
// ProcessRequestBody implements [Processor.ProcessRequestBody].
func (p passThroughProcessor) ProcessRequestBody(context.Context, *extprocv3.HttpBody) (*extprocv3.ProcessingResponse, error) {
	return &extprocv3.ProcessingResponse{Response: &extprocv3.ProcessingResponse_RequestBody{}}, nil
}

// ProcessResponseHeaders implements [Processor.ProcessResponseHeaders].
// kross: 理模型服务返回的响应 Header
func (p passThroughProcessor) ProcessResponseHeaders(context.Context, *corev3.HeaderMap) (*extprocv3.ProcessingResponse, error) {
	return &extprocv3.ProcessingResponse{Response: &extprocv3.ProcessingResponse_ResponseHeaders{}}, nil
}

// ProcessResponseBody implements [Processor.ProcessResponseBody].
// kross: 作用是处理模型服务返回的 Body
func (p passThroughProcessor) ProcessResponseBody(context.Context, *extprocv3.HttpBody) (*extprocv3.ProcessingResponse, error) {
	return &extprocv3.ProcessingResponse{Response: &extprocv3.ProcessingResponse_ResponseBody{}}, nil
}

// SetBackend implements [Processor.SetBackend].
// kross: 它在路由已经完成、Backend 已经确定后调用,它将两阶段处理关联起来:
// routerProcessor
//
//	已保存原始 Body、原始模型
//	      │
//	      │ SetBackend
//	      ▼
//
// upstreamProcessor
//
//	获得具体 Backend 配置
func (p passThroughProcessor) SetBackend(context.Context, *filterapi.RuntimeBackend, string, Processor) error {
	return nil
}
