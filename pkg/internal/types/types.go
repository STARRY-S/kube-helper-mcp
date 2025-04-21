package types

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Workload struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Status    Status `json:"status,omitempty"`
}

type Status struct {
	Replicas            int32       `json:"replicas,omitempty"`
	ReadyReplicas       int32       `json:"readyReplicas,omitempty"`
	AvailableReplicas   int32       `json:"availableReplicas,omitempty"`
	UnavailableReplicas int32       `json:"unavailableReplicas,omitempty"`
	Conditions          []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	Type    string             `json:"type"`
	Status  v1.ConditionStatus `json:"status"`
	Reason  string             `json:"reason,omitempty"`
	Message string             `json:"message,omitempty"`
}

func NewWorkload(a any) *Workload {
	b, _ := json.Marshal(a)
	w := &Workload{}
	err := json.Unmarshal(b, w)
	if err != nil {
		logrus.Warnf("failed to unmarshal JSON: %v", err)
	}
	return w
}

func UnmarsalWorkload(b []byte) *Workload {
	w := &Workload{}
	err := json.Unmarshal(b, w)
	if err != nil {
		logrus.Warnf("failed to unmarshal JSON: %v", err)
	}
	return w
}
