package types

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Workload struct {
	Kind            string `json:"kind,omitempty"`
	metav1.TypeMeta `json:",inline"`
	Metadata        ObjectMeta     `json:"metadata"`
	Status          WorkloadStatus `json:"status"`
}

func NewWorkload(a any) *Workload {
	b, _ := json.Marshal(a)
	w := &Workload{}
	err := json.Unmarshal(b, w)
	if err != nil {
		logrus.Fatal(err)
	}
	return w
}

func (w *Workload) String() string {
	return w.Metadata.Name
}
