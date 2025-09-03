package requests

import (
	"context"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nable-fusion/test-automation-lib-go/v2/pkg/errorutilities"
	"github.com/onsi/gomega"
)

type MCPClient struct {
	Client *client.Client
}

func NewMCPClient(c *client.Client) *MCPClient {
	return &MCPClient{
		Client: c,
	}
}

func (m *MCPClient) ListTools() []mcp.Tool {
	resp, err := m.Client.ListTools(context.TODO(), mcp.ListToolsRequest{})
	errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "Failed to List Tools!")

	gomega.Expect(resp).ToNot(gomega.BeNil(), "ListToolsResponse is nil")
	gomega.Expect(resp.Tools).ToNot(gomega.BeNil(), "ListToolsResponse.Tools is nil")
	gomega.Expect(resp.Tools).ToNot(gomega.BeEmpty(), "ListToolsResponse.Tools is empty")

	return resp.Tools
}
