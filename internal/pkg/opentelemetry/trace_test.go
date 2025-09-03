package opentelemetry

import (
	"context"
	"github.com/nable-fusion/fusion-cloud-common/pkg/log"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap/zaptest"
	"testing"
)

type OtelTracingTestSuite struct {
	suite.Suite
}

func TestOtelTracingTestSuite(t *testing.T) {
	suite.Run(t, new(OtelTracingTestSuite))
}

func (t *OtelTracingTestSuite) SetupTest() {
	logger := zaptest.NewLogger(t.T())
	log.SetLogger(logger)
}

func (t *OtelTracingTestSuite) TearDownTest() {
	// Reset global OTEL providers to avoid test interference
	otel.SetTracerProvider(noop.NewTracerProvider())
}

func (t *OtelTracingTestSuite) TestNewOtelTracing() {
	config := Config{
		IsOtelTracingEnabled: true,
	}
	serviceName := "test-service"

	otelTracing := NewOtelTracing(config, serviceName)

	t.NotNil(otelTracing)
	t.Equal(config, otelTracing.config)
	t.Equal(serviceName, otelTracing.serviceName)
}

func (t *OtelTracingTestSuite) TestInitialiseTraceWithTracingDisabled() {
	config := Config{
		IsOtelTracingEnabled: false,
	}
	serviceName := "test-service"

	ctx := context.Background()

	otelTracing := NewOtelTracing(config, serviceName)
	tracer, shutdown, err := otelTracing.InitialiseTrace(ctx)

	t.NoError(err)
	t.NotNil(tracer)
	t.Nil(shutdown)

	// Verify it's a noop tracer as it won't be recording spans
	_, span := tracer.Start(context.Background(), "test-span")
	t.False(span.IsRecording())
	span.End()
}

func (t *OtelTracingTestSuite) TestInitialiseTraceWithTracingEnabled() {
	config := Config{
		IsOtelTracingEnabled: true,
	}
	serviceName := "test-service"

	ctx := context.Background()

	otelTracing := NewOtelTracing(config, serviceName)
	tracer, shutdown, err := otelTracing.InitialiseTrace(ctx)

	t.NoError(err)
	t.NotNil(tracer)
	t.NotNil(shutdown)

	// Verify it's not a noop tracer
	_, span := tracer.Start(context.Background(), "test-span")
	t.True(span.IsRecording())
	span.End()

	err = shutdown(context.Background())
	t.NoError(err)
}

func (t *OtelTracingTestSuite) TestGlobalOtelConfiguration() {
	config := Config{
		IsOtelTracingEnabled: false,
	}
	serviceName := "test-service"

	ctx := context.Background()
	otelTracing := NewOtelTracing(config, serviceName)
	_, _, err := otelTracing.InitialiseTrace(ctx)
	t.NoError(err)

	globalTracer := otel.GetTracerProvider()
	t.NotNil(globalTracer)

	propagator := otel.GetTextMapPropagator()
	t.NotNil(propagator)
}
