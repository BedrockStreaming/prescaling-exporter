package handlers

import (
	"encoding/json"
	v1 "github.com/bedrockstreaming/prescaling-exporter/pkg/apis/prescaling.bedrock.tech/v1"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/services"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type CreateDTO struct {
	Name string `json:"name" example:"prescaling-event-1"`
	v1.PrescalingEventSpec
}

type UpdateDTO struct {
	Name string `json:"name" example:"prescaling-event-1"`
	v1.PrescalingEventSpec
}

type IEventHandlers interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Current(w http.ResponseWriter, r *http.Request)
}

type EventHandlers struct {
	eventService services.IPrescalingEventService
}

func NewEventHandlers(userService services.IPrescalingEventService) IEventHandlers {
	return &EventHandlers{
		eventService: userService,
	}
}

// Create
// @Summary      Create a prescaling Event
// @Description  Create a prescaling Event
// @Tags         prescalingevent
// @Accept       json
// @Produce      json
// @Param 		 data body CreateDTO true "The Request body"
// @Success      200  {object}  services.PrescalingEventOutput
// @Failure      500  {object}  string
// @Router       /api/v1/events/ [post]
func (e *EventHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var query CreateDTO
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prescalingevent := &v1.PrescalingEvent{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PrescalingEventOutput",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: query.Name,
		},
		Spec: v1.PrescalingEventSpec{
			Date:        query.Date,
			StartTime:   query.StartTime,
			EndTime:     query.EndTime,
			Multiplier:  query.Multiplier,
			Description: query.Description,
		},
	}
	prescalingEventCreate, err := e.eventService.Create(prescalingevent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info("New event created: ", prescalingEventCreate)
	utils.WriteResponse(w, http.StatusCreated, prescalingEventCreate)
}

// Delete
// @Summary      Delete a prescaling Event by name
// @Description  Delete a prescaling Event by name
// @Tags         prescalingevent
// @Accept       json
// @Produce      json
// @Param        name path string true  "event-name-1"
// @Success      200  {object}  nil
// @Failure      404  {object}  string
// @Failure      400  {object}  string
// @Router       /api/v1/events/{name} [delete]
func (e *EventHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	var prescalingEventName = mux.Vars(r)["name"]

	err := e.eventService.Delete(prescalingEventName)
	if err != nil {
		if err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonNotFound {
			log.Info(err)
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof("event %s has been deleted", prescalingEventName)
	utils.WriteResponse(w, http.StatusOK, "OK")
}

// List
// @Summary      List all prescaling Events
// @Description  List all prescaling Events
// @Tags         prescalingevent
// @Accept       json
// @Produce      json
// @Success      200  {object}  services.PrescalingEventListOutput
// @Failure      400  {object}  string
// @Router       /api/v1/events/ [get]
func (e *EventHandlers) List(w http.ResponseWriter, r *http.Request) {
	prescalingEvents, err := e.eventService.List()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	utils.WriteResponse(w, http.StatusOK, prescalingEvents)
}

// Get
// @Summary      Get a prescaling Events by name
// @Description  Get a prescaling Events by name
// @Tags         prescalingevent
// @Accept       json
// @Produce      json
// @Param        name path string true  "event-name-1"
// @Success      200  {object}  services.PrescalingEventOutput
// @Failure      404  {object}  string
// @Router       /api/v1/events/{name} [get]
func (e *EventHandlers) Get(w http.ResponseWriter, r *http.Request) {
	var prescalingEventName = mux.Vars(r)["name"]

	event, err := e.eventService.Get(prescalingEventName)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	utils.WriteResponse(w, http.StatusOK, event)
}

// Update
// @Summary      Update a prescaling Event by name
// @Description  Update a prescaling Event by name
// @Tags         prescalingevent
// @Accept       json
// @Produce      json
// @Param 		 data body UpdateDTO true "The Request body"
// @Param        name path string true  "event-name-1"
// @Success      200  {object}  services.PrescalingEventOutput
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Router       /api/v1/events/{name} [put]
func (e *EventHandlers) Update(w http.ResponseWriter, r *http.Request) {
	var query UpdateDTO

	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	prescalingevent := &v1.PrescalingEvent{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PrescalingEventOutput",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: query.Name,
		},
		Spec: v1.PrescalingEventSpec{
			Date:        query.Date,
			StartTime:   query.StartTime,
			EndTime:     query.EndTime,
			Multiplier:  query.Multiplier,
			Description: query.Description,
		},
	}
	prescalingEventUpdated, err := e.eventService.Update(prescalingevent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof("Event %s updated", prescalingEventUpdated.Name)

	utils.WriteResponse(w, http.StatusOK, prescalingEventUpdated)
}

// Current
// @Summary      Get current prescaling Event
// @Description  Get current prescaling Event
// @Tags         prescalingevent
// @Accept       json
// @Produce      json
// @Success      200  {object}  services.PrescalingEventOutput
// @Failure      204  {object}  string
// @Router       /api/v1/events/current/ [get]
func (e *EventHandlers) Current(w http.ResponseWriter, r *http.Request) {
	currentPrescalingEvent, err := e.eventService.Current()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	utils.WriteResponse(w, http.StatusOK, currentPrescalingEvent)
}
