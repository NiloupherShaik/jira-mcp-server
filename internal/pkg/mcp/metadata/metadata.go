package metadata

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/log"
)

func SSOIdFromMetadata(ctx context.Context, metadata *mcp.Meta) (string, error) {
	if metadata == nil || metadata.AdditionalFields == nil {
		log.WithContext(ctx).Error("metadata cannot be empty or nil.")
		return "", errors.New("no metadata provided in request")
	}

	ssoId, ok := metadata.AdditionalFields["ssoId"].(string)
	if !ok || ssoId == "" {
		log.WithContext(ctx).Error("no SSO ID provided in metadata.")
		return "", errors.New("ssoId is required to use the greeting tool")
	}

	return ssoId, nil
}
