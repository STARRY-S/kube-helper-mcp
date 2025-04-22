package types

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Workload struct {
	Kind            string `json:"kind,omitempty"`
	metav1.TypeMeta `json:",inline"`
	Metadata        ObjectMeta `json:"metadata,omitempty"`
	Status          Status     `json:"status,omitempty"`
}

type ObjectMeta struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type Status struct {
	Replicas          int32       `json:"replicas,omitempty"`
	AvailableReplicas int32       `json:"availableReplicas,omitempty"`
	Conditions        []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

func NewWorkload(a any) *Workload {
	b, _ := json.Marshal(a)
	w := &Workload{}
	err := json.Unmarshal(b, w)
	if err != nil {
		logrus.Fatal(err)
		// logrus.Warnf("failed to unmarshal JSON: %v", err)
	}
	return w
}
