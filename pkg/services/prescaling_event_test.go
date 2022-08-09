package services

import (
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
	testclock "k8s.io/utils/clock/testing"

	fakeclient "github.com/BedrockStreaming/prescaling-exporter/generated/client/clientset/versioned/fake"
	prescalingv1client "github.com/BedrockStreaming/prescaling-exporter/generated/client/clientset/versioned/typed/prescaling.bedrock.tech/v1"
	"github.com/BedrockStreaming/prescaling-exporter/generated/client/clientset/versioned/typed/prescaling.bedrock.tech/v1/fake"
	v1 "github.com/BedrockStreaming/prescaling-exporter/pkg/apis/prescaling.bedrock.tech/v1"
)

func TestNewEventService(t *testing.T) {

	fakeClock := testclock.NewFakeClock(time.Date(2022, time.March, 2, 21, 0, 0, 0, time.Now().Local().Location()))
	fakePrescalingEvents := &fake.FakePrescalingEvents{}

	type args struct {
		prescalingEventRepository prescalingv1client.PrescalingEventInterface
		clock                     clock.PassiveClock
	}
	tests := []struct {
		name string
		args args
		want IPrescalingEventService
	}{
		{
			name: "test-1",
			args: args{
				prescalingEventRepository: fakePrescalingEvents,
				clock:                     fakeClock,
			},
			want: &PrescalingEventService{
				fakePrescalingEvents,
				fakeClock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventService(tt.args.prescalingEventRepository, tt.args.clock); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrescalingEventService_Current(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Paris")
	fakeClock := testclock.NewFakeClock(time.Date(2022, time.July, 2, 23, 0, 1, 0, loc))

	cs := fakeclient.NewSimpleClientset(
		&v1.PrescalingEvent{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PrescalingEvent",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "project-event-1",
				Namespace: "default",
			},
			Spec: v1.PrescalingEventSpec{
				Date:        "2022-07-02",
				StartTime:   "20:00:00",
				EndTime:     "21:59:59",
				Multiplier:  0,
				Description: "",
			},
		},
		&v1.PrescalingEvent{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PrescalingEvent",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "project-event-2",
				Namespace: "default",
			},
			Spec: v1.PrescalingEventSpec{
				Date:        "2022-07-02",
				StartTime:   "22:00:00",
				EndTime:     "23:59:59",
				Multiplier:  0,
				Description: "",
			},
		},
	)
	prescaling := cs.PrescalingV1().PrescalingEvents("default")

	type fields struct {
		prescalingEventRepository prescalingv1client.PrescalingEventInterface
		clock                     clock.PassiveClock
	}
	tests := []struct {
		name    string
		fields  fields
		want    *PrescalingEventOutput
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				prescalingEventRepository: prescaling,
				clock:                     fakeClock,
			},
			want: &PrescalingEventOutput{
				Name: "project-event-2",
				PrescalingEventSpec: v1.PrescalingEventSpec{
					Date:        "2022-07-02",
					StartTime:   "22:00:00",
					EndTime:     "23:59:59",
					Multiplier:  0,
					Description: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &PrescalingEventService{
				prescalingEventRepository: tt.fields.prescalingEventRepository,
				clock:                     tt.fields.clock,
			}
			got, err := e.Current()
			if (err != nil) != tt.wantErr {
				t.Errorf("Current() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Current() got = %v, want %v", got, tt.want)
			}
		})
	}
}
