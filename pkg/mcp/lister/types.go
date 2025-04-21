package lister

import (
	"github.com/STARRY-S/kube-helper-mcp/pkg/internal/types"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
)

type listResult struct {
	results []*types.Workload `json:"results,omitempty"`
}

func (r *listResult) String() string {
	return utils.Print(r.results)
}
