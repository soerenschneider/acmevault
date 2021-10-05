package internal

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
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

	CertificatesRetrieved = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificates_retrieved_total",
		Help:      "Total amount of certificates retrieved from ACME provider",
	}, []string{"domain"})

	CertificatesRetrievalErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificate_retrieve_errors_total",
		Help:      "Total errors while trying to retrieve certificates from ACME provider",
	}, []string{"domain"})

	CertificatesRenewed = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificates_renewals_total",
		Help:      "Total number of renewed certificates",
	}, []string{"domain"})

	CertificatesRenewErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "certificates_renewal_errors_total",
		Help:      "Total errors while trying to renew certificates",
	}, []string{"domain"})

	CertWrites = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "certificates_written_total",
		Help:      "Total number of certificates written total",
	}, []string{"domain", "subsystem"})

	CertWriteError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "certificates_write_errors_total",
		Help:      "Total number of errors while writing the certificate",
	}, []string{"domain", "subsystem"})

	CertErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "certificate_errors_total",
		Help:      "Total number of errors while handling certificates",
	}, []string{"domain", "desc"})

	CertExpiryTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "certificate_expiry_time",
		Help:      "Timestamp of certificate expiry",
	}, []string{"domain"})

	HooksExecutionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "client",
		Name:      "hooks_invocation_errors",
		Help:      "Errors while invoking the hooks",
	})

	VaultTokenExpiryTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "timestamp",
		Help:      "Date of last measure",
	})
)

func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal().Msgf("Can not start metrics server at %s: %v", addr, err)
	}
}

func WriteMetrics(path string) error {
	log.Info().Msgf("Dumping metrics to %s", path)
	metrics, err := dumpMetrics()
	if err != nil {
		log.Info().Msgf("Error dumping metrics: %v", err)
		return err
	}

	err = ioutil.WriteFile(path, []byte(metrics), 0644)
	if err != nil {
		log.Info().Msgf("Error writing metrics to '%s': %v", path, err)
	}
	return err
}

func dumpMetrics() (string, error) {
	var buf = &bytes.Buffer{}
	enc := expfmt.NewEncoder(buf, expfmt.FmtText)

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return "", err
	}

	for _, f := range families {
		if err := enc.Encode(f); err != nil {
			log.Info().Msgf("could not encode metric: %s", err.Error())
		}
	}

	return buf.String(), nil
}
