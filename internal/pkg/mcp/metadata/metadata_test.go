package metadata

import (
	"context"
	"errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nable-fusion/fusion-cloud-common/pkg/log"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"testing"
)

type MetadataTest struct {
	suite.Suite
}

func TestMetadataTestSuite(t *testing.T) {
	suite.Run(t, new(MetadataTest))
}

func (t *MetadataTest) SetupTest() {
	logger := zaptest.NewLogger(t.T())
	log.SetLogger(logger)
}

func (t *MetadataTest) TestSSOIdFromMetadata() {

	expectedEmptyMetadataErr := errors.New("no metadata provided in request")

	t.Run("Returns an error if metadata is nil", func() {
		ctx := context.Background()
		ssoId, err := SSOIdFromMetadata(ctx, nil)

		t.Error(err, "Expected error when metadata is nil")
		t.EqualError(err, expectedEmptyMetadataErr.Error(), "Expected specific error for nil metadata")
		t.Empty(ssoId, "Expected empty SSO ID when metadata is nil")
	})

	t.Run("Returns an error if metadata.AdditionalFields is nil", func() {
		ctx := context.Background()

		metadata := &mcp.Meta{}
		ssoId, err := SSOIdFromMetadata(ctx, metadata)

		t.Error(err, "Expected error when metadata is nil")
		t.EqualError(err, expectedEmptyMetadataErr.Error(), "Expected specific error for nil metadata")
		t.Empty(ssoId, "Expected empty SSO ID when metadata is nil")
	})

	t.Run("Returns an error if metadata doesnt have ssoId key", func() {
		ctx := context.Background()

		metadata := &mcp.Meta{
			AdditionalFields: map[string]any{
				"otherKey": "someValue",
			},
		}
		ssoId, err := SSOIdFromMetadata(ctx, metadata)

		t.Error(err, "Expected error when ssoId is nil")
		t.EqualError(err, "ssoId is required to use the greeting tool", "Expected specific error for nil metadata")
		t.Empty(ssoId, "Expected empty SSO ID when metadata is nil")
	})

	t.Run("Returns an error if ssoId is empty", func() {
		ctx := context.Background()

		metadata := &mcp.Meta{
			AdditionalFields: map[string]any{
				"ssoId": "",
			},
		}
		ssoId, err := SSOIdFromMetadata(ctx, metadata)

		t.Error(err, "Expected error when ssoId is nil")
		t.EqualError(err, "ssoId is required to use the greeting tool", "Expected specific error for nil metadata")
		t.Empty(ssoId, "Expected empty SSO ID when metadata is nil")
	})

	t.Run("Returns a string if metadata has an ssoId key", func() {
		ctx := context.Background()

		metadata := &mcp.Meta{
			AdditionalFields: map[string]any{
				"ssoId": "someValue",
			},
		}
		ssoId, err := SSOIdFromMetadata(ctx, metadata)

		t.Nil(err, "Expected no error when SSO id is present")
		t.NotEmpty(ssoId, "Expected SSO ID when metadata is nil")
		t.Equal("someValue", ssoId, "Expected SSO ID to match the value in metadata")
	})
}
