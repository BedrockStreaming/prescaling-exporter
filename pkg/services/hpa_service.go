package services

import (
	"context"
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/BedrockStreaming/prescaling-exporter/pkg/config"
	log "github.com/sirupsen/logrus"
)

// HPAInfo représente les informations d'un HPA avec ses minimums
type HPAInfo struct {
	Namespace       string `json:"namespace"`
	Name            string `json:"name"`
	TargetKind      string `json:"targetKind"`
	TargetName      string `json:"targetName"`
	HpaMinReplicas  int    `json:"hpaMinReplicas"`
	MinReplicasAnnotation int    `json:"minReplicasAnnotation"`
	MeetsMinimum    bool   `json:"meetsMinimum"`
}

// HPAMinimumStatus représente l'état des minimums des HPA dans le cluster
type HPAMinimumStatus struct {
	AllHPAsMeetMinimum bool      `json:"allHPAsMeetMinimum"`
	HPAInfos           []HPAInfo `json:"hpaInfos,omitempty"`
}

// PlatformScalingStatus représente l'état de scaling de la plateforme
type PlatformScalingStatus struct {
	IsPlatformScaled bool `json:"isPlatformScaled"`
}

// IHPAService définit l'interface pour le service HPA
type IHPAService interface {
	CheckHPAMinimums() (*HPAMinimumStatus, error)
	GetPlatformScalingStatus() (*PlatformScalingStatus, error)
}

// HPAService implémente IHPAService
type HPAService struct {
	clientset kubernetes.Interface
}

// NewHPAService crée une nouvelle instance de HPAService
func NewHPAService(clientset kubernetes.Interface) IHPAService {
	return &HPAService{
		clientset: clientset,
	}
}

// CheckHPAMinimums vérifie si tous les HPA du cluster respectent leurs minimums
func (h *HPAService) CheckHPAMinimums() (*HPAMinimumStatus, error) {
	ctx := context.Background()
	result := &HPAMinimumStatus{
		AllHPAsMeetMinimum: true,
		HPAInfos:           []HPAInfo{},
	}

	// Récupérer tous les HPA
	hpaList, err := h.clientset.AutoscalingV2().HorizontalPodAutoscalers("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des HPA: %w", err)
	}

	// Si aucun HPA, tout est OK
	if len(hpaList.Items) == 0 {
		return result, nil
	}

	// Vérifier chaque HPA avec l'annotation de minimum de réplicas
	for _, hpa := range hpaList.Items {
		// Vérifier si l'annotation de minimum de réplicas est présente
		minReplicasStr, hasAnnotation := hpa.ObjectMeta.Annotations[config.Config.AnnotationMinReplicas]
		if !hasAnnotation || minReplicasStr == "" {
			continue // Ignorer les HPA sans annotation de minimum
		}

		// Convertir l'annotation en entier
		minReplicasAnnotation, err := strconv.Atoi(minReplicasStr)
		if err != nil {
			log.Warnf("Impossible de convertir l'annotation %s en entier pour le HPA %s/%s: %v",
				config.Config.AnnotationMinReplicas, hpa.Namespace, hpa.Name, err)
			continue // Ignorer les HPA avec une annotation invalide
		}

		// Récupérer le nombre minimum de réplicas configuré dans le HPA
		var hpaMinReplicas int
		if hpa.Spec.MinReplicas != nil {
			hpaMinReplicas = int(*hpa.Spec.MinReplicas)
		} else {
			hpaMinReplicas = 1 // Par défaut, si MinReplicas n'est pas spécifié, c'est 1
		}

		// Vérifier si le minimum configuré dans le HPA est >= au minimum requis par l'annotation
		meetsMinimum := hpaMinReplicas >= minReplicasAnnotation

		// Créer l'info HPA
		hpaInfo := HPAInfo{
			Namespace:       hpa.Namespace,
			Name:            hpa.Name,
			TargetKind:      hpa.Spec.ScaleTargetRef.Kind,
			TargetName:      hpa.Spec.ScaleTargetRef.Name,
			HpaMinReplicas:  hpaMinReplicas,
			MinReplicasAnnotation: minReplicasAnnotation,
			MeetsMinimum:    meetsMinimum,
		}

		// Ajouter l'info à la liste
		result.HPAInfos = append(result.HPAInfos, hpaInfo)

		// Mettre à jour le statut global
		if !meetsMinimum {
			result.AllHPAsMeetMinimum = false
			log.Warnf("HPA %s/%s: %d réplicas configurés < %d minimum requis (cible: %s/%s)",
				hpa.Namespace, hpa.Name, hpaMinReplicas, minReplicasAnnotation, hpa.Spec.ScaleTargetRef.Kind, hpa.Spec.ScaleTargetRef.Name)
		}
	}

	return result, nil
}

func (h *HPAService) GetPlatformScalingStatus() (*PlatformScalingStatus, error) {

	hpaInfo, err := h.CheckHPAMinimums()
	if err != nil {
		return nil, err
	}

	for _, hpa := range hpaInfo.HPAInfos {
		if !hpa.MeetsMinimum {
			return &PlatformScalingStatus{
				IsPlatformScaled: false,
			}, nil
		}
	}

	return &PlatformScalingStatus{
		IsPlatformScaled: true,
	}, nil
}
