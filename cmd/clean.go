package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/utils/clock"

	"github.com/BedrockStreaming/prescaling-exporter/generated/client/clientset/versioned"
	"github.com/BedrockStreaming/prescaling-exporter/pkg/config"
	"github.com/BedrockStreaming/prescaling-exporter/pkg/k8s"
	"github.com/BedrockStreaming/prescaling-exporter/pkg/services"
)

var retention int

var clean = &cobra.Command{
	Use:   "clean",
	Short: "Run clean CRD prescaling events",
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

		return eventService.Clean(retention)
	},
}

func init() {
	rootCmd.AddCommand(clean)
	clean.PersistentFlags().IntVarP(&retention, "retention", "r", 10, "retention days for CRD")
}
