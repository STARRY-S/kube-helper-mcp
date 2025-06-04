package k8sgpt

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (h *Helper) remediateClusterHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	var progressToken mcp.ProgressToken
	if request.Params.Meta != nil {
		progressToken = request.Params.Meta.ProgressToken
	}

	// Get the current session from context
	session := server.ClientSessionFromContext(ctx)
	server := server.ServerFromContext(ctx)
	if session == nil {
		return mcp.NewToolResultError("No active session"), nil
	}
	if server == nil {
		return mcp.NewToolResultError("Failed to get server"), nil
	}
	logrus.Infof("Trigger remediateClusterHandler %v", session.SessionID())

	err := h.RemediateCluster()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	logrus.Debugf("Done trigger remediate cluster %v", session.SessionID())

	if progressToken != nil {
		if err := server.SendNotificationToClient(
			ctx,
			"notifications/message",
			map[string]any{
				"message": "Successfully triggered K8sGPT RemediateCluster actions",
			},
		); err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to send notification: %w", err).Error()), nil
		}
		// TODO: check and wait for cluster remediate results...
	}

	result, err := h.GetMutationResult()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func (h *Helper) RemediateCluster() error {
	result, err := h.Wctx.K8sGPT.K8sGPT().Get(
		defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get k8sgpt resources: %w", err)
		}
		return nil
	}

	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err = h.Wctx.K8sGPT.K8sGPT().Get(
			defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get k8sgpt: %w", err)
		}
		result = result.DeepCopy()
		result.Spec.AI.AutoRemediation.Enabled = true
		_, err = h.Wctx.K8sGPT.K8sGPT().Update(result)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to update k8sgpt: %w", err)
	}
	return nil
}

func (h *Helper) DisableRemediateCluster() error {
	result, err := h.Wctx.K8sGPT.K8sGPT().Get(
		defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get k8sgpt resources: %w", err)
		}
		return nil
	}

	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err = h.Wctx.K8sGPT.K8sGPT().Get(
			defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get k8sgpt: %w", err)
		}
		result = result.DeepCopy()
		result.Spec.AI.AutoRemediation.Enabled = false
		_, err = h.Wctx.K8sGPT.K8sGPT().Update(result)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to update k8sgpt: %w", err)
	}
	return nil
}
