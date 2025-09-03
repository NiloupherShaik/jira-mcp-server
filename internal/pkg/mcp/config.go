package mcp

import "github.com/nable-fusion/fusion-cloud-common/pkg/appconfig"

const (
	defaultName     = "mcp-server"
	defaultBasePath = "mcp"
)

const (
	name     = "SERVER_NAME"
	basePath = "SERVER_BASE_PATH"
)

type Config struct {
	Name     string
	BasePath string
}

func NewConfig(envVarPrefix string) Config {
	av := appconfig.NewAppViper()
	av.SetEnvPrefix(envVarPrefix)

	av.SetAndBindDefaults(&map[string]interface{}{
		name:     defaultName,
		basePath: defaultBasePath,
	})

	av.AutomaticEnv()

	config := Config{
		Name:     av.GetString(name),
		BasePath: av.GetString(basePath),
	}

	return config
}
