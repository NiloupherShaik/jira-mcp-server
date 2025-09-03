package config

import "github.com/nable-fusion/fusion-cloud-common/pkg/appconfig"

const (
	defaultPort = "8080"
)

const (
	port = "PORT"
)

type Config struct {
	Port string
}

func NewConfig(envVarPrefix string) Config {
	av := appconfig.NewAppViper()
	av.SetEnvPrefix(envVarPrefix)

	av.SetAndBindDefaults(&map[string]interface{}{
		port: defaultPort,
	})

	av.AutomaticEnv()

	config := Config{
		Port: av.GetString(port),
	}

	return config
}
