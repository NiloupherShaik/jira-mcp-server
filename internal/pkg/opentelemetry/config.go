package opentelemetry

import "github.com/nable-fusion/fusion-cloud-common/pkg/appconfig"

const (
	isOtelTracingEnabledKey = "IS_OTEL_TRACING_ENABLED"
)

type Config struct {
	IsOtelTracingEnabled bool
}

func NewConfig(envVarPrefix string) Config {
	av := appconfig.NewAppViper()
	av.SetEnvPrefix(envVarPrefix)

	av.SetAndBindDefaults(&map[string]interface{}{
		isOtelTracingEnabledKey: true,
	})

	av.AutomaticEnv()

	config := Config{
		IsOtelTracingEnabled: av.GetBool(isOtelTracingEnabledKey),
	}

	return config
}
