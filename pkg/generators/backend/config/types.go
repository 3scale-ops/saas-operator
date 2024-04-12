package config

type Secret string

func (s Secret) String() string { return string(s) }

const (
	BackendInternalApiSecret     Secret = "backend-internal-api"
	BackendErrorMonitoringSecret Secret = "backend-error-monitoring"
	BackendSystemEventsSecret    Secret = "backend-system-events-hook"
)
