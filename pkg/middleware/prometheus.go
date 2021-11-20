package middleware

import (
	"net/http"
	"strconv"
	"strings"

	promlib "github.com/prometheus/client_golang/prometheus"

	"jsonstore/pkg/lib"
	"jsonstore/pkg/prometheus"
)

func HTTPMetrics(p *prometheus.Prometheus, serviceName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := "/" + strings.Split(r.URL.Path, "/")[1]
			timer := promlib.NewTimer(p.HTTPHandlerLatencyHistogram().WithLabelValues(serviceName, path))
			defer timer.ObserveDuration()

			ww := lib.NewRecordingWriter(w)
			next.ServeHTTP(ww, r)

			p.HTTPStatusCodeCounter().WithLabelValues(serviceName, path, strconv.Itoa(ww.Status), ww.Err.Code).Inc()
		})
	}
}
