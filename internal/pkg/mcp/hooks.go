package mcp

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/log"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/opentelemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func RegisterHooks(otelTracer trace.Tracer, serverName string) *server.Hooks {
	hooks := &server.Hooks{}

	var listToolsSpan trace.Span
	hooks.AddBeforeListTools(func(ctx context.Context, _ any, _ *mcp.ListToolsRequest) {
		_, listToolsSpan = otelTracer.Start(ctx, opentelemetry.McpListToolsKey)
		listToolsSpan.SetAttributes(
			attribute.String(opentelemetry.McpServerNameKey, serverName),
		)
	})

	hooks.AddAfterListTools(func(ctx context.Context, _ any, _ *mcp.ListToolsRequest, _ *mcp.ListToolsResult) {
		if listToolsSpan != nil && listToolsSpan.IsRecording() {
			listToolsSpan.AddEvent(opentelemetry.McpListToolsKey + " completed")
			listToolsSpan.End()
		} else {
			log.WithContext(ctx).Warn("Span not recording for ListTools completion")
		}
	})

	var callToolSpan trace.Span
	hooks.AddBeforeCallTool(func(ctx context.Context, _ any, toolReq *mcp.CallToolRequest) {
		_, callToolSpan = otelTracer.Start(ctx, opentelemetry.McpToolCallKey)

		callToolSpan.SetAttributes(
			attribute.String(opentelemetry.ToolNameKey, toolReq.Method),
		)
	})

	hooks.AddAfterCallTool(func(ctx context.Context, _ any, _ *mcp.CallToolRequest, _ *mcp.CallToolResult) {
		if callToolSpan != nil && callToolSpan.IsRecording() {
			callToolSpan.AddEvent(opentelemetry.McpToolCallKey + " completed")
			callToolSpan.End()
		} else {
			log.WithContext(ctx).Warn("Span not recording for CallTool completion")
		}
	})

	return hooks
}
