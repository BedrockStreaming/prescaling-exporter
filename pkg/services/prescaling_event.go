package services

import (
	"context"
	"errors"
	prescalingv1 "github.com/bedrockstreaming/prescaling-exporter/generated/client/clientset/versioned/typed/prescaling.bedrock.tech/v1"
	v1 "github.com/bedrockstreaming/prescaling-exporter/pkg/apis/prescaling.bedrock.tech/v1"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/utils"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
	"time"
)

type PrescalingEventOutput struct {
	Name string `json:"name"`
	v1.PrescalingEventSpec
}

type PrescalingEventListOutput []PrescalingEventOutput

type IPrescalingEventService interface {
	Create(prescalingevent *v1.PrescalingEvent) (*PrescalingEventOutput, error)
	Delete(name string) error
	List() (*PrescalingEventListOutput, error)
	Get(name string) (*PrescalingEventOutput, error)
	Update(prescalingevent *v1.PrescalingEvent) (*PrescalingEventOutput, error)
	Current() (*PrescalingEventOutput, error)
	Clean(retentionDays int) error
	GetClock() clock.PassiveClock
}

type PrescalingEventService struct {
	prescalingEventRepository prescalingv1.PrescalingEventInterface
	clock                     clock.PassiveClock
}

func NewEventService(prescalingEventRepository prescalingv1.PrescalingEventInterface, clock clock.PassiveClock) IPrescalingEventService {
	return &PrescalingEventService{
		prescalingEventRepository: prescalingEventRepository,
		clock:                     clock,
	}
}

func (e *PrescalingEventService) Create(prescalingevent *v1.PrescalingEvent) (*PrescalingEventOutput, error) {
	create, err := e.prescalingEventRepository.Create(context.Background(), prescalingevent, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	r := &PrescalingEventOutput{
		Name: create.Name,
		PrescalingEventSpec: v1.PrescalingEventSpec{
			Date:        create.Spec.Date,
			StartTime:   create.Spec.StartTime,
			EndTime:     create.Spec.EndTime,
			Multiplier:  create.Spec.Multiplier,
			Description: create.Spec.Description,
		},
	}

	return r, nil
}

func (e *PrescalingEventService) Delete(name string) error {
	err := e.prescalingEventRepository.Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (e *PrescalingEventService) List() (*PrescalingEventListOutput, error) {
	events, err := e.prescalingEventRepository.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make(PrescalingEventListOutput, 0)
	for _, item := range events.Items {
		result = append(result, PrescalingEventOutput{
			Name: item.Name,
			PrescalingEventSpec: v1.PrescalingEventSpec{
				Date:        item.Spec.Date,
				StartTime:   item.Spec.StartTime,
				EndTime:     item.Spec.EndTime,
				Multiplier:  item.Spec.Multiplier,
				Description: item.Spec.Description,
			},
		})
	}

	return &result, nil
}

func (e *PrescalingEventService) Get(name string) (*PrescalingEventOutput, error) {
	event, err := e.prescalingEventRepository.Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	result := &PrescalingEventOutput{
		Name: event.Name,
		PrescalingEventSpec: v1.PrescalingEventSpec{
			Date:        event.Spec.Date,
			StartTime:   event.Spec.StartTime,
			EndTime:     event.Spec.EndTime,
			Multiplier:  event.Spec.Multiplier,
			Description: event.Spec.Description,
		},
	}

	return result, nil
}

func (e *PrescalingEventService) Update(prescalingevent *v1.PrescalingEvent) (*PrescalingEventOutput, error) {
	find, err := e.prescalingEventRepository.Get(context.Background(), prescalingevent.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	prescalingevent.ResourceVersion = find.ResourceVersion

	event, err := e.prescalingEventRepository.Update(context.Background(), prescalingevent, metav1.UpdateOptions{})

	if err != nil {
		return nil, err
	}

	response := &PrescalingEventOutput{
		Name: event.Name,
		PrescalingEventSpec: v1.PrescalingEventSpec{
			Date:        event.Spec.Date,
			StartTime:   event.Spec.StartTime,
			EndTime:     event.Spec.EndTime,
			Multiplier:  event.Spec.Multiplier,
			Description: event.Spec.Description,
		},
	}

	return response, nil
}

func (e *PrescalingEventService) Current() (*PrescalingEventOutput, error) {
	events, err := e.prescalingEventRepository.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	filtered := filter(events.Items, func(event v1.PrescalingEvent) bool {
		date := e.clock.Now().Format("2006-01-02")

		start, _ := utils.SetTime(event.Spec.StartTime, e.clock.Now())
		end, _ := utils.SetTime(event.Spec.EndTime, e.clock.Now())

		return event.Spec.Date == date && utils.InRangeTime(start, end, e.clock.Now())
	})

	if len(filtered) == 0 {
		return nil, errors.New("no events found")
	}

	event := filtered[0]

	response := &PrescalingEventOutput{
		Name: event.Name,
		PrescalingEventSpec: v1.PrescalingEventSpec{
			Date:        event.Spec.Date,
			StartTime:   event.Spec.StartTime,
			EndTime:     event.Spec.EndTime,
			Multiplier:  event.Spec.Multiplier,
			Description: event.Spec.Description,
		},
	}
	return response, nil
}

func (e *PrescalingEventService) Clean(retentionDays int) error {

	if retentionDays < 2 {
		return errors.New("retention days must be > 1")
	}

	events, err := e.prescalingEventRepository.List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, item := range events.Items {
		itemDate, err := time.Parse("2006-01-02", item.Spec.Date)

		if err != nil {
			log.Error(err)
			continue
		}

		if utils.DaysBetweenDates(e.clock.Now(), itemDate) < retentionDays {
			continue
		}

		err = e.prescalingEventRepository.Delete(context.Background(), item.Name, metav1.DeleteOptions{})
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("prescaling events %s has been deleted", item.Name)
	}
	return nil
}

func (e *PrescalingEventService) GetClock() clock.PassiveClock {
	return e.clock
}

func filter(data []v1.PrescalingEvent, f func(v1.PrescalingEvent) bool) []v1.PrescalingEvent {
	fltd := make([]v1.PrescalingEvent, 0)
	for _, e := range data {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}
