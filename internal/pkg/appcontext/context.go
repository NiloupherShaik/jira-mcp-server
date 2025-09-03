// Package appcontext provides decorators for context.Context specific to this application.
package appcontext

import (
	"context"
	"time"
)

// using private structs to ensure keys are always different
// as context values should only be set and retrieved using this package, there's no risk of someone adding a key with the same value and overwriting it accidentally.
type (
	trackingIDKey struct{}
	startTimeKey  struct{}
	endtimeKey    struct{}
	serverNameKey struct{}
	ssoIDKey      struct{}
)

func WithServerName(ctx context.Context, serverName string) context.Context {
	return context.WithValue(ctx, serverNameKey{}, serverName)
}

func WithSSOId(ctx context.Context, ssoID string) context.Context {
	return context.WithValue(ctx, ssoIDKey{}, ssoID)
}

// WithTrackingID adds tracking ID to the context.
func WithTrackingID(ctx context.Context, trackingID string) context.Context {
	return context.WithValue(ctx, trackingIDKey{}, trackingID)
}

// WithStartTime adds start time to the context.
func WithStartTime(ctx context.Context, startTime time.Time) context.Context {
	return context.WithValue(ctx, startTimeKey{}, startTime)
}

// WithEndTime adds end time to the context.
func WithEndTime(ctx context.Context, endTime time.Time) context.Context {
	return context.WithValue(ctx, endtimeKey{}, endTime)
}

// ServerNameFromContext returns the server name from the context.
func ServerNameFromContext(ctx context.Context) string {
	return valueFromContext(ctx, serverNameKey{})
}

// TrackingIDFromContext returns the tracking ID from the context.
func TrackingIDFromContext(ctx context.Context) string {
	return valueFromContext(ctx, trackingIDKey{})
}

func SSOIdFromContext(ctx context.Context) string {
	return valueFromContext(ctx, ssoIDKey{})
}

// StartTimeFromContext returns the start time from the context.
func StartTimeFromContext(ctx context.Context) time.Time {
	if startTime, ok := ctx.Value(startTimeKey{}).(time.Time); ok {
		return startTime
	}

	return time.Time{}
}

// EndTimeFromContext returns the end time from the context.
func EndTimeFromContext(ctx context.Context) time.Time {
	if endTime, ok := ctx.Value(endtimeKey{}).(time.Time); ok {
		return endTime
	}

	return time.Time{}
}

func valueFromContext(ctx context.Context, key any) string {
	if valueStr, ok := ctx.Value(key).(string); ok {
		return valueStr
	}

	return ""
}
