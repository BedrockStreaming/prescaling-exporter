package services

import (
	"context"
	"reflect"
	"testing"

	"github.com/BedrockStreaming/prescaling-exporter/pkg/config"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func init() {
	// Initialiser la configuration pour les tests
	config.Config.AnnotationMinReplicas = "annotations.scaling.exporter.replica.min"
}

func TestNewHPAService(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()

	type args struct {
		clientset kubernetes.Interface
	}
	tests := []struct {
		name string
		args args
		want IHPAService
	}{
		{
			name: "test-new-hpa-service",
			args: args{
				clientset: fakeClientset,
			},
			want: &HPAService{
				clientset: fakeClientset,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHPAService(tt.args.clientset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHPAService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHPAService_CheckHPAMinimums(t *testing.T) {
	// Créer un clientset factice avec des données de test
	fakeClientset := fake.NewSimpleClientset()

	// Définir les valeurs de test
	minReplicas := int32(2)

	// Créer un HPA de test avec l'annotation de minimum de réplicas
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			Annotations: map[string]string{
				config.Config.AnnotationMinReplicas: "2", // Définir le minimum via l'annotation
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: "test-deployment",
			},
			MinReplicas: &minReplicas,
		},
		Status: autoscalingv2.HorizontalPodAutoscalerStatus{
			CurrentReplicas: 3,
			DesiredReplicas: 3,
		},
	}

	// Ajouter les objets au clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Erreur lors de la création du HPA de test: %v", err)
	}

	// Créer le service HPA
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Exécuter le test
	status, err := hpaService.CheckHPAMinimums()
	if err != nil {
		t.Fatalf("Erreur lors de l'appel à CheckHPAMinimums: %v", err)
	}

	// Vérifier que tous les HPA respectent leurs minimums
	if !status.AllHPAsMeetMinimum {
		t.Errorf("AllHPAsMeetMinimum devrait être true car le HPA a 3 réplicas et le minimum est 2")
	}

	// Vérifier que les informations du HPA sont correctes
	if len(status.HPAInfos) != 1 {
		t.Errorf("HPAInfos devrait contenir 1 élément, mais en contient %d", len(status.HPAInfos))
	} else {
		hpaInfo := status.HPAInfos[0]
		if hpaInfo.Namespace != "default" || hpaInfo.Name != "test-hpa" || hpaInfo.TargetKind != "Deployment" ||
			hpaInfo.TargetName != "test-deployment" || hpaInfo.HpaMinReplicas != 3 || hpaInfo.MinReplicasAnnotation != 2 || !hpaInfo.MeetsMinimum {
			t.Errorf("HPAInfo incorrect: %+v", hpaInfo)
		}
	}
}

// TestHPAService_CheckHPAMinimums_NoHPAs teste le cas où il n'y a pas de HPA dans le cluster
func TestHPAService_CheckHPAMinimums_NoHPAs(t *testing.T) {
	// Créer un clientset factice sans HPA
	fakeClientset := fake.NewSimpleClientset()

	// Créer le service HPA
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Exécuter le test
	status, err := hpaService.CheckHPAMinimums()
	if err != nil {
		t.Fatalf("Erreur lors de l'appel à CheckHPAMinimums: %v", err)
	}

	// Vérifier que AllHPAsMeetMinimum est true quand il n'y a pas de HPA
	if !status.AllHPAsMeetMinimum {
		t.Errorf("AllHPAsMeetMinimum devrait être true quand il n'y a pas de HPA")
	}

	// Vérifier que HPAInfos est vide
	if len(status.HPAInfos) != 0 {
		t.Errorf("HPAInfos devrait être vide, mais contient %d éléments", len(status.HPAInfos))
	}
}

// TestHPAService_CheckHPAMinimums_HPABelowMinimum teste le cas où un HPA est en dessous du minimum requis
func TestHPAService_CheckHPAMinimums_HPABelowMinimum(t *testing.T) {
	// Créer un clientset factice avec des données de test
	fakeClientset := fake.NewSimpleClientset()

	// Créer un HPA de test avec l'annotation de minimum de réplicas
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			Annotations: map[string]string{
				config.Config.AnnotationMinReplicas: "3", // Définir le minimum via l'annotation (3 réplicas)
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: "test-deployment",
			},
		},
		Status: autoscalingv2.HorizontalPodAutoscalerStatus{
			CurrentReplicas: 2, // Seulement 2 réplicas, en dessous du minimum de 3
			DesiredReplicas: 2,
		},
	}

	// Ajouter les objets au clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Erreur lors de la création du HPA de test: %v", err)
	}

	// Créer le service HPA
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Exécuter le test
	status, err := hpaService.CheckHPAMinimums()
	if err != nil {
		t.Fatalf("Erreur lors de l'appel à CheckHPAMinimums: %v", err)
	}

	// Vérifier que AllHPAsMeetMinimum est false quand un HPA est en dessous du minimum
	if status.AllHPAsMeetMinimum {
		t.Errorf("AllHPAsMeetMinimum devrait être false quand un HPA est en dessous du minimum")
	}

	// Vérifier que les informations du HPA sont correctes
	if len(status.HPAInfos) != 1 {
		t.Errorf("HPAInfos devrait contenir 1 élément, mais en contient %d", len(status.HPAInfos))
	} else {
		hpaInfo := status.HPAInfos[0]
		if hpaInfo.Namespace != "default" || hpaInfo.Name != "test-hpa" || hpaInfo.TargetKind != "Deployment" ||
			hpaInfo.TargetName != "test-deployment" || hpaInfo.HpaMinReplicas != 2 || hpaInfo.MinReplicasAnnotation != 3 || hpaInfo.MeetsMinimum {
			t.Errorf("HPAInfo incorrect: %+v", hpaInfo)
		}
	}
}

