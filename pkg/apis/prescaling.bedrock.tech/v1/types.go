package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PrescalingEvent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PrescalingEventSpec `json:"spec"`
}

type PrescalingEventSpec struct {
	Date        string `json:"date" example:"2022-05-25"`
	StartTime   string `json:"start_time" example:"20:00:00"`
	EndTime     string `json:"end_time" example:"23:59:59"`
	Multiplier  int    `json:"multiplier" example:"2"`
	Description string `json:"description" example:"a good description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PrescalingEventList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PrescalingEvent `json:"items"`
}
