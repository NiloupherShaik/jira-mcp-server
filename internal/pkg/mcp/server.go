package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/mcp/options"
	"go.opentelemetry.io/otel/trace"
)

func NewServer(
	cfg Config,
	toolList []server.ServerTool,
	otelTracer trace.Tracer,
) *server.StreamableHTTPServer {

	s := server.NewMCPServer(
		cfg.Name,
		"",
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(RegisterHooks(otelTracer, cfg.Name)),
	)

	for _, tool := range toolList {
		s.AddTool(tool.Tool, tool.Handler)
	}

	srv := server.NewStreamableHTTPServer(s, server.WithHTTPContextFunc(options.HeadersToContext()))

	return srv
}
