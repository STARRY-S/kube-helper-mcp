package helper

import (
	"encoding/json"

	"github.com/STARRY-S/kube-helper-mcp/pkg/internal/types"
)

type listResult struct {
	Workloads []*types.Workload `json:"workloads"`
}

func (r *listResult) AddWorkload(w *types.Workload) {
	r.Workloads = append(r.Workloads, w)
}

func (r *listResult) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
