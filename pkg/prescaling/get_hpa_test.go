package prescaling

import (
	"testing"
	"time"

	clock "k8s.io/utils/clock/testing"

	"github.com/bedrockstreaming/prescaling-exporter/generated/client/clientset/versioned/fake"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/k8s"
	"github.com/bedrockstreaming/prescaling-exporter/pkg/services"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclientk8s "k8s.io/client-go/kubernetes/fake"
)

var loc = time.Now().Local().Location()

func TestGetHpa(t *testing.T) {
	fakeClock := clock.NewFakeClock(time.Date(2022, time.March, 2, 21, 0, 0, 0, loc))

	var client k8s.Client
	client.Clientset = testclientk8s.NewSimpleClientset(
		&v2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "project-a",
				Namespace: "default",
				Annotations: map[string]string{
					"annotations.scaling.exporter.replica.min": "10",
					"annotations.scaling.exporter.time.start":  "20:00:00",
					"annotations.scaling.exporter.time.end":    "23:00:00",
				},
				Labels: map[string]string{"project": "project-a"},
			},
			Spec: v2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: v2beta2.CrossVersionObjectReference{
					Name: "project-a",
				},
			},
			Status: v2beta2.HorizontalPodAutoscalerStatus{
				CurrentReplicas: 2,
			},
		},
		&v2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "project-b",
				Namespace: "default",
				Annotations: map[string]string{
					"annotations.scaling.exporter.replica.min": "20",
					"annotations.scaling.exporter.time.start":  "18:00:00",
					"annotations.scaling.exporter.time.end":    "23:30:00",
				},
				Labels: map[string]string{"project": "project-b"},
			},
			Spec: v2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: v2beta2.CrossVersionObjectReference{
					Name: "project-b",
				},
			},
			Status: v2beta2.HorizontalPodAutoscalerStatus{
				CurrentReplicas: 2,
			},
		},
	)

	prescalingClient := fake.NewSimpleClientset()
	prescalingEvents := prescalingClient.PrescalingV1().PrescalingEvents("default")
	prescalingEventService := services.NewEventService(prescalingEvents, fakeClock)
	prescaling := NewPrescaling(&client, prescalingEventService)

	expected := []Hpa{
		{
			Replica:         10,
			CurrentReplicas: 2,
			Project:         "project-a",
			Namespace:       "default",
			Deployment:      "project-a",
			Start:           time.Date(2022, time.March, 2, 20, 00, 0, 0, loc),
			End:             time.Date(2022, time.March, 2, 23, 00, 0, 0, loc),
		},
		{
			Replica:         20,
			CurrentReplicas: 2,
			Project:         "project-b",
			Namespace:       "default",
			Deployment:      "project-b",
			Start:           time.Date(2022, time.March, 2, 18, 00, 0, 0, loc),
			End:             time.Date(2022, time.March, 2, 23, 30, 0, 0, loc),
		},
	}

	hpaList := prescaling.GetHpa()
	assert.Equal(t, hpaList, expected, "KO - result is not equal to input")
}

func TestGetHpaError(t *testing.T) {
	fakeClock := clock.NewFakeClock(time.Date(2022, time.March, 2, 21, 0, 0, 0, loc))

	var client k8s.Client
	client.Clientset = testclientk8s.NewSimpleClientset(
		&v2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "project-replica-nil",
				Namespace: "default",
				Annotations: map[string]string{
					"annotations.scaling.exporter.replica.min": "0",
					"annotations.scaling.exporter.time.start":  "20:00:00",
					"annotations.scaling.exporter.time.end":    "23:00:00",
				},
				Labels: map[string]string{"project": "project-replica-nil"},
			},
			Spec: v2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: v2beta2.CrossVersionObjectReference{
					Name: "project-replica-nil",
				},
			},
			Status: v2beta2.HorizontalPodAutoscalerStatus{
				CurrentReplicas: 2,
			},
		},
		&v2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "project-no-start",
				Namespace: "default",
				Annotations: map[string]string{
					"annotations.scaling.exporter.replica.min": "20",
					"annotations.scaling.exporter.time.end":    "23:30:00",
				},
				Labels: map[string]string{"project": "project-no-start"},
			},
			Spec: v2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: v2beta2.CrossVersionObjectReference{
					Name: "project-no-start",
				},
			},
			Status: v2beta2.HorizontalPodAutoscalerStatus{
				CurrentReplicas: 2,
			},
		},
		&v2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "project-no-end",
				Namespace: "default",
				Annotations: map[string]string{
					"annotations.scaling.exporter.replica.min": "20",
					"annotations.scaling.exporter.time.start":  "23:30:00",
				},
				Labels: map[string]string{"project": "project-no-end"},
			},
			Spec: v2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: v2beta2.CrossVersionObjectReference{
					Name: "project-no-end",
				},
			},
			Status: v2beta2.HorizontalPodAutoscalerStatus{
				CurrentReplicas: 2,
			},
		},
		&v2beta2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "project-replica-misconfigured",
				Namespace: "default",
				Annotations: map[string]string{
					"annotations.scaling.exporter.replica.min": "err",
					"annotations.scaling.exporter.time.start":  "20:00:00",
					"annotations.scaling.exporter.time.end":    "23:00:00",
				},
				Labels: map[string]string{"project": "project-replica-misconfigured"},
			},
			Spec: v2beta2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: v2beta2.CrossVersionObjectReference{
					Name: "project-replica-misconfigured",
				},
			},
			Status: v2beta2.HorizontalPodAutoscalerStatus{
				CurrentReplicas: 2,
			},
		},
	)

	prescalingClient := fake.NewSimpleClientset()
	prescalingEvents := prescalingClient.PrescalingV1().PrescalingEvents("default")
	prescalingEventService := services.NewEventService(prescalingEvents, fakeClock)
	precaling := NewPrescaling(&client, prescalingEventService)

	hpaList := precaling.GetHpa()
	assert.Nil(t, hpaList, "KO - result is not nul")
}

func TestCheckAnnotationsKO(t *testing.T) {
	testCases := []struct {
		name       string
		expected   string
		prescaling Hpa
	}{
		{
			name: "OK",
			prescaling: Hpa{
				Replica: 1,
				Start:   time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
				End:     time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
			},
		},
		{
			name:     "KO - Replica is null",
			expected: "annotation replica min is misconfigured",
			prescaling: Hpa{
				Replica: 0,
				Start:   time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
				End:     time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
			},
		},
		{
			name:     "KO - Start is null",
			expected: "annotation time start is misconfigured",
			prescaling: Hpa{
				Replica: 1,
				Start:   time.Time{},
				End:     time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
			},
		},
		{
			name:     "KO - End is null",
			expected: "annotation time start is misconfigured",
			prescaling: Hpa{
				Replica: 1,
				Start:   time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
				End:     time.Time{},
			},
		},
	}

	for _, testCase := range testCases {
		err := testCase.prescaling.Validate()
		if err == nil {
			assert.Equal(t, nil, err, testCase.name)
		} else {
			assert.EqualError(t, err, testCase.expected, testCase.name)
		}
	}
}

func TestCheckAnnotationsOK(t *testing.T) {
	testCases := []struct {
		name       string
		expected   string
		prescaling Hpa
	}{
		{
			name: "OK - checkAnnotation return nil",
			prescaling: Hpa{
				Replica: 1,
				Start:   time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
				End:     time.Date(2022, time.March, 2, 20, 30, 0, 0, loc),
			},
		},
	}

	for _, testCase := range testCases {
		err := testCase.prescaling.Validate()
		assert.Nil(t, err, testCase.name)
	}
}
