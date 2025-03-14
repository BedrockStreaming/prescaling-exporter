package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/BedrockStreaming/prescaling-exporter/docs"
	"github.com/BedrockStreaming/prescaling-exporter/pkg/config"
	"github.com/BedrockStreaming/prescaling-exporter/pkg/handlers"
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
	hpaHandlers   handlers.IHPAHandlers
}

func NewServer(statusHandler handlers.IStatusHandlers, eventHandlers handlers.IEventHandlers, hpaHandlers handlers.IHPAHandlers) IServer {
	return &Server{
		statusHandler: statusHandler,
		eventHandlers: eventHandlers,
		hpaHandlers:   hpaHandlers,
	}
}

func (s *Server) Initialize() error {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/status", s.statusHandler.Index)

	// Routes pour les événements de prescaling
	apiv1 := router.PathPrefix("/api/v1/events").Subrouter()
	apiv1.HandleFunc("/", s.eventHandlers.List).Methods(http.MethodGet)
	apiv1.HandleFunc("/", s.eventHandlers.Create).Methods(http.MethodPost)
	apiv1.HandleFunc("/current", s.eventHandlers.Current).Methods(http.MethodGet)
	apiv1.HandleFunc("/{name}", s.eventHandlers.Get).Methods(http.MethodGet)
	apiv1.HandleFunc("/{name}", s.eventHandlers.Update).Methods(http.MethodPut)
	apiv1.HandleFunc("/{name}", s.eventHandlers.Delete).Methods(http.MethodDelete)

	// Route pour la vérification des minimums des HPA
	router.HandleFunc("/api/v1/hpas", s.hpaHandlers.CheckHPAMinimums).Methods(http.MethodGet)

	// Route for checking the platform scaling status
	router.HandleFunc("/api/v1/hpas/check", s.hpaHandlers.GetPlatformScalingStatus).Methods(http.MethodGet)

	log.Info("Listen on port: ", config.Config.Port)

	return http.ListenAndServe(":"+config.Config.Port, router)
}
