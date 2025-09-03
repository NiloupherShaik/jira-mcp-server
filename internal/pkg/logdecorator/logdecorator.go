// Package logdecorator provides logger decorates based on context values.
package logdecorator

import (
	"context"

	"github.com/nable-fusion/mcp-server-template/internal/pkg/appcontext"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/log"

	"go.uber.org/zap"
)

// ContextHandler is the context handler.
type ContextHandler struct{}

// SetupLogger sets up the logger.
func SetupLogger() {
	log.SetContextHandler(&ContextHandler{})
}

// DecorateFromContext decorates the logger with values from context.
func (h *ContextHandler) DecorateFromContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if val := appcontext.ServerNameFromContext(ctx); val != "" {
		logger = logger.With(zap.String(log.McpServerNameLogKey, val))
	}

	if val := appcontext.SSOIdFromContext(ctx); val != "" {
		logger = logger.With(zap.String(log.SSOIdKey, val))
	}

	if val := appcontext.TrackingIDFromContext(ctx); val != "" {
		logger = logger.With(zap.String(log.TrackingIdKey, val))
	}

	return logger
}
