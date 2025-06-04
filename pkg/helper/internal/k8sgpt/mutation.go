package k8sgpt

import (
	"context"

	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"

	k8sgptv1alpha1 "github.com/k8sgpt-ai/k8sgpt-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	noRemediateResultFound = `The remediate action has been successfully triggered and no mutation result yet,
which indicates that the cluster is already healthy or is currently remediating resources.
If necessary, check the cluster mutation results again after a few minutes.`
)

type mutationResult struct {
	ResourceRef corev1.ObjectReference        `json:"resource,omitempty"`
	Status      k8sgptv1alpha1.MutationStatus `json:"status"`
}

func (h *Helper) getRemediateResultHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	_ = ctx
	_ = request
	result, err := h.GetMutationResult()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func (h *Helper) GetMutationResult() (string, error) {
	results, err := h.Wctx.K8sGPT.Mutation().List(defaultK8sGPTNamespace, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	if results == nil || len(results.Items) == 0 {
		return noRemediateResultFound, nil
	}

	res := make([]mutationResult, 0, len(results.Items))
	for _, item := range results.Items {
		res = append(res, mutationResult{
			ResourceRef: item.Spec.ResourceRef,
			Status:      item.Status,
		})
	}
	return utils.Print(res), nil
}
