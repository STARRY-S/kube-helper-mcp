package k8s

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *KubeHelper) getDeployment(ns string, n string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Apps.Deployment().Get(ns, n, opts)
}

func (h *KubeHelper) getDaemonSet(ns string, n string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Apps.DaemonSet().Get(ns, n, opts)
}
func (h *KubeHelper) getStatefulSet(ns string, n string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Apps.StatefulSet().Get(ns, n, opts)
}
func (h *KubeHelper) getPod(ns string, n string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Core.Pod().Get(ns, n, opts)
}
func (h *KubeHelper) getJob(ns string, n string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Batch.Job().Get(ns, n, opts)
}
func (h *KubeHelper) getCronJob(ns string, n string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Batch.CronJob().Get(ns, n, opts)
}
func (h *KubeHelper) getNamespace(ns string, _ string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Core.Namespace().Get(ns, opts)
}
func (h *KubeHelper) getNode(_ string, n string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Core.Node().Get(n, opts)
}
func (h *KubeHelper) getService(ns string, n string, opts metav1.GetOptions) (metav1.Object, error) {
	return h.Wctx.Core.Service().Get(ns, n, opts)
}

func (h *KubeHelper) GetResource(
	workload string,
	name string,
	ns string,
) (string, error) {
	opts := metav1.GetOptions{}
	var getFunc func(ns string, n string, o metav1.GetOptions) (metav1.Object, error)
	switch strings.TrimSuffix(strings.ToLower(workload), "s") {
	case "deployment":
		getFunc = h.getDeployment
	case "daemonset":
		getFunc = h.getDaemonSet
	case "statefulset":
		getFunc = h.getStatefulSet
	case "job":
		getFunc = h.getJob
	case "cronjob":
		getFunc = h.getCronJob
	case "pod", "":
		getFunc = h.getPod
	case "namespace":
		getFunc = h.getNamespace
	case "node":
		getFunc = h.getNode
	case "service":
		getFunc = h.getService
	default:
		return "", fmt.Errorf("unsupported workload type: %s", workload)
	}

	ns = strings.ToLower(ns)
	switch ns {
	case "*":
		ns = ""
	}

	result, err := getFunc(ns, name, opts)
	if err != nil {
		return "", fmt.Errorf("failed to list %v: %w", workload, err)
	}
	result.SetManagedFields(nil)
	return utils.Print(result), nil
}

func (h *KubeHelper) getResourceHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	_ = ctx
	arguments := request.GetArguments()
	if arguments == nil {
		return nil, fmt.Errorf("failed to get request arguments")
	}
	workload, ok := arguments["resource"].(string)
	if !ok {
		return nil, errors.New("resource not provided")
	}
	name, ok := arguments["name"].(string)
	if !ok {
		return nil, errors.New("resource name not provided, use list_resources to list resources")
	}
	namespace, _ := arguments["namespace"].(string)

	result, err := h.GetResource(workload, name, namespace)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	logrus.Debugf("handle the getResource Handler")
	return mcp.NewToolResultText(result), nil
}
