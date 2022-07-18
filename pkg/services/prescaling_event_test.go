package services

import (
	"reflect"
	"testing"
	"time"

	"k8s.io/utils/clock"
	testclock "k8s.io/utils/clock/testing"

	prescalingv1 "github.com/bedrockstreaming/prescaling-exporter/generated/client/clientset/versioned/typed/prescaling.bedrock.tech/v1"
	"github.com/bedrockstreaming/prescaling-exporter/generated/client/clientset/versioned/typed/prescaling.bedrock.tech/v1/fake"
)

func TestNewEventService(t *testing.T) {

	fakeClock := testclock.NewFakeClock(time.Date(2022, time.March, 2, 21, 0, 0, 0, time.UTC))
	fakePrescalingEvents := &fake.FakePrescalingEvents{}

	type args struct {
		prescalingEventRepository prescalingv1.PrescalingEventInterface
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

//func TestPrescalingEventService_Clean(t *testing.T) {
//	fakeClock := testclock.NewFakeClock(time.Date(2022, time.March, 2, 21, 0, 0, 0, time.UTC))
//	fakeEventRepository := fake.FakePrescalingEvents{}
//	fakeEventRepository.Fake.ReactionChain = []testingFake.Reactor{
//	}
//	actionlist := testingFake.ListActionImpl{
//
//	}
//	type fields struct {
//		prescalingEventRepository prescalingv1.PrescalingEventInterface
//		clock                     clock.PassiveClock
//	}
//	type args struct {
//		retentionDays int
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr error
//	}{
//		{
//			name: "test-1",
//			args: args{
//				retentionDays: 1,
//			},
//			wantErr: errors.New("retention days must be > 1"),
//
//		},
//		{
//			name: "test-2",
//			args: args{
//				retentionDays: 2,
//			},
//			fields: fields{
//				fakeEventRepository,
//				fakeClock,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &PrescalingEventService{
//				prescalingEventRepository: tt.fields.prescalingEventRepository,
//				clock:                     tt.fields.clock,
//			}
//			if err := e.Clean(tt.args.retentionDays); err.Error() != tt.wantErr.Error() {
//				t.Errorf("Clean() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

//func TestPrescalingEventService_Create(t *testing.T) {
//	type fields struct {
//		prescalingEventRepository prescalingv1.PrescalingEventInterface
//		clock                     clock.PassiveClock
//	}
//	type args struct {
//		prescalingevent *v1.PrescalingEvent
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *PrescalingEventOutput
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &PrescalingEventService{
//				prescalingEventRepository: tt.fields.prescalingEventRepository,
//				clock:                     tt.fields.clock,
//			}
//			got, err := e.Create(tt.args.prescalingevent)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Create() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPrescalingEventService_Current(t *testing.T) {
//	type fields struct {
//		prescalingEventRepository prescalingv1.PrescalingEventInterface
//		clock                     clock.PassiveClock
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    *PrescalingEventOutput
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &PrescalingEventService{
//				prescalingEventRepository: tt.fields.prescalingEventRepository,
//				clock:                     tt.fields.clock,
//			}
//			got, err := e.Current()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Current() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Current() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPrescalingEventService_Delete(t *testing.T) {
//	type fields struct {
//		prescalingEventRepository prescalingv1.PrescalingEventInterface
//		clock                     clock.PassiveClock
//	}
//	type args struct {
//		name string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &PrescalingEventService{
//				prescalingEventRepository: tt.fields.prescalingEventRepository,
//				clock:                     tt.fields.clock,
//			}
//			if err := e.Delete(tt.args.name); (err != nil) != tt.wantErr {
//				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestPrescalingEventService_Get(t *testing.T) {
//	type fields struct {
//		prescalingEventRepository prescalingv1.PrescalingEventInterface
//		clock                     clock.PassiveClock
//	}
//	type args struct {
//		name string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *PrescalingEventOutput
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &PrescalingEventService{
//				prescalingEventRepository: tt.fields.prescalingEventRepository,
//				clock:                     tt.fields.clock,
//			}
//			got, err := e.Get(tt.args.name)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Get() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPrescalingEventService_GetClock(t *testing.T) {
//	type fields struct {
//		prescalingEventRepository prescalingv1.PrescalingEventInterface
//		clock                     clock.PassiveClock
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   clock.PassiveClock
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &PrescalingEventService{
//				prescalingEventRepository: tt.fields.prescalingEventRepository,
//				clock:                     tt.fields.clock,
//			}
//			if got := e.GetClock(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetClock() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPrescalingEventService_List(t *testing.T) {
//	type fields struct {
//		prescalingEventRepository prescalingv1.PrescalingEventInterface
//		clock                     clock.PassiveClock
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    *PrescalingEventListOutput
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &PrescalingEventService{
//				prescalingEventRepository: tt.fields.prescalingEventRepository,
//				clock:                     tt.fields.clock,
//			}
//			got, err := e.List()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("List() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPrescalingEventService_Update(t *testing.T) {
//	type fields struct {
//		prescalingEventRepository prescalingv1.PrescalingEventInterface
//		clock                     clock.PassiveClock
//	}
//	type args struct {
//		prescalingevent *v1
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *PrescalingEventOutput
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &PrescalingEventService{
//				prescalingEventRepository: tt.fields.prescalingEventRepository,
//				clock:                     tt.fields.clock,
//			}
//			got, err := e.Update(tt.args.prescalingevent)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Update() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_filter(t *testing.T) {
//	type args struct {
//		data []v1.PrescalingEvent
//		f    func(v1.PrescalingEvent) bool
//	}
//	tests := []struct {
//		name string
//		args args
//		want []v1.PrescalingEvent
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := filter(tt.args.data, tt.args.f); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("filter() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
