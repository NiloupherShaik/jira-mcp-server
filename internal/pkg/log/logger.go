// Package logging provides a logger that includes values from context in traces.
package log

import (
	"context"

	"github.com/nable-fusion/fusion-cloud-common/pkg/log"
	"go.uber.org/zap"
)

var contextHandler ContextHandler

// ContextHandler is an interface for handling context values.
type ContextHandler interface {
	DecorateFromContext(ctx context.Context, logger *zap.Logger) *zap.Logger
}

// SetContextHandler sets the context handler.
func SetContextHandler(handler ContextHandler) {
	contextHandler = handler
}

// NewLogger creates a new logger in the fusion common log package.
func NewLogger() {
	log.NewLogger(log.NewConfig())
}

// WithContext returns a new logger where values from context are then automatically included in traces.
func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return log.WithContext(context.Background()) //nolint:contextcheck // this is fine
	}

	logger := log.WithContext(ctx)

	if contextHandler == nil {
		return logger
	}

	return contextHandler.DecorateFromContext(ctx, logger)
}
