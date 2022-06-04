//go:build client
// +build client

package metrics

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/expfmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

var (
	AuthErrors = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "client",
		Name:      "auth_errors",
		Help:      "Errors while authenticating against the backend",
	})

	CertReadErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "client",
		Name:      "certificate_read_errors",
		Help:      "Errors while trying to read the certificate from the backend",
	}, []string{"domain"})

	CertClientExpiryTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "client",
		Name:      "certificate_expiry_time",
		Help:      "Timestamp of certificate expiry",
	}, []string{"domain"})

	HooksExecutionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "client",
		Name:      "hooks_invocation_errors",
		Help:      "Errors while invoking the hooks",
	})

	CertClientErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "client",
		Name:      "certificate_format_errors_total",
		Help:      "Total number of certificate errors",
	}, []string{"domain", "desc"})

	PersistCertErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "client",
		Name:      "write_certificate_errors_total",
		Help:      "Errors while trying to write received certificate to backend",
	}, []string{"domain"})
)

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
