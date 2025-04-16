package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/STARRY-S/learn-mcp/pkg/utils"
	"github.com/STARRY-S/learn-mcp/pkg/utils/kubehelper"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rancher/wrangler/v3/pkg/kubeconfig"
	"github.com/sirupsen/logrus"
)

var (
	version        bool
	sse            bool
	bind           string
	port           int
	kubeConfigFile string

	h *kubehelper.Helper
)

func init() {
	flag.BoolVar(&version, "version", false, "Show version")
	flag.BoolVar(&sse, "sse", false, "Use SSE for streaming output")
	flag.StringVar(&bind, "bind", "0.0.0.0", "Bind address")
	flag.IntVar(&port, "port", 8188, "Bind port")
	flag.StringVar(&kubeConfigFile, "kubeconfig", "", "Kube-config file (optional)")
	flag.Parse()
}

func main() {
	if version {
		fmt.Printf("%v - %v\n", utils.Version, utils.Commit)
		return
	}

	cfg, err := kubeconfig.GetNonInteractiveClientConfig(kubeConfigFile).ClientConfig()
	if err != nil {
		logrus.Fatalf("Error building kubeconfig: %v", err)
	}
	h = kubehelper.NewHelper(cfg)

	// Create MCP server
	s := server.NewMCPServer(
		"kube_helper",
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
	s.AddTool(kubeCheckTool, kubeCheckHandler)

	if sse {
		listen := fmt.Sprintf("%v:%v", bind, port)
		u := fmt.Sprintf("http://%v", listen)
		sseServer := server.NewSSEServer(s,
			server.WithBaseURL(u),
		)
		logrus.Infof("Listen on %q", u)
		err = sseServer.Start(listen)
	} else {
		err = server.ServeStdio(s)
	}
	if err != nil {
		logrus.Fatalf("%v", err)
	}
}

func kubeCheckHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	if h == nil {
		return nil, errors.New("kubehelper is not initialized")
	}

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
