package k8s

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/STARRY-S/kube-helper-mcp/pkg/helper/internal/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *KubeHelper) listDeployment(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.Wctx.Apps.Deployment().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Add(types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listDaemonSet(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.Wctx.Apps.DaemonSet().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Add(types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listStatefulSet(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.Wctx.Apps.StatefulSet().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Add(types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listJob(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.Wctx.Batch.Job().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Add(types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listCronJob(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.Wctx.Batch.CronJob().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Add(types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listPod(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.Wctx.Core.Pod().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Add(types.NewWorkload(item))
	}
	return result, err
}

func (h *KubeHelper) listNamespace(
	_ string, opts metav1.ListOptions,
) (*listResult, error) {
	list, err := h.Wctx.Core.Namespace().List(opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, list := range list.Items {
		result.Add(types.NewResource(list))
	}
	return result, nil
}

func (h *KubeHelper) listNode(
	_ string, opts metav1.ListOptions,
) (*listResult, error) {
	list, err := h.Wctx.Core.Node().List(opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, list := range list.Items {
		result.Add(types.NewNode(list))
	}
	return result, nil
}

func (h *KubeHelper) listService(
	ns string, opts metav1.ListOptions,
) (*listResult, error) {
	list, err := h.Wctx.Core.Service().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, list := range list.Items {
		result.Add(types.NewService(list))
	}
	return result, nil
}

func (h *KubeHelper) listEvent(
	ns string, opts metav1.ListOptions,
) (*listResult, error) {
	list, err := h.Wctx.Core.Event().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, list := range list.Items {
		result.Add(types.NewEvent(list))
	}
	return result, nil
}

func (h *KubeHelper) ListResource(
	resource string,
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
	switch strings.TrimSuffix(strings.ToLower(resource), "s") {
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
	case "node":
		listFunc = h.listNode
	case "namespace":
		listFunc = h.listNamespace
	case "service":
		listFunc = h.listService
	case "event":
		listFunc = h.listEvent
	default:
		return "", fmt.Errorf("unsupported workload type: %s", resource)
	}

	ns = strings.ToLower(ns)
	switch ns {
	case "*":
		ns = ""
	}

	result, err := listFunc(ns, opts)
	if err != nil {
		return "", fmt.Errorf("failed to list %v: %w", resource, err)
	}
	return result.String(), nil
}

func (h *KubeHelper) listResourceHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	_ = ctx
	arguments := request.GetArguments()
	if arguments == nil {
		return nil, fmt.Errorf("failed to get request arguments")
	}
	name, ok := arguments["resource"].(string)
	if !ok {
		return nil, errors.New("resource not provided")
	}
	namespace, _ := arguments["namespace"].(string)
	labels, _ := arguments["labels"].([]any)
	limit, _ := arguments["limit"].(float64)

	s := make([]string, 0, len(labels))
	for _, label := range labels {
		if str, ok := label.(string); ok {
			s = append(s, str)
		}
	}
	result, err := h.ListResource(name, namespace, s, int64(limit))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	logrus.Debugf("handle the ListResource Handler")
	return mcp.NewToolResultText(result), nil
}
