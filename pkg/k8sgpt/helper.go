package k8sgpt

import (
	"strings"

	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/STARRY-S/kube-helper-mcp/pkg/wrangler"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"k8s.io/client-go/rest"
)

type Helper struct {
	wctx *wrangler.Context
}

type Options struct {
	Cfg *rest.Config
}

func NewK8sGPTHelper(o *Options) *Helper {
	wctx := wrangler.NewContextOrDie(o.Cfg)
	return &Helper{
		wctx: wctx,
	}
}

func (h *Helper) Server() *server.MCPServer {
	// Create MCP server
	s := server.NewMCPServer(
		"k8sgpt_helper",
		strings.TrimPrefix(utils.Version, "v"), // version does not has 'v' prefix
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add list_workload tool
	s.AddTool(mcp.NewTool(
		"check_cluster",
		mcp.WithDescription(
			`Trigger the K8sGPT cluster self-check (inspection) actions, do not return the status result.`),
	), h.checkClusterHandler)

	s.AddTool(mcp.NewTool(
		"get_check_results",
		mcp.WithDescription(`Get all K8sGPT cluster self-check (inspection) results in JSON format.`),
	), h.getCheckResultsHandler)

	s.AddTool(mcp.NewTool(
		"remediate_cluster",
		mcp.WithDescription(`Use K8sGPT to remediate the cluster, only triggers the actions, do not return the remediate mutation result.`),
	), h.remediateClusterHandler)

	s.AddTool(mcp.NewTool(
		"get_mutation_result",
		mcp.WithDescription(`Get the K8sGPT remediate mutation result in JSON format.`),
	), h.getRemediateResultHandler)

	return s
}
