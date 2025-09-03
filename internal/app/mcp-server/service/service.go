package service

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/nable-fusion/mcp-server-template/internal/pkg/log"
	"github.com/nable-fusion/mcp-server-template/internal/pkg/mcp/metadata"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nable-fusion/mcp-server-template/internal/app/mcp-server/arguments/weather"
	"github.com/nable-fusion/mcp-server-template/internal/app/mcp-server/config"
)

type McpServerToolService struct {
	cfg config.Config
}

func NewMcpServerToolService(cfg config.Config) *McpServerToolService {
	return &McpServerToolService{
		cfg: cfg,
	}
}

// UnboxTools will return all tools to provide to your MCP server.
// NOTE: your MCP server should be scoped to a given domain / function - This is a weather MCP server, and should only include tools for that domain.
// This includes a greeting tool purely to showcase a really simple tool.
func (m *McpServerToolService) UnboxTools() []server.ServerTool {
	tools := []server.ServerTool{
		m.GreetingTool(),
		m.WeatherTool(),
		m.UVIndexTool(),
	}

	return tools
}

func (m *McpServerToolService) UVIndexTool() server.ServerTool {
	// define a new tool
	tool := server.ServerTool{}

	// build up the underlying MCP tool via the options pattern - easier than trying to build up input schema yourself.
	mcpTool := mcp.NewTool("uv-index-tool",
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

	// this tool handler takes in strongly typed arguments, you can do this without however this is clearer.
	uvIndexTool := func(ctx context.Context, request mcp.CallToolRequest, args weather.CurrentUVIndexToolargs) (*mcp.CallToolResult, error) {
		if args.UVIndex < 1 || args.UVIndex > 8 {
			log.WithContext(ctx).Error("your UV index input must be within 1-8")
			return mcp.NewToolResultError("your UV index input must be within 1-8"), nil
		}

		currentIndex := 7

		if currentIndex > args.UVIndex {
			return mcp.NewToolResultText(fmt.Sprintf("UV Index is too high, current UV index: %v, desired UV index: %v", currentIndex, args.UVIndex)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("UV Index is within acceptable levels, current UV index: %v, desired UV index: %v", currentIndex, args.UVIndex)), nil
	}

	tool.Tool = mcpTool
	tool.Handler = mcp.NewTypedToolHandler(uvIndexTool)

	return tool
}

// WeatherTool interacts with a weather API to return the weather for a given location.
func (m *McpServerToolService) WeatherTool() server.ServerTool {
	// define a new tool
	tool := server.ServerTool{}

	// build up the underlying MCP tool via the options pattern - easier than trying to build up input schema yourself.
	mcpTool := mcp.NewTool("get-weather",
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

	// this tool handler takes in strongly typed arguments, you can do this without however this is clearer.
	weatherTool := func(ctx context.Context, request mcp.CallToolRequest, args weather.CurrentWeatherToolArgs) (*mcp.CallToolResult, error) {
		if args.Location == "" {
			log.WithContext(ctx).Error("location parameter is required.")
			return mcp.NewToolResultError("location parameter is required."), nil
		}

		// avoid weird case match issues
		location := strings.ToLower(args.Location)

		switch location {
		case "dundee":
			return mcp.NewToolResultText("Weather is great, 26° celsius."), nil
		case "edinburgh":
			return mcp.NewToolResultText("Weather is terrible, -3° celsius."), nil
		default:
			log.WithContext(ctx).Error("this isn't a real api, but you get the point. Build up your api request and return the response, and handle errors with ToolResultError")
			return mcp.NewToolResultError("Sorry, we can't get the weather for that location."), nil
		}
	}

	tool.Tool = mcpTool
	tool.Handler = mcp.NewTypedToolHandler(weatherTool)

	return tool
}

// GreetingTool is a very simple tool that takes in no parameters.
func (m *McpServerToolService) GreetingTool() server.ServerTool {
	// define a new tool
	tool := server.ServerTool{
		Tool:    mcp.NewTool("greeting"),
		Handler: nil,
	}

	// describe tool, this allows the agent to understand what a tool does and will influence whether or not it uses it.
	tool.Tool.Description = "Returns a greeting."

	// define tool behaviour / handler - this always returns hello world / hello sso-id.
	tool.Handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ssoId, err := metadata.SSOIdFromMetadata(ctx, request.Params.Meta)
		if err != nil {
			log.WithContext(ctx).Error("failed to get SSO ID from metadata", zap.Error(err))
			return mcp.NewToolResultError("ssoId is required to use the greeting tool"), nil
		}

		log.WithContext(ctx).Info("called into greeting tool", zap.Any(log.MetadataKey, request.Params.Meta.AdditionalFields))

		greeting := fmt.Sprintf("Hello, %s!", ssoId)
		return mcp.NewToolResultText(greeting), nil
	}

	return tool
}
