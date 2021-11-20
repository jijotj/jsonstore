package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"jsonstore/pkg/lib"
	"jsonstore/pkg/prometheus"
)

func Recovery(p *prometheus.Prometheus, serviceName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					p.PanicCounter().WithLabelValues(serviceName).Inc()
					zap.S().Errorf("panic recovery: %v", err)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					lib.WriteResponseJSON(w, map[string]string{"error": "There was an internal server error"})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
