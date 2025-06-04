package k8sgpt

import (
	"context"

	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type clusterResult struct {
	Name    string `json:"name"`
	Details string `json:"details"`
	Kind    string `json:"kind"`
}

func (h *Helper) getCheckResultsHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	_ = ctx
	_ = request
	result, err := h.GetCheckClusterResults()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func (h *Helper) GetCheckClusterResults() (string, error) {
	results, err := h.Wctx.K8sGPT.Result().List(defaultK8sGPTNamespace, metav1.ListOptions{
		// LabelSelector: "k8sgpts.k8sgpt.ai/name=" + defaultK8sGPTName,
	})
	if err != nil {
		return "", err
	}
	if results == nil || len(results.Items) == 0 {
		return "No abnormality has been found in the cluster yet.", nil
	}

	res := make([]clusterResult, 0, len(results.Items))
	for _, item := range results.Items {
		res = append(res, clusterResult{
			Name:    item.Spec.Name,
			Details: item.Spec.Details,
			Kind:    item.Spec.Kind,
		})
	}
	return utils.Print(res), nil
}
