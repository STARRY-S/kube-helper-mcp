package types

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Event struct {
	Kind            string     `json:"kind,omitempty"`
	Metadata        ObjectMeta `json:"metadata"`
	metav1.TypeMeta `json:",inline"`

	Reason  string             `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	Message string             `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	Source  corev1.EventSource `json:"source,omitempty" protobuf:"bytes,5,opt,name=source"`
	Count   int32              `json:"count,omitempty" protobuf:"varint,8,opt,name=count"`
	Type    string             `json:"type,omitempty" protobuf:"bytes,9,opt,name=type"`
}

func NewEvent(a any) *Event {
	b, _ := json.Marshal(a)
	w := &Event{}
	err := json.Unmarshal(b, w)
	if err != nil {
		logrus.Fatal(err)
	}
	return w
}

func (w *Event) String() string {
	return w.Metadata.Name
}
