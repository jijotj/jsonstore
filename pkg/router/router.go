package router

import (
	"net/http"

	"github.com/gorilla/mux"
	promlib "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"jsonstore/pkg/handler"
	"jsonstore/pkg/middleware"
	"jsonstore/pkg/prometheus"
	"jsonstore/pkg/service"
)

const (
	metricsPath = "/metrics"
	healthPath  = "/health"
	configsPath = "/configs"
)

type Context struct {
	Manager service.Manager
}

func (ctx Context) New() http.Handler {
	promRegistry := promlib.NewRegistry()
	prom := prometheus.NewPrometheus(promRegistry)

	metrics := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{})

	middlewares := []middleware.Middleware{
		middleware.HTTPMetrics(prom, "jsonstore"),
		middleware.Recovery(prom, "jsonstore"),
	}

	router := mux.NewRouter()

	//Future extension: An auth middleware can be added here to protect the APIs
	router.Handle(healthPath, middleware.Wrap(handler.Health(), middlewares...)).Methods(http.MethodGet)
	router.Handle(metricsPath, middleware.Wrap(metrics, middlewares...)).Methods(http.MethodGet)

	router.Handle(configsPath+"/search", middleware.Wrap(handler.SearchConfigs(ctx.Manager),
		middlewares...)).Methods(http.MethodGet)
	router.Handle(configsPath+"/{name}", middleware.Wrap(handler.GetConfig(ctx.Manager),
		middlewares...)).Methods(http.MethodGet)
	router.Handle(configsPath, middleware.Wrap(handler.GetAllConfigs(ctx.Manager),
		middlewares...)).Methods(http.MethodGet)

	router.Handle(configsPath, middleware.Wrap(handler.CreateConfig(ctx.Manager),
		middlewares...)).Methods(http.MethodPost)
	router.Handle(configsPath+"/{name}", middleware.Wrap(handler.UpdateConfig(ctx.Manager),
		middlewares...)).Methods(http.MethodPut, http.MethodPatch)

	router.Handle(configsPath+"/{name}", middleware.Wrap(handler.DeleteConfig(ctx.Manager),
		middlewares...)).Methods(http.MethodDelete)

	return router
}
