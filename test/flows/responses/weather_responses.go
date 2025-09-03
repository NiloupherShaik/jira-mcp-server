package responses

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/onsi/gomega"
)

func ConfirmSuccesfulWeatherResponse(resp *mcp.CallToolResult) {
	gomega.Expect(resp.Content).ToNot(gomega.BeEmpty(), "weatherResp is empty")
	gomega.Expect(resp.IsError).To(gomega.BeFalse(), "Response should not return an MCP error")
}

func ConfirmErrorWeatherResponse(resp *mcp.CallToolResult) {
	gomega.Expect(resp.Content).ToNot(gomega.BeEmpty(), "weatherResp is empty")
	gomega.Expect(resp.IsError).To(gomega.BeTrue(), "Response should return an MCP error")
}
