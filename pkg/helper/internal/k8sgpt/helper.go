package k8sgpt

import (
	"context"
	"strings"

	"github.com/STARRY-S/kube-helper-mcp/pkg/helper/internal/common"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Helper struct {
	*common.Common
}

type Options struct {
	*common.Options
}

func NewK8sGPTHelper(o *Options) *Helper {
	return &Helper{
		Common: common.NewCommon(o.Options),
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

func (h *Helper) Serve(ctx context.Context) error {
	return h.Common.Start(ctx, h.Server())
}

func (h *Helper) Shutdown(ctx context.Context) error {
	return h.Common.Stop(ctx)
}
