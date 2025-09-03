package responses

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/onsi/gomega"
)

func ConfirmListToolsSuccess(toolResp, expectedTools []mcp.Tool) {
	gomega.Expect(toolResp).ToNot(gomega.BeNil(), "toolsResp is nil")
	gomega.Expect(toolResp).ToNot(gomega.BeEmpty(), "toolsResp is empty")
	gomega.Expect(toolResp).To(gomega.HaveLen(len(expectedTools)), "toolsResp length does not match expected length")

	// check if lists match
	gomega.Expect(toolResp).To(gomega.HaveExactElements(expectedTools))
}

func extractTextFromMCPResult(response *mcp.CallToolResult) string {
	gomega.Expect(response).ToNot(gomega.BeNil(), "CallToolResponse is nil")

	extractedContent, ok := response.Content[0].(mcp.TextContent)
	gomega.Expect(ok).To(gomega.BeTrue(), "Content is not of type TextContent")
	extractedText := extractedContent.Text

	gomega.Expect(extractedText).ToNot(gomega.BeEmpty(), "Extracted text is empty")

	return extractedText
}