// TestHPAService_CheckHPAMinimums_NoAnnotation teste le cas où un HPA n'a pas d'annotation de minimum
func TestHPAService_CheckHPAMinimums_NoAnnotation(t *testing.T) {
	// Créer un clientset factice avec des données de test
	fakeClientset := fake.NewSimpleClientset()

	// Créer un HPA de test sans annotation de minimum
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			// Pas d'annotation de minimum
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: "test-deployment",
			},
		},
		Status: autoscalingv2.HorizontalPodAutoscalerStatus{
			CurrentReplicas: 1,
			DesiredReplicas: 1,
		},
	}

	// Ajouter l'HPA au clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Erreur lors de la création du HPA de test: %v", err)
	}

	// Créer le service HPA
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Exécuter le test
	status, err := hpaService.CheckHPAMinimums()
	if err != nil {
		t.Fatalf("Erreur lors de l'appel à CheckHPAMinimums: %v", err)
	}

	// Vérifier que AllHPAsMeetMinimum est true car les HPA sans annotation sont ignorés
	if !status.AllHPAsMeetMinimum {
		t.Errorf("AllHPAsMeetMinimum devrait être true quand les HPA n'ont pas d'annotation de minimum")
	}

	// Vérifier que HPAInfos est vide
	if len(status.HPAInfos) != 0 {
		t.Errorf("HPAInfos devrait être vide, mais contient %d éléments", len(status.HPAInfos))
	}
}

// TestHPAService_GetPlatformScalingStatus teste la méthode GetPlatformScalingStatus
func TestHPAService_GetPlatformScalingStatus(t *testing.T) {
	// Créer un clientset factice avec des données de test
	fakeClientset := fake.NewSimpleClientset()

	// Créer un HPA de test avec l'annotation de minimum de réplicas
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			Annotations: map[string]string{
				config.Config.AnnotationMinReplicas: "2", // Définir le minimum via l'annotation
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: "test-deployment",
			},
		},
		Status: autoscalingv2.HorizontalPodAutoscalerStatus{
			CurrentReplicas: 3, // Au-dessus du minimum de 2
			DesiredReplicas: 3,
		},
	}

	// Ajouter l'HPA au clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Erreur lors de la création du HPA de test: %v", err)
	}

	// Créer le service HPA
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Exécuter le test
	status, err := hpaService.GetPlatformScalingStatus()
	if err != nil {
		t.Fatalf("Erreur lors de l'appel à GetPlatformScalingStatus: %v", err)
	}

	// Vérifier que IsPlatformScaled est true quand tous les HPA respectent leurs minimums
	if !status.IsPlatformScaled {
		t.Errorf("IsPlatformScaled devrait être true quand tous les HPA respectent leurs minimums")
	}
}

// TestHPAService_GetPlatformScalingStatus_HPABelowMinimum teste le cas où un HPA est en dessous du minimum requis
func TestHPAService_GetPlatformScalingStatus_HPABelowMinimum(t *testing.T) {
	// Créer un clientset factice avec des données de test
	fakeClientset := fake.NewSimpleClientset()

	// Créer un HPA de test avec l'annotation de minimum de réplicas
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			Annotations: map[string]string{
				config.Config.AnnotationMinReplicas: "3", // Définir le minimum via l'annotation (3 réplicas)
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: "test-deployment",
			},
		},
		Status: autoscalingv2.HorizontalPodAutoscalerStatus{
			CurrentReplicas: 2, // Seulement 2 réplicas, en dessous du minimum de 3
			DesiredReplicas: 2,
		},
	}

	// Ajouter l'HPA au clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Erreur lors de la création du HPA de test: %v", err)
	}

	// Créer le service HPA
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Exécuter le test
	status, err := hpaService.GetPlatformScalingStatus()
	if err != nil {
		t.Fatalf("Erreur lors de l'appel à GetPlatformScalingStatus: %v", err)
	}

	// Vérifier que IsPlatformScaled est false quand un HPA est en dessous du minimum
	if status.IsPlatformScaled {
		t.Errorf("IsPlatformScaled devrait être false quand un HPA est en dessous du minimum")
	}
}

// TestHPAService_GetPlatformScalingStatus_NoHPAs teste le cas où il n'y a pas de HPA dans le cluster
func TestHPAService_GetPlatformScalingStatus_NoHPAs(t *testing.T) {
	// Créer un clientset factice sans HPA
	fakeClientset := fake.NewSimpleClientset()

	// Créer le service HPA
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Exécuter le test
	status, err := hpaService.GetPlatformScalingStatus()
	if err != nil {
		t.Fatalf("Erreur lors de l'appel à GetPlatformScalingStatus: %v", err)
	}

	// Vérifier que IsPlatformScaled est true quand il n'y a pas de HPA
	if !status.IsPlatformScaled {
		t.Errorf("IsPlatformScaled devrait être true quand il n'y a pas de HPA")
	}
}
