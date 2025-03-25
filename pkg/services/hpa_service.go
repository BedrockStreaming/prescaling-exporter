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

// HPAInfo represents information about an HPA with its minimums
type HPAInfo struct {
	Namespace             string `json:"namespace"`
	Name                  string `json:"name"`
	TargetKind            string `json:"targetKind"`
	TargetName            string `json:"targetName"`
	HpaMinReplicas        int    `json:"hpaMinReplicas"`
	MinReplicasAnnotation int    `json:"minReplicasAnnotation"`
	MeetsMinimum          bool   `json:"meetsMinimum"`
}

// HPAMinimumStatus represents the status of HPA minimums in the cluster
type HPAMinimumStatus struct {
	AllHPAsMeetMinimum bool      `json:"allHPAsMeetMinimum"`
	HPAInfos           []HPAInfo `json:"hpaInfos,omitempty"`
}

// PlatformScalingStatus represents the scaling status of the platform
type PlatformScalingStatus struct {
	IsPlatformScaled bool `json:"isPlatformScaled"`
}

// IHPAService defines the interface for the HPA service
type IHPAService interface {
	CheckHPAMinimums() (*HPAMinimumStatus, error)
	GetPlatformScalingStatus() (*PlatformScalingStatus, error)
}

// HPAService implements IHPAService
type HPAService struct {
	clientset kubernetes.Interface
}

// NewHPAService creates a new instance of HPAService
func NewHPAService(clientset kubernetes.Interface) IHPAService {
	return &HPAService{
		clientset: clientset,
	}
}

// CheckHPAMinimums checks if all HPAs in the cluster meet their minimums
func (h *HPAService) CheckHPAMinimums() (*HPAMinimumStatus, error) {
	ctx := context.Background()
	result := &HPAMinimumStatus{
		AllHPAsMeetMinimum: true,
		HPAInfos:           []HPAInfo{},
	}

	// Get all HPAs
	hpaList, err := h.clientset.AutoscalingV2().HorizontalPodAutoscalers("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving HPAs: %w", err)
	}

	// If no HPAs, everything is OK
	if len(hpaList.Items) == 0 {
		return result, nil
	}

	// Check each HPA with the minimum replicas annotation
	for _, hpa := range hpaList.Items {
		// Check if the minimum replicas annotation is present
		minReplicasStr, hasAnnotation := hpa.ObjectMeta.Annotations[config.Config.AnnotationMinReplicas]
		if !hasAnnotation || minReplicasStr == "" {
			continue // Ignore HPAs without the minimum annotation
		}

		// Convert the annotation to an integer
		minReplicasAnnotation, err := strconv.Atoi(minReplicasStr)
		if err != nil {
			log.Warnf("Unable to convert annotation %s to integer for HPA %s/%s: %v",
				config.Config.AnnotationMinReplicas, hpa.Namespace, hpa.Name, err)
			continue // Ignore HPAs with invalid annotation
		}

		// Get the minimum number of replicas configured in the HPA
		var hpaMinReplicas int
		if hpa.Spec.MinReplicas != nil {
			hpaMinReplicas = int(*hpa.Spec.MinReplicas)
		} else {
			hpaMinReplicas = 1 // By default, if MinReplicas is not specified, it's 1
		}

		// Check if the minimum configured in the HPA is >= the minimum required by the annotation
		meetsMinimum := hpaMinReplicas >= minReplicasAnnotation

		// Create the HPA info
		hpaInfo := HPAInfo{
			Namespace:             hpa.Namespace,
			Name:                  hpa.Name,
			TargetKind:            hpa.Spec.ScaleTargetRef.Kind,
			TargetName:            hpa.Spec.ScaleTargetRef.Name,
			HpaMinReplicas:        hpaMinReplicas,
			MinReplicasAnnotation: minReplicasAnnotation,
			MeetsMinimum:          meetsMinimum,
		}

		// Add the info to the list
		result.HPAInfos = append(result.HPAInfos, hpaInfo)

		// Update the global status
		if !meetsMinimum {
			result.AllHPAsMeetMinimum = false
			log.Warnf("HPA %s/%s: %d configured replicas < %d minimum required (target: %s/%s)",
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
