package helper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/STARRY-S/kube-helper-mcp/pkg/internal/types"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *KubeHelper) listDeployment(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Apps.Deployment().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Workloads = append(result.Workloads, types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listDaemonSet(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Apps.DaemonSet().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Workloads = append(result.Workloads, types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listStatefulSet(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Apps.StatefulSet().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Workloads = append(result.Workloads, types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listJob(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Batch.Job().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Workloads = append(result.Workloads, types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listCronJob(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Batch.CronJob().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Workloads = append(result.Workloads, types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listPod(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Core.Pod().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Workloads = append(result.Workloads, types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) ListWorkload(
	workload string,
	ns string,
	labels []string,
	limit int64,
) (string, error) {
	opts := metav1.ListOptions{
		Limit: limit,
	}
	if len(labels) > 0 {
		opts.LabelSelector = strings.Join(labels, ",")
	}

	var listFunc func(string, metav1.ListOptions) (*listResult, error)
	switch strings.TrimSuffix(strings.ToLower(workload), "s") {
	case "deployment":
		listFunc = h.listDeployment
	case "daemonset":
		listFunc = h.listDaemonSet
	case "statefulset":
		listFunc = h.listStatefulSet
	case "job":
		listFunc = h.listJob
	case "cronjob":
		listFunc = h.listCronJob
	case "pod", "":
		listFunc = h.listPod
	default:
		return "", fmt.Errorf("unsupported workload type: %s", workload)
	}

	ns = strings.ToLower(ns)
	switch ns {
	case "*":
		ns = ""
	}

	result, err := listFunc(ns, opts)
	if err != nil {
		return "", fmt.Errorf("failed to list %v: %w", workload, err)
	}
	return result.String(), nil
}

func (h *KubeHelper) ListNamespace(
	limit int64,
) (string, error) {
	opts := metav1.ListOptions{
		Limit: limit,
	}

	list, err := h.wctx.Core.Namespace().List(opts)
	if err != nil {
		return "", err
	}
	result := []string{}
	for _, list := range list.Items {
		result = append(result, list.Name)
	}
	return strings.Join(result, ","), nil
}

func (h *KubeHelper) listWorkloadHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	_ = ctx
	workload, ok := request.Params.Arguments["workload"].(string)
	if !ok {
		return nil, errors.New("workload not provided")
	}
	namespace, _ := request.Params.Arguments["namespace"].(string)
	labels, _ := request.Params.Arguments["labels"].([]any)
	limit, _ := request.Params.Arguments["limit"].(float64)

	s := make([]string, 0, len(labels))
	for _, label := range labels {
		if str, ok := label.(string); ok {
			s = append(s, str)
		}
	}
	result, err := h.ListWorkload(workload, namespace, s, int64(limit))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(result), nil
}

func (h *KubeHelper) listNamespaceHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	_ = ctx
	limit, _ := request.Params.Arguments["limit"].(float64)
	result, err := h.ListNamespace(int64(limit))
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(result), nil
}
