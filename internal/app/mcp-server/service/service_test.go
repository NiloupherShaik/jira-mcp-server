package service

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nable-fusion/fusion-cloud-common/pkg/log"
	"github.com/nable-fusion/mcp-server-template/internal/app/mcp-server/config"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"testing"
)

type ServiceTest struct {
	suite.Suite

	subject *McpServerToolService
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTest))
}

func (st *ServiceTest) SetupTest() {
	cfg := config.Config{
		Port: "8080",
	}

	logger := zaptest.NewLogger(st.T())
	log.SetLogger(logger)

	st.subject = NewMcpServerToolService(cfg)
}

func (st *ServiceTest) TestUnboxToolsReturnsTools() {
	resp := st.subject.UnboxTools()

	st.NotEmpty(resp, "Expected tools to be returned")
	st.Len(resp, 3, "Expected 3 tools to be returned")
}

func (st *ServiceTest) TestUVIndexToolHandler() {

	st.Run("Returns an MCP tool with correct fields", func() {
		uvTool := st.subject.UVIndexTool()

		expectedTool := mcp.NewTool("uv-index-tool",
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

		st.NotEmpty(uvTool.Tool, "Expected UVIndex tool to be returned")

		st.Equal(expectedTool.Name, uvTool.Tool.Name)
		st.Equal(expectedTool.Description, uvTool.Tool.Description)
		st.Equal(expectedTool.InputSchema, uvTool.Tool.InputSchema)
		st.Equal(expectedTool.Annotations, uvTool.Tool.Annotations)

		st.NotEmpty(uvTool.Handler, "Expected tool handler to be set")

	})

	st.Run("Returns an MCP tool error when invalid arguments are passing in", func() {
		ctx := context.Background()
		uvTool := st.subject.UVIndexTool()
		req := mcp.CallToolRequest{}

		uv := 9

		req.Params.Arguments = map[string]int{
			"uv-index": uv,
		}

		expectedResp := mcp.NewToolResultError("your UV index input must be within 1-8")

		resp, err := uvTool.Handler(ctx, req)

		// validate we get a ToolResultError
		st.True(resp.IsError)
		st.Nil(err)
		// validate error = expected error
		st.Equal(expectedResp, resp)
	})

	st.Run("Return an MCP tool result when current UV is too high", func() {
		ctx := context.Background()
		uvTool := st.subject.UVIndexTool()
		req := mcp.CallToolRequest{}

		uv := 6
		req.Params.Arguments = map[string]int{
			"uv-index": uv,
		}

		// note the current UV Index is hardcoded to 7 in the tool handler
		expectedResp := mcp.NewToolResultText(fmt.Sprintf("UV Index is too high, current UV index: %v, desired UV index: %v", 7, uv))

		resp, err := uvTool.Handler(ctx, req)

		// validate we get a ToolResultError
		st.False(resp.IsError)
		st.Nil(err)
		// validate error = expected error
		st.Equal(expectedResp, resp)
	})

	st.Run("Returns an MCP tool result when current UV is within range", func() {
		ctx := context.Background()
		uvTool := st.subject.UVIndexTool()
		req := mcp.CallToolRequest{}

		uv := 8

		// If you're using strongly typed arguemnts via the NewTypedToolHandler method, then your struct should have the below json annotation on it.
		req.Params.Arguments = map[string]int{
			"uv-index": uv,
		}

		// note the current UV Index is hardcoded to 7 in the tool handler
		expectedResp := mcp.NewToolResultText(fmt.Sprintf("UV Index is within acceptable levels, current UV index: %v, desired UV index: %v", 7, uv))

		resp, err := uvTool.Handler(ctx, req)

		// validate we get a ToolResultError
		st.False(resp.IsError)
		st.Nil(err)
		// validate error = expected error
		st.Equal(expectedResp, resp)
	})

}

func (st *ServiceTest) TestWeatherToolHandler() {
	st.Run("Returns an MCP tool with correct fields", func() {
		uvTool := st.subject.WeatherTool()

		expectedTool := mcp.NewTool("get-weather",
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

		st.NotEmpty(uvTool.Tool, "Expected Weather tool to be returned")

		st.Equal(expectedTool.Name, uvTool.Tool.Name)

		st.Equal(expectedTool.Description, uvTool.Tool.Description)
		st.Equal(expectedTool.InputSchema, uvTool.Tool.InputSchema)
		st.Equal(expectedTool.Annotations, uvTool.Tool.Annotations)

		st.NotEmpty(uvTool.Handler, "Expected tool handler to be set")
	})

	st.Run("Returns an MCP tool error when invalid arguments are passing in", func() {
		ctx := context.Background()
		weatherTool := st.subject.WeatherTool()
		req := mcp.CallToolRequest{}

		location := ""

		req.Params.Arguments = map[string]string{
			"location": location,
		}

		expectedResp := mcp.NewToolResultError("location parameter is required.")

		resp, err := weatherTool.Handler(ctx, req)

		// validate we get a ToolResultError
		st.True(resp.IsError)
		st.Nil(err)
		// validate error = expected error
		st.Equal(expectedResp, resp)
	})

	st.Run("Returns an MCP tool result when valid location is passed in", func() {
		ctx := context.Background()
		weatherTool := st.subject.WeatherTool()
		req := mcp.CallToolRequest{}

		location := "dundee"

		req.Params.Arguments = map[string]string{
			"location": location,
		}

		expectedResp := mcp.NewToolResultText(fmt.Sprintf("Weather is great, 26° celsius."))

		resp, err := weatherTool.Handler(ctx, req)

		st.False(resp.IsError)
		st.Nil(err)
		st.Equal(expectedResp, resp)
	})

	st.Run("Returns an MCP tool error when invalid location is passed in", func() {
		ctx := context.Background()
		weatherTool := st.subject.WeatherTool()
		req := mcp.CallToolRequest{}

		location := "london"

		req.Params.Arguments = map[string]string{
			"location": location,
		}

		expectedResp := mcp.NewToolResultError(fmt.Sprintf("Sorry, we can't get the weather for that location."))

		resp, err := weatherTool.Handler(ctx, req)

		st.True(resp.IsError)
		st.Nil(err)
		st.Equal(expectedResp, resp)
	})
}

func (st *ServiceTest) TestGreetingTool() {
	st.Run("Returns an MCP tool with correct fields", func() {
		greetingTool := st.subject.GreetingTool()

		expectedTool := mcp.NewTool("greeting",
			mcp.WithDescription("Returns a greeting."),
		)

		st.NotEmpty(greetingTool.Tool, "Expected tool to be returned")

		st.Equal(expectedTool.Name, greetingTool.Tool.Name)
		st.Equal(expectedTool.Description, greetingTool.Tool.Description)
		st.Equal(expectedTool.InputSchema, greetingTool.Tool.InputSchema)
		st.Equal(expectedTool.Annotations, greetingTool.Tool.Annotations)

		st.NotEmpty(greetingTool.Handler, "Expected tool handler to be set")

	})

	st.Run("Returns an MCP tool result with the users sso id when called", func() {
		ctx := context.Background()
		greetingTool := st.subject.GreetingTool()
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

		expectedResp := mcp.NewToolResultText(fmt.Sprintf("Hello, fake-sso-id!"))

		resp, err := greetingTool.Handler(ctx, req)

		st.False(resp.IsError)
		st.Nil(err)
		st.Equal(expectedResp, resp)
	})

	st.Run("Returns an MCP tool result error when no metadata field is passed in", func() {
		ctx := context.Background()
		greetingTool := st.subject.GreetingTool()
		req := mcp.CallToolRequest{}

		expectedResp := mcp.NewToolResultError(fmt.Sprintf("ssoId is required to use the greeting tool"))

		resp, err := greetingTool.Handler(ctx, req)

		st.True(resp.IsError)
		st.Nil(err)
		st.Equal(expectedResp, resp)
	})

	st.Run("Returns an MCP tool result error when no metadata AdditionalFields field is passed in", func() {
		ctx := context.Background()
		greetingTool := st.subject.GreetingTool()
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Meta: &mcp.Meta{},
			},
		}

		expectedResp := mcp.NewToolResultError(fmt.Sprintf("ssoId is required to use the greeting tool"))

		resp, err := greetingTool.Handler(ctx, req)

		st.True(resp.IsError)
		st.Nil(err)
		st.Equal(expectedResp, resp)
	})

	st.Run("Returns an MCP tool result error when no ssoId is passed in", func() {
		ctx := context.Background()
		greetingTool := st.subject.GreetingTool()
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Meta: &mcp.Meta{
					ProgressToken: nil,
					AdditionalFields: map[string]any{
						"notAnSsoId": "",
					},
				},
			},
		}

		expectedResp := mcp.NewToolResultError(fmt.Sprintf("ssoId is required to use the greeting tool"))

		resp, err := greetingTool.Handler(ctx, req)

		st.True(resp.IsError)
		st.Nil(err)
		st.Equal(expectedResp, resp)
	})

	st.Run("Returns an MCP tool result error when empty ssoId is passed in", func() {
		ctx := context.Background()
		greetingTool := st.subject.GreetingTool()
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Meta: &mcp.Meta{
					ProgressToken: nil,
					AdditionalFields: map[string]any{
						"ssoId": "",
					},
				},
			},
		}

		expectedResp := mcp.NewToolResultError(fmt.Sprintf("ssoId is required to use the greeting tool"))

		resp, err := greetingTool.Handler(ctx, req)

		st.True(resp.IsError)
		st.Nil(err)
		st.Equal(expectedResp, resp)
	})
}
