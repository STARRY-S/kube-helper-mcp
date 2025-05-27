package types

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Node struct {
	Kind            string `json:"kind,omitempty"`
	metav1.TypeMeta `json:",inline"`
	Metadata        ObjectMeta `json:"metadata"`
	Spec            NodeSpec   `json:"spec"`
	Status          NodeStatus `json:"status"`
}

type NodeSpec struct {
	PodCIDRs []string `json:"podCIDRs,omitempty"`
}

type NodeStatus struct {
	Phase      string                `json:"phase,omitempty"`
	Conditions []Condition           `json:"conditions,omitempty"`
	Addresses  []NodeAddress         `json:"addresses,omitempty"`
	NodeInfo   corev1.NodeSystemInfo `json:"nodeInfo,omitempty"`
}

type NodeAddress struct {
	Type    string `json:"type,omitempty"`
	Address string `json:"address,omitempty"`
}

func NewNode(a any) *Node {
	b, _ := json.Marshal(a)
	w := &Node{}
	err := json.Unmarshal(b, w)
	if err != nil {
		logrus.Fatal(err)
		// logrus.Warnf("failed to unmarshal JSON: %v", err)
	}
	return w
}

func (w *Node) String() string {
	return w.Metadata.Name
}
