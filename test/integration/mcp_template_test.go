package integration

import (
	"context"
	"github.com/nable-fusion/test-automation-lib-go/v2/pkg/errorutilities"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nable-fusion/mcp-server-template/test/flows/responses"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Template MCP server >", func() {
	ginkgo.It("should successfully list tools", func() {
		resp := mcpClient.ListTools()

		expectedTools := getTools()

		responses.ConfirmListToolsSuccess(resp, expectedTools)
	})

	ginkgo.Describe("Tool: greeting >", func() {
		ginkgo.It("should successfully return greeting with an SSO id meta field", func() {
			ctx := context.Background()

			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Meta: &mcp.Meta{
						ProgressToken: nil,
						AdditionalFields: map[string]any{
							"ssoId": "fake-sso-id",
						},
					},
				},
			}
			req.Params.Name = "greeting"

			resp, err := mcpClient.Client.CallTool(ctx, req)
			errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "failed to call greeting tool")
			responses.ConfirmGreetingResponseWithSsoId(resp)
		})

		ginkgo.It("should return error without metadata", func() {
			ctx := context.Background()

			req := mcp.CallToolRequest{}

			req.Params.Name = "greeting"

			resp, err := mcpClient.Client.CallTool(ctx, req)
			errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "failed to call greeting tool")
			responses.ConfirmGreetingResponseWithoutMetadata(resp)
		})

		ginkgo.It("should return error without ssoId", func() {
			ctx := context.Background()

			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Meta: &mcp.Meta{
						ProgressToken: nil,
						AdditionalFields: map[string]any{
							"noSsoId": "",
						},
					},
				},
			}
			req.Params.Name = "greeting"

			resp, err := mcpClient.Client.CallTool(ctx, req)
			errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "failed to call greeting tool")
			responses.ConfirmGreetingResponseWithoutSsoId(resp)
		})
	})

	ginkgo.Describe("Tool: get-weather >", func() {
		ginkgo.It("Should return a successful response", func() {
			ctx := context.Background()
			req := mcp.CallToolRequest{}

			req.Params.Name = "get-weather"

			req.Params.Arguments = map[string]interface{}{
				"location": "dundee",
			}

			resp, err := mcpClient.Client.CallTool(ctx, req)
			errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "failed to call get-weather tool")

			responses.ConfirmSuccesfulWeatherResponse(resp)
		})

		ginkgo.It("Should return an error response when a bad request is made", func() {
			ctx := context.Background()
			req := mcp.CallToolRequest{}

			req.Params.Name = "get-weather"

			req.Params.Arguments = map[string]interface{}{}

			resp, err := mcpClient.Client.CallTool(ctx, req)
			errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "failed to call get-weather tool")

			responses.ConfirmErrorWeatherResponse(resp)
		})

		ginkgo.It("Should return error MCP response when no request body is provided", func() {
			ctx := context.Background()
			req := mcp.CallToolRequest{}

			req.Params.Name = "get-weather"

			resp, err := mcpClient.Client.CallTool(ctx, req)
			errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "failed to call get-weather tool")

			gomega.Expect(resp.IsError).To(gomega.BeTrue(), "Response should return an MCP error")
		})
	})
})

func getTools() []mcp.Tool {
	weatherTool := mcp.NewTool("get-weather",
		// describe tool, this allows the agent to understand what a tool does and will influence whether or not it uses it.
		mcp.WithDescription("Gets the weather forecast for a given location"),
		// define required tool parameters, for this example the location is required.
		mcp.WithString("location",
			// signify that it's required.
			mcp.Required(),
			// describe the tool parameter clearly, this allows the agent to build this up.
			mcp.Description("The location for which you want the weather forecast."),
		),
	)

	uvIndexTool := mcp.NewTool("uv-index-tool",
		// describe tool, this allows the agent to understand what a tool does and will influence whether or not it uses it.
		mcp.WithDescription("Determines whether or not the current UV index is too high when compared to an input UV index"),
		// define required tool parameters, for this example the location is required.
		mcp.WithNumber("uv-index",
			// signify that it's required.
			mcp.Required(),
			// describe the tool parameter clearly, this allows the agent to build this up.
			mcp.Description("The max UV index you want to compare."),
			// you can set min/max param values.
			mcp.Min(1),
			mcp.Max(8),
		),
	)

	greetingTool := mcp.NewTool("greeting", mcp.WithDescription("Returns a greeting."))

	return []mcp.Tool{
		weatherTool,
		greetingTool,
		uvIndexTool,
	}
}
