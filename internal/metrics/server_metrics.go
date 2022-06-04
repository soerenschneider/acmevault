//go:build server
// +build server

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
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

	CertServerExpiryTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
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
