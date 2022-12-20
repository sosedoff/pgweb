package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func Handler() http.Handler {
	return promhttp.Handler()
}

func StartServer(logger *logrus.Logger, path string, addr string) error {
	logger.WithField("addr", addr).WithField("path", path).Info("starting prometheus metrics server")

	http.Handle(path, Handler())
	return http.ListenAndServe(addr, nil)
}
