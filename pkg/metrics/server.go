package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	registry = prometheus.NewRegistry()

	handler = promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: false,
		},
	)
)

func init() {
	registry.MustRegister(
		sessionsGauge,
		queriesCounter,
		healtyGauge,
	)
}

func Handler() http.Handler {
	return handler
}

func StartServer(logger *logrus.Logger, path string, addr string) error {
	logger.WithField("addr", addr).WithField("path", path).Info("starting prometheus metrics server")

	http.Handle(path, handler)
	return http.ListenAndServe(addr, nil)
}
