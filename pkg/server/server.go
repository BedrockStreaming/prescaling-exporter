package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/bedrockstreaming/prescaling-exporter/docs"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/config"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/handlers"
)

// @title        Prescaling API
// @version      1.0.0
// @description  This API was built with FastAPI to deal with prescaling recordings in CRD

type IServer interface {
	Initialize() error
}

type Server struct {
	statusHandler handlers.IStatusHandlers
	eventHandlers handlers.IEventHandlers
}

func NewServer(statusHandler handlers.IStatusHandlers, eventHandlers handlers.IEventHandlers) IServer {
	return &Server{
		statusHandler: statusHandler,
		eventHandlers: eventHandlers,
	}
}

func (s *Server) Initialize() error {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9101/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/status", s.statusHandler.Index)

	apiv1 := router.PathPrefix("/api/v1/events").Subrouter()
	apiv1.HandleFunc("/", s.eventHandlers.List).Methods(http.MethodGet)
	apiv1.HandleFunc("/", s.eventHandlers.Create).Methods(http.MethodPost)
	apiv1.HandleFunc("/current", s.eventHandlers.Current).Methods(http.MethodGet)
	apiv1.HandleFunc("/{name}", s.eventHandlers.Get).Methods(http.MethodGet)
	apiv1.HandleFunc("/{name}", s.eventHandlers.Update).Methods(http.MethodPut)
	apiv1.HandleFunc("/{name}", s.eventHandlers.Delete).Methods(http.MethodDelete)

	log.Info("Listen on port: ", config.Config.Port)

	return http.ListenAndServe(":"+config.Config.Port, router)
}
