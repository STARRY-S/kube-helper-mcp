package kubehelper

import (
	"fmt"
	"strings"

	"github.com/STARRY-S/learn-mcp/pkg/utils"
	"github.com/STARRY-S/learn-mcp/pkg/wrangler"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Helper struct {
	wctx *wrangler.Context
}

type listResult struct {
	Name   string
	Status string
}

func NewHelper(c *rest.Config) *Helper {
	wctx := wrangler.NewContextOrDie(c)
	return &Helper{
		wctx: wctx,
	}
}

func (h *Helper) listDeployment(ns string, opts v1.ListOptions) (any, error) {
	list, err := h.wctx.Apps.Deployment().List(ns, opts)
	if err != nil {
		return nil, err
	}
	var result []listResult
	for _, item := range list.Items {
		r := listResult{
			Name:   item.Name,
			Status: item.Status.String(),
		}
		result = append(result, r)
	}
	return result, err
}

func (h *Helper) listDaemonSet(ns string, opts v1.ListOptions) (any, error) {
	list, err := h.wctx.Apps.DaemonSet().List(ns, opts)
	if err != nil {
		return nil, err
	}
	var result []listResult
	for _, item := range list.Items {
		r := listResult{
			Name:   item.Name,
			Status: item.Status.String(),
		}
		result = append(result, r)
	}
	return result, err
}

func (h *Helper) listStatefulSet(ns string, opts v1.ListOptions) (any, error) {
	list, err := h.wctx.Apps.StatefulSet().List(ns, opts)
	if err != nil {
		return nil, err
	}
	var result []listResult
	for _, item := range list.Items {
		r := listResult{
			Name:   item.Name,
			Status: item.Status.String(),
		}
		result = append(result, r)
	}
	return result, err
}

func (h *Helper) listPod(ns string, opts v1.ListOptions) (any, error) {
	list, err := h.wctx.Core.Pod().List(ns, opts)
	if err != nil {
		return nil, err
	}
	var result []listResult
	for _, item := range list.Items {
		r := listResult{
			Name:   item.Name,
			Status: item.Status.String(),
		}
		result = append(result, r)
	}
	return result, err
}

func (h *Helper) ListWorkload(
	workload string,
	ns string,
	labels []string,
	limit int64,
) (string, error) {
	opts := v1.ListOptions{
		Limit: limit,
	}
	if len(labels) > 0 {
		opts.LabelSelector = strings.Join(labels, ",")
	}

	var listFunc func(string, v1.ListOptions) (interface{}, error)
	switch workload {
	case "deployment":
		listFunc = h.listDeployment
	case "daemonset":
		listFunc = h.listDaemonSet
	case "statefulset":
		listFunc = h.listStatefulSet
	case "pod":
		listFunc = h.listPod
	default:
		return "", fmt.Errorf("unsupported workload type: %s", workload)
	}

	a, err := listFunc(ns, opts)
	if err != nil {
		return "", fmt.Errorf("failed to list %v: %w", workload, err)
	}
	return utils.Print(a), nil
}
