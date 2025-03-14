package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/BedrockStreaming/prescaling-exporter/pkg/services"
	"github.com/BedrockStreaming/prescaling-exporter/pkg/utils"
)

// IHPAHandlers définit l'interface pour les handlers HPA
type IHPAHandlers interface {
	CheckHPAMinimums(w http.ResponseWriter, r *http.Request)
	GetPlatformScalingStatus(w http.ResponseWriter, r *http.Request)
}

// HPAHandlers implémente IHPAHandlers
type HPAHandlers struct {
	hpaService services.IHPAService
}

// NewHPAHandlers crée une nouvelle instance de HPAHandlers
func NewHPAHandlers(hpaService services.IHPAService) IHPAHandlers {
	return &HPAHandlers{
		hpaService: hpaService,
	}
}

// CheckHPAMinimums
// @Summary      Vérifier si les HPA respectent leurs minimums requis
// @Description  Vérifie si tous les HPA du cluster ont au moins le nombre minimum de pods en cours d'exécution
// @Tags         cluster
// @Accept       json
// @Produce      json
// @Success      200  {object}  services.HPAMinimumStatus
// @Router       /api/v1/hpa-check [get]
func (h *HPAHandlers) CheckHPAMinimums(w http.ResponseWriter, r *http.Request) {
	status, err := h.hpaService.CheckHPAMinimums()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Si tous les HPA ne respectent pas leurs minimums, renvoyer un code 200 OK
	if !status.AllHPAsMeetMinimum {
		utils.WriteResponse(w, http.StatusOK, status)
		return
	}

	utils.WriteResponse(w, http.StatusOK, status)
}

// GetPlatformScalingStatus
// @Summary      Vérifier l'état de scaling de la plateforme
// @Description  Vérifie si la plateforme est correctement scalée et renvoie les détails des HPA
// @Tags         cluster
// @Accept       json
// @Produce      json
// @Success      200  {object}  services.PlatformScalingStatus
// @Router       /api/v1/platform-scaling [get]
func (h *HPAHandlers) GetPlatformScalingStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.hpaService.GetPlatformScalingStatus()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Si la plateforme n'est pas correctement scalée, renvoyer un code 200 OK
	if !status.IsPlatformScaled {
		utils.WriteResponse(w, http.StatusOK, status)
		return
	}

	utils.WriteResponse(w, http.StatusOK, status)
}
