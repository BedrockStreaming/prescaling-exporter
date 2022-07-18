package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/utils/clock"

	"github.com/bedrockstreaming/prescaling-exporter/generated/client/clientset/versioned"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/config"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/exporter"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/handlers"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/k8s"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/prescaling"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/server"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/services"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serv = &cobra.Command{
	Use:   "serv",
	Short: "Run prescaling server",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := k8s.NewClient()
		if err != nil {
			log.Error("error - k8s client login error: ", err)
			return err
		}

		prescalingClient, err := versioned.NewForConfig(client.Config)
		if err != nil {
			log.Error("error - k8s prescaling client login error: ", err)
			return err
		}

		prescalingEvents := prescalingClient.PrescalingV1().PrescalingEvents(config.Config.Namespace)
		time := clock.RealClock{}

		eventService := services.NewEventService(prescalingEvents, time)
		eventHandler := handlers.NewEventHandlers(eventService)
		statusHandler := handlers.NewStatusHandlers()

		collector := exporter.NewPrescalingCollector(
			prescaling.NewPrescaling(client, eventService),
		)

		prometheus.MustRegister(collector)

		return server.NewServer(statusHandler, eventHandler).Initialize()
	},
}

func init() {
	rootCmd.AddCommand(serv)
}
