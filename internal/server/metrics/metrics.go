package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
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

	ServerLatestIterationTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "latest_iteration_time_seconds",
		Help:      "Latest invocation of the server",
	})

	CertificatesRetrieved = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificates_retrieved_total",
		Help:      "Total amount of certificates retrieved from ACME provider",
	})

	CertificatesRetrievalErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificate_retrieve_errors_total",
		Help:      "Total errors while trying to retrieve certificates from ACME provider",
	})

	CertificatesRenewals = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificates_renewals_total",
		Help:      "Total number of renewed certificates",
	})

	CertificatesRenewErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificates_renewal_errors_total",
		Help:      "Total errors while trying to renew certificates",
	})

	CertWrites = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificates_written_total",
		Help:      "Total number of certificates written total",
	}, []string{"subsystem"})

	CertWriteError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificates_write_errors_total",
		Help:      "Total number of errors while writing the certificate",
	}, []string{"subsystem"})

	CertErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificate_errors_total",
		Help:      "Total number of errors while handling certificates",
	}, []string{"desc"})

	CertExpiryTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificate_expiry_time",
		Help:      "Timestamp of certificate expiry",
	}, []string{"domain"})
)

func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal().Msgf("Can not start metrics server at %s: %v", addr, err)
	}
}
