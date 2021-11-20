package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	promlib "github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"jsonstore/pkg/middleware"
	"jsonstore/pkg/prometheus"
)

func TestHTTPMetrics(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/resource/status", nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, "hello")
		require.NoError(t, err, "Unexpected write error")
	})

	registry := promlib.NewRegistry()
	httpMetricsMiddleware := middleware.HTTPMetrics(prometheus.NewPrometheus(registry), "jsonstore")
	httpMetricsMiddleware(nextHandler).ServeHTTP(w, r)

	metricFamilies, err := registry.Gather()
	require.NoError(t, err, "Incorrect gather error")

	t.Run("response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, w.Code, "Incorrect HTTP status code")
		assert.Equal(t, "hello", w.Body.String(), "Incorrect body")
	})

	t.Run("latency metrics", func(t *testing.T) {
		metrics := findMetricsByName(t, metricFamilies, "jsonstore_http_handler_latency")

		labels := metrics.GetLabel()
		require.Len(t, labels, 2, "Incorrect number of labels")
		assert.Equal(t, "path", labels[0].GetName(), "Incorrect label name")
		assert.Equal(t, "/resource", labels[0].GetValue(), "Incorrect label value")
		assert.Equal(t, "service", labels[1].GetName(), "Incorrect label name")
		assert.Equal(t, "jsonstore", labels[1].GetValue(), "Incorrect label value")
	})

	t.Run("status code metrics", func(t *testing.T) {
		metrics := findMetricsByName(t, metricFamilies, "jsonstore_http_status_code")

		labels := metrics.GetLabel()
		require.Len(t, labels, 4, "Incorrect number of labels")
		assert.Equal(t, "code", labels[0].GetName(), "Incorrect label name")
		assert.Equal(t, "200", labels[0].GetValue(), "Incorrect label value")
		assert.Equal(t, "error", labels[1].GetName(), "Incorrect label name")
		assert.Empty(t, labels[1].GetValue(), "Incorrect label value")
		assert.Equal(t, "path", labels[2].GetName(), "Incorrect label name")
		assert.Equal(t, "/resource", labels[2].GetValue(), "Incorrect label value")
		assert.Equal(t, "service", labels[3].GetName(), "Incorrect label name")
		assert.Equal(t, "jsonstore", labels[3].GetValue(), "Incorrect label value")

		assert.Equal(t, 1.0, metrics.GetCounter().GetValue(), "Incorrect count for HTTP status code 203")
	})
}

func findMetricsByName(t *testing.T, metrics []*dto.MetricFamily, name string) *dto.Metric {
	for _, m := range metrics {
		if m.GetName() == name {
			metrics := m.GetMetric()
			require.Len(t, metrics, 1, "Incorrect number of metrics")
			return metrics[0]
		}
	}

	return nil
}
