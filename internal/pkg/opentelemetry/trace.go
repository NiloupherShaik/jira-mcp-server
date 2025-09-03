package opentelemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/nable-fusion/mcp-server-template/internal/pkg/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

type OtelTracing struct {
	config      Config
	serviceName string
}

const batchTimeout = 5 * time.Second

type TracingContext struct {
	TraceParent string `json:"traceparent"`
	TraceState  string `json:"tracestate"`
}

func NewOtelTracing(config Config, serviceName string) *OtelTracing {
	return &OtelTracing{
		config:      config,
		serviceName: serviceName,
	}
}

func (ot *OtelTracing) InitialiseTrace(ctx context.Context) (tracer trace.Tracer, shutdown func(ctx context.Context) error, err error) {
	contextLogger := log.WithContext(ctx).
		With(zap.Bool(log.IsOpenTelemetryEnabledLogKey, ot.config.IsOtelTracingEnabled))

	if !ot.config.IsOtelTracingEnabled {
		contextLogger.Warn("Using NOOP OTEL Trace Provider")

		tp := noop.NewTracerProvider()

		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

		return tp.Tracer(ot.serviceName), nil, nil
	}

	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			// TODO: Insecure is maybe fine since we hit the pod directly?
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithTimeout(time.Second),
		),
	)
	if err != nil {
		err = fmt.Errorf("failed to set open telemetry exporter: %w", err)
		contextLogger.Error("Error creating OTEL exporter", zap.Error(err))
		return nil, nil, err
	}

	contextLogger.Info("OpenTelemetry exporter created successfully")
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(batchTimeout)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracer = tp.Tracer(ot.serviceName)

	return tracer, exporter.Shutdown, nil
}
