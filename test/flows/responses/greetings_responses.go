package responses

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/onsi/gomega"
)

func ConfirmGreetingResponseWithSsoId(resp *mcp.CallToolResult) {
	text := extractTextFromMCPResult(resp)

	gomega.Expect(text).ToNot(gomega.BeEmpty(), "greetingResp is empty")
	gomega.Expect(text).To(gomega.Equal("Hello, fake-sso-id!"))
}

func ConfirmGreetingResponseWithoutMetadata(resp *mcp.CallToolResult) {
	text := extractTextFromMCPResult(resp)

	gomega.Expect(text).ToNot(gomega.BeEmpty(), "greetingResp is empty")
	gomega.Expect(text).To(gomega.Equal("ssoId is required to use the greeting tool"))
}

func ConfirmGreetingResponseWithoutSsoId(resp *mcp.CallToolResult) {
	text := extractTextFromMCPResult(resp)

	gomega.Expect(text).ToNot(gomega.BeEmpty(), "greetingResp is empty")
	gomega.Expect(text).To(gomega.Equal("ssoId is required to use the greeting tool"))
}
