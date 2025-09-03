package bootstrap

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/client/transport"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nable-fusion/mcp-server-template/test/flows/requests"
	"github.com/nable-fusion/test-automation-lib-go/v2/pkg/environmentutilities"
	"github.com/nable-fusion/test-automation-lib-go/v2/pkg/errorutilities"
	"github.com/nable-fusion/test-automation-lib-go/v2/pkg/jsonutilities"
	"github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"
)

type Environment struct {
	MCPServerHost string
}

func RegisterSuite(t *testing.T) {
	// mark as a helper function
	t.Helper()
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Integration Test Suite")
}

func SuiteSetup() []byte {
	environment := Environment{
		MCPServerHost: environmentutilities.ConfirmEnvVariableIsSet("MCP_SERVER_HOST"),
	}

	return jsonutilities.ParseModelAsJson(environment)
}

func ParseEnvironment(envBytes []byte) *Environment {
	env := &Environment{}
	jsonutilities.ParseJsonAsModel(envBytes, &env)

	return env
}

func SetupMCPClient(env *Environment) *requests.MCPClient {
	// TODO: find a way of doing this on each tool invocation
	headers := map[string]string{
		"X-Tracking-ID": "fake-tracking-id",
	}

	c := createMCPClient(env.MCPServerHost, headers)

	return requests.NewMCPClient(c)
}

func createMCPClient(host string, headers map[string]string) *client.Client {
	mcpClient, err := client.NewStreamableHttpClient(
		fmt.Sprintf("%v/%v", host, "mcp"),
		transport.WithHTTPHeaders(headers),
	)

	ctx := context.Background()

	errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "Failed to Create New MCP Client!")

	err = mcpClient.Start(ctx)
	errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "Failed to Start MCP Client!")

	initRequest := mcp.InitializeRequest{}

	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp-test-client",
		Version: "1.0.0",
	}

	_, err = mcpClient.Initialize(context.TODO(), initRequest)

	errorutilities.ConfirmErrorHasNotOccurredWithMessage(err, "Failed to Initialize MCP Client!")

	return mcpClient
}
