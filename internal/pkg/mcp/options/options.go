package options

import (
	"context"
	"net/http"

	"github.com/nable-fusion/mcp-server-template/internal/pkg/appcontext"
)

func HeadersToContext() func(ctx context.Context, r *http.Request) context.Context {
	return func(ctx context.Context, r *http.Request) context.Context {
		ssoId := r.Header.Get("X-Sso-Id")
		trackingId := r.Header.Get("X-Tracking-Id")

		// if we find a header, add it to the request context, so we can read it in our toolHandler func
		if ssoId != "" {
			ctx = appcontext.WithSSOId(ctx, ssoId)
		}

		if trackingId != "" {
			ctx = appcontext.WithTrackingID(ctx, trackingId)
		}

		return ctx
	}
}
