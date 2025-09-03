package mcp_test

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	internalMCP "github.com/nable-fusion/mcp-server-template/internal/pkg/mcp"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"testing"
)

type ServerTest struct {
	suite.Suite

	mockTracer trace.Tracer
	mockConfig internalMCP.Config
	ctx        context.Context
	toolList   []server.ServerTool
}

func (t *ServerTest) SetupTest() {
	t.ctx = context.Background()

	toolDef := mcp.NewTool(
		"test-tool",
		mcp.WithDescription("A test tool"))

	sampleTool := server.ServerTool{
		Tool: toolDef,
		Handler: server.ToolHandlerFunc(func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return &mcp.CallToolResult{}, nil
		}),
	}

	t.toolList = []server.ServerTool{sampleTool}

	t.mockTracer = noop.NewTracerProvider().Tracer("test-tracer")
	t.mockConfig = internalMCP.Config{
		Name: "Test MCP Server",
	}
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTest))
}

func (t *ServerTest) TestNewServerInitialisation() {
	mcpServer := internalMCP.NewServer(t.mockConfig, t.toolList, t.mockTracer)

	t.NotNil(mcpServer)
	t.IsType(&server.StreamableHTTPServer{}, mcpServer, "Expected StreamableHTTPServer type")
	t.Equal(1, len(t.toolList), "Tool list should contain 1 tool")
}
