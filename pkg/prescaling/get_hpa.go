package prescaling

import (
	"context"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/config"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (p Prescaling) GetHpa() []Hpa {
	hpaList, err := p.client.Clientset.AutoscalingV2beta2().
		HorizontalPodAutoscalers("").
		List(context.Background(), metav1.ListOptions{})

	if err != nil { // todo: retry ?
		panic(err.Error())
	}

	var preScalingList []Hpa

	now := p.prescalingEventService.GetClock().Now()

	for _, hpa := range hpaList.Items {
		if hpa.ObjectMeta.Annotations[config.Config.AnnotationMinReplicas] != "" {
			replica, err := strconv.Atoi(hpa.ObjectMeta.Annotations[config.Config.AnnotationMinReplicas])
			if err != nil {
				log.Infof("error - cannot convert string to int %s", err)
			}

			preScaling := Hpa{
				Replica:         replica,
				CurrentReplicas: hpa.Status.CurrentReplicas,
				Project:         hpa.ObjectMeta.Labels[config.Config.LabelProject],
				Namespace:       hpa.ObjectMeta.Namespace,
				Deployment:      hpa.Spec.ScaleTargetRef.Name,
			}

			preScaling.Start, err = utils.SetTime(hpa.ObjectMeta.Annotations[config.Config.AnnotationStartTime], now)
			if err != nil {
				log.Infof("error - %s", err)
			}

			preScaling.End, err = utils.SetTime(hpa.ObjectMeta.Annotations[config.Config.AnnotationEndTime], now)
			if err != nil {
				log.Infof("error - %s", err)
			}

			err = preScaling.Validate()
			if err != nil {
				log.Infof("error - %s", err)
			} else {
				preScalingList = append(preScalingList, preScaling)
			}
		}
	}

	return preScalingList
}
