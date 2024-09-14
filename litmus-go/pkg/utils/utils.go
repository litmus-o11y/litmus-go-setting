package utils

import "net/url"

const (
	OTELExporterOTLPEndpoint = "otel-collector.observability.svc.cluster.local:4317"
)

func HttpTimeout(err error) bool {
	httpErr := err.(*url.Error)
	if httpErr != nil {
		return httpErr.Timeout()
	}
	return false
}
