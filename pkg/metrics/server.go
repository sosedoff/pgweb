package metrics

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func StartServer(logger *logrus.Logger, path string, addr string) error {
	logger.WithField("addr", addr).WithField("path", path).Info("starting prometheus metrics server")

	http.Handle(path, NewHandler())
	return http.ListenAndServe(addr, nil)
}
