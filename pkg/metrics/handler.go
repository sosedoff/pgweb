package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	startTime   time.Time
	promHandler http.Handler
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	uptimeGauge.Set(time.Since(h.startTime).Seconds())

	h.promHandler.ServeHTTP(rw, req)
}

func NewHandler() http.Handler {
	return Handler{
		startTime:   time.Now(),
		promHandler: promhttp.Handler(),
	}
}
