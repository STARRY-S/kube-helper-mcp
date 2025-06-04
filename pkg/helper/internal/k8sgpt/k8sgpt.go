package k8sgpt

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/retry"

	k8sgptv1alpha1 "github.com/k8sgpt-ai/k8sgpt-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultK8sGPTAIBackend  = "openai"
	defaultK8sGPTName       = "k8sgpt-cluster-check"
	defaultK8sGPTNamespace  = "k8sgpt-operator-system"
	defaultK8sGPTRepository = "ghcr.io/k8sgpt-ai/k8sgpt"
	defaultK8sGPTVersion    = "v0.4.17"               // TODO: edit here if needed
	defaultK8sGPTSecretName = "k8sgpt-openai-api-key" // #nosec G101
	defaultK8sGPTSecretKey  = "api-key"
)

func (h *Helper) checkClusterHandler(
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
	logrus.Debugf("Trigger checkClusterHandler %v", session.SessionID())

	err := h.TriggerClusterCheck()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	logrus.Debugf("Done trigger clustercheck %v", session.SessionID())

	if progressToken != nil {
		if err := server.SendNotificationToClient(
			ctx,
			"notifications/message",
			map[string]any{
				"message": "Successfully triggered K8sGPT ClusterCheck actions",
			},
		); err != nil {
			return mcp.NewToolResultError(fmt.Errorf("failed to send notification: %w", err).Error()), nil
		}
		// TODO: check and wait for cluster check results...
	}

	result, err := h.GetCheckClusterResults()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(result), nil
}

func (h *Helper) TriggerClusterCheck() error {
	// Create or update k8sgpt check resource.
	result, err := h.Wctx.K8sGPT.K8sGPT().Get(defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get k8sgpt: %w", err)
		}
		k := newK8sGPT()
		_, err := h.Wctx.K8sGPT.K8sGPT().Create(k)
		if err != nil {
			return fmt.Errorf("failed to create k8sgpt: %w", err)
		}
	} else {
		logrus.Infof("K8sGPT [%s/%s] already exists.", result.Namespace, result.Name)
		if needUpdateK8sGPT(result) {
			logrus.Infof("Changes detected, updating the K8sGPT resource.")
			if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				result, err = h.Wctx.K8sGPT.K8sGPT().Get(
					defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("failed to get k8sgpt: %w", err)
				}
				result = result.DeepCopy()
				updateK8sGPT(result)
				_, err = h.Wctx.K8sGPT.K8sGPT().Update(result)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
				return fmt.Errorf("failed to update k8sgpt: %w", err)
			}
		}
	}

	// Waiting for k8sgpt resources
	return nil
}

var (
	// TODO: Make these values configurable
	defaultK8sGPTSpec = k8sgptv1alpha1.K8sGPTSpec{
		Version:    defaultK8sGPTVersion,
		Repository: defaultK8sGPTRepository,
		AI: &k8sgptv1alpha1.AISpec{
			AutoRemediation: k8sgptv1alpha1.AutoRemediation{
				Enabled:   false,
				Resources: []string{"Pod", "Deployment"}, // TODO: only check pod, deployments currently
			},
			Backend: defaultK8sGPTAIBackend,
			Enabled: true,
			Model:   "gpt-4.1-mini",
			Secret: &k8sgptv1alpha1.SecretRef{
				Name: defaultK8sGPTSecretName,
				Key:  defaultK8sGPTSecretKey,
			},
			ProxyEndpoint: os.Getenv("HTTPS_PROXY"),
		},
		NoCache: false,
		Filters: []string{"Pod", "Deployment"}, // TODO: only check pod, deployments currently
	}
)

func newK8sGPT() *k8sgptv1alpha1.K8sGPT {
	return &k8sgptv1alpha1.K8sGPT{
		TypeMeta: metav1.TypeMeta{
			APIVersion: k8sgptv1alpha1.GroupVersion.String(),
			Kind:       "K8sGPT",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      defaultK8sGPTName,
			Namespace: defaultK8sGPTNamespace,
		},
		Spec: defaultK8sGPTSpec,
	}
}

func needUpdateK8sGPT(result *k8sgptv1alpha1.K8sGPT) bool {
	// TODO:
	return !reflect.DeepEqual(result.Spec, defaultK8sGPTSpec)
}

func updateK8sGPT(result *k8sgptv1alpha1.K8sGPT) {
	if result == nil {
		return
	}
	result.Spec = defaultK8sGPTSpec
}
