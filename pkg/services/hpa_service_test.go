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
	// Initialize configuration for tests
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
	// Create a fake clientset with test data
	fakeClientset := fake.NewSimpleClientset()

	// Define test values
	minReplicas := int32(2)

	// Create a test HPA with the minimum replicas annotation
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			Annotations: map[string]string{
				config.Config.AnnotationMinReplicas: "2", // Set the minimum via annotation
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

	// Add objects to the clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test HPA: %v", err)
	}

	// Create the HPA service
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Run the test
	status, err := hpaService.CheckHPAMinimums()
	if err != nil {
		t.Fatalf("Error calling CheckHPAMinimums: %v", err)
	}

	// Check that all HPAs meet their minimums
	if !status.AllHPAsMeetMinimum {
		t.Errorf("AllHPAsMeetMinimum should be true because the HPA has 3 replicas and the minimum is 2")
	}

	// Check that the HPA information is correct
	if len(status.HPAInfos) != 1 {
		t.Errorf("HPAInfos should contain 1 element, but contains %d", len(status.HPAInfos))
	} else {
		hpaInfo := status.HPAInfos[0]
		if hpaInfo.Namespace != "default" || hpaInfo.Name != "test-hpa" || hpaInfo.TargetKind != "Deployment" ||
			hpaInfo.TargetName != "test-deployment" || hpaInfo.HpaMinReplicas != 3 || hpaInfo.MinReplicasAnnotation != 2 || !hpaInfo.MeetsMinimum {
			t.Errorf("Incorrect HPAInfo: %+v", hpaInfo)
		}
	}
}

// TestHPAService_CheckHPAMinimums_NoHPAs tests the case where there are no HPAs in the cluster
func TestHPAService_CheckHPAMinimums_NoHPAs(t *testing.T) {
	// Create a fake clientset without HPAs
	fakeClientset := fake.NewSimpleClientset()

	// Create the HPA service
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Run the test
	status, err := hpaService.CheckHPAMinimums()
	if err != nil {
		t.Fatalf("Error calling CheckHPAMinimums: %v", err)
	}

	// Check that AllHPAsMeetMinimum is true when there are no HPAs
	if !status.AllHPAsMeetMinimum {
		t.Errorf("AllHPAsMeetMinimum should be true when there are no HPAs")
	}

	// Check that HPAInfos is empty
	if len(status.HPAInfos) != 0 {
		t.Errorf("HPAInfos should be empty, but contains %d elements", len(status.HPAInfos))
	}
}

// TestHPAService_CheckHPAMinimums_HPABelowMinimum tests the case where an HPA is below the required minimum
func TestHPAService_CheckHPAMinimums_HPABelowMinimum(t *testing.T) {
	// Create a fake clientset with test data
	fakeClientset := fake.NewSimpleClientset()

	// Create a test HPA with the minimum replicas annotation
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			Annotations: map[string]string{
				config.Config.AnnotationMinReplicas: "3", // Set the minimum via annotation (3 replicas)
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: "test-deployment",
			},
		},
		Status: autoscalingv2.HorizontalPodAutoscalerStatus{
			CurrentReplicas: 2, // Only 2 replicas, below the minimum of 3
			DesiredReplicas: 2,
		},
	}

	// Add the HPA to the clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test HPA: %v", err)
	}

	// Create the HPA service
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Run the test
	status, err := hpaService.CheckHPAMinimums()
	if err != nil {
		t.Fatalf("Error calling CheckHPAMinimums: %v", err)
	}

	// Check that AllHPAsMeetMinimum is false when an HPA is below the minimum
	if status.AllHPAsMeetMinimum {
		t.Errorf("AllHPAsMeetMinimum should be false when an HPA is below the minimum")
	}

	// Check that the HPA information is correct
	if len(status.HPAInfos) != 1 {
		t.Errorf("HPAInfos should contain 1 element, but contains %d", len(status.HPAInfos))
	} else {
		hpaInfo := status.HPAInfos[0]
		if hpaInfo.Namespace != "default" || hpaInfo.Name != "test-hpa" || hpaInfo.TargetKind != "Deployment" ||
			hpaInfo.TargetName != "test-deployment" || hpaInfo.HpaMinReplicas != 2 || hpaInfo.MinReplicasAnnotation != 3 || hpaInfo.MeetsMinimum {
			t.Errorf("Incorrect HPAInfo: %+v", hpaInfo)
		}
	}
}

