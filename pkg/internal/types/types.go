package types

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type ObjectMeta struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type WorkloadStatus struct {
	Replicas          int32       `json:"replicas,omitempty"`
	AvailableReplicas int32       `json:"availableReplicas,omitempty"`
	Conditions        []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type Resource struct {
	ObjectMeta
}

func NewResource(a any) *Resource {
	b, _ := json.Marshal(a)
	w := &Resource{}
	err := json.Unmarshal(b, w)
	if err != nil {
		logrus.Fatal(err)
	}
	return w
}

func (w *Resource) String() string {
	return w.Name
}

func (w *Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(w)
}
