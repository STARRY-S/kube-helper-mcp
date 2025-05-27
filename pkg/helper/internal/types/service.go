package types

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {
	Kind            string `json:"kind,omitempty"`
	metav1.TypeMeta `json:",inline"`
	Metadata        ObjectMeta    `json:"metadata"`
	Status          ServiceStatus `json:"status"`
}

type ServiceSpec struct {
	Ports     []corev1.ServicePort `json:"ports,omitempty"`
	ClusterIP string               `json:"clusterIP,omitempty"`
	Type      string               `json:"type,omitempty"`
}

type ServiceStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

func NewService(a any) *Service {
	b, _ := json.Marshal(a)
	w := &Service{}
	err := json.Unmarshal(b, w)
	if err != nil {
		logrus.Fatal(err)
	}
	return w
}

func (w *Service) String() string {
	return w.Metadata.Name
}
