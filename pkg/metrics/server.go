package metrics

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func StartServer(logger *logrus.Logger, path string, addr string) error {
	logger.WithField("addr", addr).WithField("path", path).Info("starting prometheus metrics server")

	mux := http.NewServeMux()
	mux.Handle(path, NewHandler())

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return server.ListenAndServe()
}
