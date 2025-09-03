package integration

import (
	"testing"

	"github.com/nable-fusion/mcp-server-template/test/flows/bootstrap"
	"github.com/nable-fusion/mcp-server-template/test/flows/requests"

	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/format"
)

func TestSuite(t *testing.T) {
	bootstrap.RegisterSuite(t)
}

var mcpClient *requests.MCPClient

var _ = SynchronizedBeforeSuite(
	//nolint:gocritic // We cant inline this or integration tests fail
	func() []byte {
		return bootstrap.SuiteSetup()
	},
	func(envBytes []byte) {
		env := bootstrap.ParseEnvironment(envBytes)

		mcpClient = bootstrap.SetupMCPClient(env)

		format.TruncatedDiff = false
	},
)

var _ = SynchronizedAfterSuite(
	func() {
		// This runs on each parallel node after all specs complete
	},
	func() {
		// This runs only once after all parallel nodes are done
		if mcpClient != nil && mcpClient.Client != nil {
			_ = mcpClient.Client.Close() //nolint:errcheck // We don't care if this fails, as the test suite is ending
		}
	},
)
