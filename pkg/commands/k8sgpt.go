package commands

import (
	"fmt"

	"github.com/STARRY-S/kube-helper-mcp/pkg/k8sgpt"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rancher/wrangler/v3/pkg/kubeconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type k8sGPTCmd struct {
	*baseCmd

	sse    bool
	listen string
	port   int
}

func newK8sGPTCmd() *k8sGPTCmd {
	cc := &k8sGPTCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "k8sgpt",
		Short: "MCP server to make K8sGPT operator actions",
		Long:  "",
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.SetupLogrus(cc.hideLogTime)
			if cc.debug {
				logrus.SetLevel(logrus.DebugLevel)
				logrus.Debugf("Debug output enabled")
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.run()
		},
	})

	flags := cc.baseCmd.cmd.Flags()
	flags.BoolVarP(&cc.sse, "sse", "", false, "Use SSE protocol instead of stdio")
	flags.StringVarP(&cc.listen, "listen", "l", "127.0.0.1", "SSE Listen Address")
	flags.IntVarP(&cc.port, "port", "p", 8000, "SSE Listen Port")
	addCommands(cc.cmd)
	return cc
}

func (cc *k8sGPTCmd) run() error {
	cfg, err := kubeconfig.GetNonInteractiveClientConfig(cc.baseCmd.kubeConfig).ClientConfig()
	if err != nil {
		return fmt.Errorf("building kubeconfig: %w", err)
	}
	h := k8sgpt.NewK8sGPTHelper(&k8sgpt.Options{
		Cfg: cfg,
	})
	s := h.Server()

	if cc.sse {
		addr := fmt.Sprintf("%v:%v", cc.listen, cc.port)
		u := fmt.Sprintf("http://%v", addr)
		sseServer := server.NewSSEServer(h.Server(), server.WithBaseURL(u))
		logrus.Printf("SSE server listening on %q", u)
		return sseServer.Start(addr)
	}
	return server.ServeStdio(s)
}
