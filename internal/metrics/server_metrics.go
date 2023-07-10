package metrics

import (
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
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
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msgf("Can not start metrics server")
	}
}
