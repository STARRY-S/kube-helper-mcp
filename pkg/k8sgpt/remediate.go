package k8sgpt

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (h *Helper) remediateClusterHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	_ = ctx
	_ = request
	result, err := h.RemediateCluster()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func (h *Helper) RemediateCluster() (string, error) {
	result, err := h.wctx.K8sGPT.K8sGPT().Get(
		defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return "", fmt.Errorf("failed to get k8sgpt: %w", err)
		}
		return "The K8sGPT self-check action not triggered, use 'check_cluster' to trigger the cluster self-check before remediate.", nil
	}

	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err = h.wctx.K8sGPT.K8sGPT().Get(
			defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get k8sgpt: %w", err)
		}
		result = result.DeepCopy()
		result.Spec.AI.AutoRemediation.Enabled = true
		_, err = h.wctx.K8sGPT.K8sGPT().Update(result)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", fmt.Errorf("failed to update k8sgpt: %w", err)
	}
	return "Successfully enabled the k8sgpt auto remediation, please check the remediate results in few minutes.", nil
}
