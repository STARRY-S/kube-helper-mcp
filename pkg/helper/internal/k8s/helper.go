package k8s

import (
	"context"
	"strings"

	"github.com/STARRY-S/kube-helper-mcp/pkg/helper/internal/common"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// KubeHelper is a struct that provides methods to make Kubernetes API calls.
type KubeHelper struct {
	*common.Common
}

type Options struct {
	*common.Options
}

func NewKubeHelper(o *Options) *KubeHelper {
	return &KubeHelper{
		Common: common.NewCommon(o.Options),
	}
}

func (h *KubeHelper) Server() *server.MCPServer {
	// Create MCP server
	s := server.NewMCPServer(
		"kubernetes_helper",
		strings.TrimPrefix(utils.Version, "v"), // version does not has 'v' prefix
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add list_workload tool
	s.AddTool(mcp.NewTool(
		"list_resources",
		mcp.WithDescription(`List the kubernetes  with status information in JSON format`),
		mcp.WithString(
			"resource",
			mcp.Required(),
			mcp.Description("The kubernetes workload kind to query "+
				"(pod, deployment, statefulset, daemonset, job, cronjob, service, namespace, node)"),
			mcp.Enum("pod", "deployment", "statefulset", "daemonset", "job", "cronjob", "service", "namespace", "node", "event"),
			mcp.DefaultString("pod"),
		),
		mcp.WithString(
			"namespace",
			mcp.Description("The kubernetes namespace to query"),
			mcp.DefaultString(""),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("The limit of the workload to query"),
			mcp.DefaultNumber(50),
		),
	), h.listResourceHandler)

	// Add get_workload tool
	s.AddTool(mcp.NewTool(
		"get_single_resource",
		mcp.WithDescription(`Get one kubernetes resource detailed information in JSON format`),
		mcp.WithString(
			"resource",
			mcp.Required(),
			mcp.Description("The kubernetes resource kind to query "+
				"(pod, deployment, statefulset, daemonset, job, cronjob, service, namespace, node, event)"),
			mcp.Enum("pod", "deployment", "statefulset", "daemonset", "job", "cronjob", "service", "namespace", "node", "event"),
			mcp.DefaultString("pod"),
		),
		mcp.WithString(
			"namespace",
			mcp.Description("The kubernetes namespace of the resource to query"),
			mcp.DefaultString(""),
		),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The kubernetes resource name to get, must be provided"),
		),
	), h.getResourceHandler)

	return s
}

func (h *KubeHelper) Serve(ctx context.Context) error {
	return h.Common.Start(ctx, h.Server())
}

func (h *KubeHelper) Shutdown(ctx context.Context) error {
	return h.Common.Stop(ctx)
}
