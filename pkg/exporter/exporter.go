package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/BedrockStreaming/prescaling-exporter/pkg/prescaling"
	"github.com/BedrockStreaming/prescaling-exporter/pkg/utils"
)

type prescalingCollector struct {
	prescaleMetrics *prometheus.Desc
	minMetrics      *prometheus.Desc
	prescaling      prescaling.IPrescaling
}

func NewPrescalingCollector(p prescaling.IPrescaling) prometheus.Collector {
	return &prescalingCollector{
		prescaleMetrics: prometheus.NewDesc(
			"prescale_metric",
			"Number used for prescale application",
			[]string{"project", "deployment", "namespace"},
			nil,
		), minMetrics: prometheus.NewDesc(
			"min_replica",
			"Number of pod desired for prescale",
			[]string{"project", "deployment", "namespace"},
			nil,
		),
		prescaling: p,
	}
}

func (collector *prescalingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.prescaleMetrics
	ch <- collector.minMetrics
}

func (collector *prescalingCollector) Collect(ch chan<- prometheus.Metric) {
	log.Info("Collect")
	hpaList := collector.prescaling.GetHpa()
	if len(hpaList) == 0 {
		log.Error("error - no prescaling hpa configuration found")
		return
	}

	currentPrescalingEvent, err := collector.prescaling.GetEventService().Current()

	now := collector.prescaling.GetEventService().GetClock().Now()
	for _, hpa := range hpaList {
		multiplier := 1
		if err == nil && currentPrescalingEvent.StartTime != "" && currentPrescalingEvent.EndTime != "" {
			hpa.Start, _ = utils.SetTime(currentPrescalingEvent.StartTime, now)
			hpa.End, _ = utils.SetTime(currentPrescalingEvent.EndTime, now)
			multiplier = currentPrescalingEvent.Multiplier
		}

		collector.addDataToMetrics(ch, multiplier, hpa)
	}

}

func (collector *prescalingCollector) addDataToMetrics(ch chan<- prometheus.Metric, multiplier int, hpa prescaling.Hpa) {
	eventInRangeTime := utils.InRangeTime(hpa.Start, hpa.End, collector.prescaling.GetEventService().GetClock().Now())
	desiredScalingType := prescaling.DesiredScaling(eventInRangeTime, multiplier, hpa.Replica, hpa.CurrentReplicas)
	prescaleMetric := prometheus.MustNewConstMetric(collector.prescaleMetrics, prometheus.GaugeValue, float64(desiredScalingType), hpa.Project, hpa.Deployment, hpa.Namespace)
	minMetric := prometheus.MustNewConstMetric(collector.minMetrics, prometheus.GaugeValue, float64(hpa.Replica), hpa.Project, hpa.Deployment, hpa.Namespace)
	ch <- prescaleMetric
	ch <- minMetric
}
