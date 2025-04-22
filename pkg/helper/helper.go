package helper

import (
	"strings"

	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/STARRY-S/kube-helper-mcp/pkg/wrangler"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/rest"
)

// KubeHelper is a struct that provides methods to make Kubernetes API calls.
type KubeHelper struct {
	wctx *wrangler.Context
}

func NewKubeHelper(c *rest.Config) *KubeHelper {
	wctx := wrangler.NewContextOrDie(c)
	return &KubeHelper{
		wctx: wctx,
	}
}

func (h *KubeHelper) Server() *server.MCPServer {
	// Create MCP server
	s := server.NewMCPServer(
		"kube_api_call",
		strings.TrimPrefix(utils.Version, "v"), // version does not has 'v' prefix
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add list_workload tool
	s.AddTool(mcp.NewTool(
		"list_workloads",
		mcp.WithDescription(`List the real-time kubernetes cluster workloads with status information in JSON format,
different workloads and namespaces will produce different results.`),
		mcp.WithString(
			"workload",
			mcp.Required(),
			mcp.Description("The kubernetes workload kind to query (pod, deployment, statefulset, daemonset)"),
			mcp.Enum("pod", "deployment", "statefulset", "daemonset", "job", "cronjob"),
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
	), h.listWorkloadHandler)

	// Add list_namespace tool
	s.AddTool(mcp.NewTool(
		"list_namespaces",
		mcp.WithDescription(`List the real-time kubernetes cluster namespaces JSON format`),
		mcp.WithNumber(
			"limit",
			mcp.Description("The limit of the namespaces to query"),
			mcp.DefaultNumber(50),
		),
	), h.listNamespaceHandler)

	// Add get_workload tool
	s.AddTool(mcp.NewTool(
		"get_workload",
		mcp.WithDescription(`Get the real-time kubernetes workload detailed information in JSON format`),
		mcp.WithString(
			"workload",
			mcp.Required(),
			mcp.Description("The kubernetes workload kind to query (pod, deployment, statefulset, daemonset)"),
			mcp.Enum("pod", "deployment", "statefulset", "daemonset", "job", "cronjob"),
			mcp.DefaultString("pod"),
		),
		mcp.WithString(
			"namespace",
			mcp.Description("The kubernetes namespace to query"),
			mcp.DefaultString(""),
		),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("The kubernetes workload resource name to get, must be provided"),
		),
	), h.getWorkloadHandler)

	return s
}
