package k8sgpt

import (
	"context"
	"fmt"
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

func (h *Helper) Server(ctx context.Context) (*server.MCPServer, error) {
	// Create MCP server
	s := server.NewMCPServer(
		"k8sgpt_helper",
		strings.TrimPrefix(utils.Version, "v"), // version does not has 'v' prefix
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
		server.WithInstructions("Trigger/Get the cluster self-check & remediate results, need to convert the result in table formats."),
	)

	// s.AddPrompt() // TODO: Add Prompts
	// s.AddResources() // TODO: Add Resources

	// Add session
	session := common.NewSession("k8sgpt-session-1")
	if err := s.RegisterSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to register session for server: %w", err)
	}

	// Add list_workload tool
	s.AddTool(mcp.NewTool(
		"check_cluster",
		mcp.WithDescription(
			`Trigger the K8sGPT cluster self-check actions, and return self-check results, need to convert the result in table format for better reading.`),
	), h.checkClusterHandler)

	s.AddTool(mcp.NewTool(
		"get_check_results",
		mcp.WithDescription(`Get all K8sGPT cluster self-check (inspection) results in JSON format, need to convert the result in table format for better reading.`),
	), h.getCheckResultsHandler)

	s.AddTool(mcp.NewTool(
		"remediate_cluster",
		mcp.WithDescription(`Use K8sGPT to remediate the cluster, and return mutation results, need to convert the result in table format for better reading.`),
	), h.remediateClusterHandler)

	s.AddTool(mcp.NewTool(
		"get_mutation_result",
		mcp.WithDescription(`Get the K8sGPT remediate mutation result in JSON format, need to convert the result in table format for better reading..`),
	), h.getRemediateResultHandler)

	return s, nil
}

func (h *Helper) Serve(ctx context.Context) error {
	server, err := h.Server(ctx)
	if err != nil {
		return err
	}
	return h.Common.Start(ctx, server)
}

func (h *Helper) Shutdown(ctx context.Context) error {
	return h.Common.Stop(ctx)
}
