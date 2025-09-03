package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nable-fusion/mcp-server-template/internal/pkg/appcontext"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/logdecorator"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/nable-fusion/mcp-server-template/internal/app/mcp-server/config"
	mcpServerService "github.com/nable-fusion/mcp-server-template/internal/app/mcp-server/service"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/log"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/mcp"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/opentelemetry"
	"go.uber.org/zap"
)

const (
	envVarPrefix = "MCP_SERVER"
)

func main() {
	// create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// app level config
	cfg := config.NewConfig(envVarPrefix)

	// create a global logger
	log.NewLogger()

	logdecorator.SetupLogger()

	// create config for MCP server
	serverCfg := mcp.NewConfig(envVarPrefix)

	// attach server name to loggers context
	ctx = appcontext.WithServerName(ctx, serverCfg.Name)

	// create tracer for open telemetry
	otelConfig := opentelemetry.NewConfig(envVarPrefix)
	otelTraceProvider := opentelemetry.NewOtelTracing(otelConfig, serverCfg.Name)

	otelTracer, shutdown, err := otelTraceProvider.InitialiseTrace(ctx)
	if err != nil {
		log.WithContext(ctx).Fatal("Error initialising OpenTelemetry tracer", zap.Error(err))
	}
	// shutdown the tracer when the application exits
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.WithContext(ctx).Error("Error during shutdown", zap.Error(err))
		}
	}()

	// attaches tools to the server
	mcpToolService := mcpServerService.NewMcpServerToolService(cfg)
	log.WithContext(ctx).Info("Registering tools with MCP server", zap.Int("toolCount", len(mcpToolService.UnboxTools())))

	mcpServer := mcp.NewServer(serverCfg, mcpToolService.UnboxTools(), otelTracer)

	log.WithContext(ctx).Info("MCP server started successfully", zap.String("serverName", serverCfg.Name), zap.String("port", cfg.Port))
	// receive on signal channel to gracefully terminate
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// start server in non-blocking goroutine
	go func() {
		log.WithContext(ctx).Info("Starting MCP server...")
		// spinning up a new mux to allow us to implement our own handlers, alongside the MCP handler.
		mux := http.NewServeMux()
		mux.Handle("/mcp", mcpServer)

		// adding this otel handler to the mux, so that we can automatically trace requests
		otelHandler := otelhttp.NewHandler(mux, serverCfg.Name)

		httpServer := &http.Server{
			Addr:              fmt.Sprintf(":%s", cfg.Port),
			Handler:           otelHandler,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// set up a logger with context
		log.WithContext(ctx).Info(fmt.Sprintf("🚀 Starting MCP server on http://localhost:%s", cfg.Port))

		if err := httpServer.ListenAndServe(); err != nil {
			log.WithContext(ctx).Fatal("Error starting MCP server", zap.Error(err))
		}
	}()

	// block until we receive on signal channel
	sig := <-sigChan
	log.WithContext(ctx).Info("Recevied on signal channel, gracefully shutting down...", zap.Any(log.SignalLogKey, sig))
}