// TestHPAService_CheckHPAMinimums_NoAnnotation tests the case where an HPA has no minimum annotation
func TestHPAService_CheckHPAMinimums_NoAnnotation(t *testing.T) {
	// Create a fake clientset with test data
	fakeClientset := fake.NewSimpleClientset()

	// Create a test HPA without minimum annotation
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			// No minimum annotation
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

	// Add the HPA to the clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test HPA: %v", err)
	}

	// Create the HPA service
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Run the test
	status, err := hpaService.CheckHPAMinimums()
	if err != nil {
		t.Fatalf("Error calling CheckHPAMinimums: %v", err)
	}

	// Check that AllHPAsMeetMinimum is true because HPAs without annotation are ignored
	if !status.AllHPAsMeetMinimum {
		t.Errorf("AllHPAsMeetMinimum should be true when HPAs have no minimum annotation")
	}

	// Check that HPAInfos is empty
	if len(status.HPAInfos) != 0 {
		t.Errorf("HPAInfos should be empty, but contains %d elements", len(status.HPAInfos))
	}
}

// TestHPAService_GetPlatformScalingStatus tests the GetPlatformScalingStatus method
func TestHPAService_GetPlatformScalingStatus(t *testing.T) {
	// Create a fake clientset with test data
	fakeClientset := fake.NewSimpleClientset()

	// Create a test HPA with the minimum replicas annotation
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			Annotations: map[string]string{
				config.Config.AnnotationMinReplicas: "2", // Set the minimum via annotation
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: "test-deployment",
			},
		},
		Status: autoscalingv2.HorizontalPodAutoscalerStatus{
			CurrentReplicas: 3, // Above the minimum of 2
			DesiredReplicas: 3,
		},
	}

	// Add the HPA to the clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test HPA: %v", err)
	}

	// Create the HPA service
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Run the test
	status, err := hpaService.GetPlatformScalingStatus()
	if err != nil {
		t.Fatalf("Error calling GetPlatformScalingStatus: %v", err)
	}

	// Check that IsPlatformScaled is true when all HPAs meet their minimums
	if !status.IsPlatformScaled {
		t.Errorf("IsPlatformScaled should be true when all HPAs meet their minimums")
	}
}

// TestHPAService_GetPlatformScalingStatus_HPABelowMinimum tests the case where an HPA is below the required minimum
func TestHPAService_GetPlatformScalingStatus_HPABelowMinimum(t *testing.T) {
	// Create a fake clientset with test data
	fakeClientset := fake.NewSimpleClientset()

	// Create a test HPA with the minimum replicas annotation
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hpa",
			Namespace: "default",
			Annotations: map[string]string{
				config.Config.AnnotationMinReplicas: "3", // Set the minimum via annotation (3 replicas)
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: "test-deployment",
			},
		},
		Status: autoscalingv2.HorizontalPodAutoscalerStatus{
			CurrentReplicas: 2, // Only 2 replicas, below the minimum of 3
			DesiredReplicas: 2,
		},
	}

	// Add the HPA to the clientset
	_, err := fakeClientset.AutoscalingV2().HorizontalPodAutoscalers("default").Create(context.Background(), hpa, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test HPA: %v", err)
	}

	// Create the HPA service
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Run the test
	status, err := hpaService.GetPlatformScalingStatus()
	if err != nil {
		t.Fatalf("Error calling GetPlatformScalingStatus: %v", err)
	}

	// Check that IsPlatformScaled is false when an HPA is below the minimum
	if status.IsPlatformScaled {
		t.Errorf("IsPlatformScaled should be false when an HPA is below the minimum")
	}
}

// TestHPAService_GetPlatformScalingStatus_NoHPAs tests the case where there are no HPAs in the cluster
func TestHPAService_GetPlatformScalingStatus_NoHPAs(t *testing.T) {
	// Create a fake clientset without HPAs
	fakeClientset := fake.NewSimpleClientset()

	// Create the HPA service
	hpaService := &HPAService{
		clientset: fakeClientset,
	}

	// Run the test
	status, err := hpaService.GetPlatformScalingStatus()
	if err != nil {
		t.Fatalf("Error calling GetPlatformScalingStatus: %v", err)
	}

	// Check that IsPlatformScaled is true when there are no HPAs
	if !status.IsPlatformScaled {
		t.Errorf("IsPlatformScaled should be true when there are no HPAs")
	}
}
