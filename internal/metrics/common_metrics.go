package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "acmevault"
)

var (
	AwsDynCredentialsRequested = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "vault_aws_credentials_requested_total",
		Help:      "Total amount of dynamic AWS credentials requested",
	})

	AwsDynCredentialsRequestErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "vault_aws_credentials_request_errors_total",
		Help:      "Total amount of errors while trying to acquire dynamic AWS credentials",
	})

	CertErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificate_errors_total",
		Help:      "Total number of errors while handling certificates",
	}, []string{"desc"})
)
