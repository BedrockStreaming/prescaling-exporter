package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/BedrockStreaming/prescaling-exporter/pkg/services"
	"github.com/BedrockStreaming/prescaling-exporter/pkg/utils"
)

// IHPAHandlers defines the interface for HPA handlers
type IHPAHandlers interface {
	CheckHPAMinimums(w http.ResponseWriter, r *http.Request)
	GetPlatformScalingStatus(w http.ResponseWriter, r *http.Request)
}

// HPAHandlers implements IHPAHandlers
type HPAHandlers struct {
	hpaService services.IHPAService
}

// NewHPAHandlers creates a new instance of HPAHandlers
func NewHPAHandlers(hpaService services.IHPAService) IHPAHandlers {
	return &HPAHandlers{
		hpaService: hpaService,
	}
}

// CheckHPAMinimums
// @Summary      Check if HPAs meet their required minimums
// @Description  Checks if all HPAs in the cluster have at least the minimum number of running pods
// @Tags         cluster
// @Accept       json
// @Produce      json
// @Success      200  {object}  services.HPAMinimumStatus
// @Router       /api/v1/hpas [get]
func (h *HPAHandlers) CheckHPAMinimums(w http.ResponseWriter, r *http.Request) {
	status, err := h.hpaService.CheckHPAMinimums()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteResponse(w, http.StatusOK, status)
}

// GetPlatformScalingStatus
// @Summary      Check the platform scaling status
// @Description  Checks if the platform is correctly scaled and returns HPA details
// @Tags         cluster
// @Accept       json
// @Produce      json
// @Success      200  {object}  services.PlatformScalingStatus
// @Router       /api/v1/hpas/check [get]
func (h *HPAHandlers) GetPlatformScalingStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.hpaService.GetPlatformScalingStatus()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the platform is not correctly scaled, return a 200 OK status
	if !status.IsPlatformScaled {
		utils.WriteResponse(w, http.StatusOK, status)
		return
	}

	utils.WriteResponse(w, http.StatusOK, status)
}
