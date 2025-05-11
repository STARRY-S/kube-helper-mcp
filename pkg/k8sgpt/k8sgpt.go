package k8sgpt

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/mark3labs/mcp-go/mcp"
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
	defaultK8sGPTVersion    = "latest"
	defaultK8sGPTSecretName = "k8sgpt-openai-api-key"
	defaultK8sGPTSecretKey  = "api-key"
)

func (h *Helper) checkClusterHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	_ = ctx
	_ = request
	result, err := h.CheckCluster()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func (h *Helper) CheckCluster() (string, error) {
	result, err := h.wctx.K8sGPT.K8sGPT().Get(defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return "", fmt.Errorf("failed to get k8sgpt: %w", err)
		}
		k := newK8sGPT()
		_, err := h.wctx.K8sGPT.K8sGPT().Create(k)
		if err != nil {
			return "", fmt.Errorf("failed to create k8sgpt: %w", err)
		}
		return "Successfully created the K8sGPT resource, please check the results in a few minutes.", nil
	}

	logrus.Infof("K8sGPT [%s/%s] already exists.", result.Namespace, result.Name)
	if !needUpdateK8sGPT(result) {
		return "Done, the K8sGPT resource already exists, please check the results in a few minutes.", nil
	}
	logrus.Infof("Changes detected, updating the K8sGPT resource.")

	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err = h.wctx.K8sGPT.K8sGPT().Get(
			defaultK8sGPTNamespace, defaultK8sGPTName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get k8sgpt: %w", err)
		}
		result = result.DeepCopy()
		updateK8sGPT(result)
		_, err = h.wctx.K8sGPT.K8sGPT().Update(result)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", fmt.Errorf("failed to update k8sgpt: %w", err)
	}
	return "Successfully updated the k8sgpt resource, please check the new results in a few minutes.", nil
}

var (
	// TODO: Make these values configurable
	defaultK8sGPTSpec = k8sgptv1alpha1.K8sGPTSpec{
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
		NoCache:    false,
		Repository: defaultK8sGPTRepository,
		Version:    defaultK8sGPTVersion,
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
	return !reflect.DeepEqual(result.Spec, defaultK8sGPTSpec)
}

func updateK8sGPT(result *k8sgptv1alpha1.K8sGPT) {
	if result == nil {
		return
	}
	result.Spec = defaultK8sGPTSpec
}
