package prescaling

import (
	"github.com/bedrockstreaming/prescaling-exporter/pkg/k8s"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/services"
)

type IPrescaling interface {
	GetHpa() []Hpa
	GetEventService() services.IPrescalingEventService
}

type Prescaling struct {
	client                 *k8s.Client
	prescalingEventService services.IPrescalingEventService
}

func (p Prescaling) GetEventService() services.IPrescalingEventService {
	return p.prescalingEventService
}

func NewPrescaling(client *k8s.Client, eventService services.IPrescalingEventService) IPrescaling {
	return &Prescaling{
		client:                 client,
		prescalingEventService: eventService,
	}
}
