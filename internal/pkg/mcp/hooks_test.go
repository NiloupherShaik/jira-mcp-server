package mcp

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type HooksTestSuite struct {
	suite.Suite
	tracer     *sdkTrace.TracerProvider
	exporter   *tracetest.InMemoryExporter
	testTracer trace.Tracer
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}

func (s *HooksTestSuite) SetupTest() {
	s.exporter = tracetest.NewInMemoryExporter()
	s.tracer = sdkTrace.NewTracerProvider(
		sdkTrace.WithSyncer(s.exporter),
	)
	otel.SetTracerProvider(s.tracer)
	s.testTracer = s.tracer.Tracer("test-tracer")
}

func (s *HooksTestSuite) TearDownTest() {
	s.exporter.Reset()
}

func (s *HooksTestSuite) TestRegisterHooksReturnsHooksObject() {
	hooks := RegisterHooks(s.testTracer, "test-server")

	s.NotNil(hooks)
	s.Len(hooks.OnBeforeListTools, 1)
	s.Len(hooks.OnAfterListTools, 1)
	s.Len(hooks.OnBeforeCallTool, 1)
	s.Len(hooks.OnAfterCallTool, 1)
}

func (s *HooksTestSuite) TestListToolsHooksCreatesSpansWithCorrectAttributes() {
	serverName := "test-server"
	hooks := RegisterHooks(s.testTracer, serverName)

	ctx := context.Background()
	req := &mcp.ListToolsRequest{}
	result := &mcp.ListToolsResult{}

	beforeHook := hooks.OnBeforeListTools[0]
	beforeHook(ctx, nil, req)

	ctx, span := s.testTracer.Start(ctx, "test-span")

	afterHook := hooks.OnAfterListTools[0]
	afterHook(ctx, nil, req, result)

	span.End()

	s.tracer.ForceFlush(ctx)

	spans := s.exporter.GetSpans()
	s.GreaterOrEqual(len(spans), 1)

	var testSpan *tracetest.SpanStub
	for i := range spans {
		if spans[i].Name == "test-span" {
			testSpan = &spans[i]
			break
		}
	}
	s.NotNil(testSpan, "Test span should exist")
}

func (s *HooksTestSuite) TestCallToolHooksCreatesSpansWithCorrectAttributes() {
	serverName := "test-server"
	hooks := RegisterHooks(s.testTracer, serverName)

	ctx := context.Background()
	toolName := "test-tool"
	req := &mcp.CallToolRequest{}
	req.Params.Name = toolName
	result := &mcp.CallToolResult{}

	beforeHook := hooks.OnBeforeCallTool[0]
	beforeHook(ctx, nil, req)

	ctx, span := s.testTracer.Start(ctx, "test-span")

	afterHook := hooks.OnAfterCallTool[0]
	afterHook(ctx, nil, req, result)

	span.End()

	s.tracer.ForceFlush(ctx)

	spans := s.exporter.GetSpans()
	s.GreaterOrEqual(len(spans), 1)

	var testSpan *tracetest.SpanStub
	for i := range spans {
		if spans[i].Name == "test-span" {
			testSpan = &spans[i]
			break
		}
	}

	s.NotNil(testSpan)
}
