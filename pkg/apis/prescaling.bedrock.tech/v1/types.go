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
	Date        string `json:"date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Multiplier  int    `json:"multiplier"`
	Description string `json:"description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PrescalingEventList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PrescalingEvent `json:"items"`
}
