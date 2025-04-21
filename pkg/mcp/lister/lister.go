package lister

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/STARRY-S/kube-helper-mcp/pkg/internal/types"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/STARRY-S/kube-helper-mcp/pkg/wrangler"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

// Lister is a struct that provides methods to list Kubernetes resources.
type Lister struct {
	wctx *wrangler.Context
}

func NewLister(c *rest.Config) *Lister {
	wctx := wrangler.NewContextOrDie(c)
	return &Lister{
		wctx: wctx,
	}
}

func (h *Lister) listDeployment(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Apps.Deployment().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Results = append(result.Results, types.NewWorkload(item))
	}
	return result, err
}

func (h *Lister) listDaemonSet(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Apps.DaemonSet().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Results = append(result.Results, types.NewWorkload(item))
	}
	return result, err
}

func (h *Lister) listStatefulSet(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Apps.StatefulSet().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Results = append(result.Results, types.NewWorkload(item))
	}
	return result, err
}

func (h *Lister) listPod(ns string, opts metav1.ListOptions) (*listResult, error) {
	list, err := h.wctx.Core.Pod().List(ns, opts)
	if err != nil {
		return nil, err
	}
	result := &listResult{}
	for _, item := range list.Items {
		result.Results = append(result.Results, types.NewWorkload(item))
	}
	return result, err
}

func (h *Lister) ListWorkload(
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
	return a.String(), nil
}

func (h *Lister) Server() *server.MCPServer {
	// Create MCP server
	s := server.NewMCPServer(
		"lister",
		utils.Version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add a calculator tool
	kubeCheckTool := mcp.NewTool(
		"list_workload",
		mcp.WithDescription("List kubernetes cluster workloads with detailed informations"),
		mcp.WithString(
			"workload",
			mcp.Required(),
			mcp.Description("The kubernetes workload kind to query (pod, deployment, statefulset, daemonset)"),
			mcp.Enum("pod", "deployment", "statefulset", "daemonset"),
		),
		mcp.WithString(
			"namespace",
			mcp.Description("The kubernetes namespace to query"),
			mcp.DefaultString("default"),
		),
		mcp.WithArray(
			"labels",
			mcp.Description("The label of the workload to query"),
			mcp.DefaultString(""),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("The limit of the workload to query"),
			mcp.DefaultNumber(50),
		),
	)

	// Add tool handler
	s.AddTool(kubeCheckTool, h.kubeCheckHandler)
	return s
}

func (h *Lister) kubeCheckHandler(
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
