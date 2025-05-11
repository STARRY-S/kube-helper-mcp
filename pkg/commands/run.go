package commands

import (
	"fmt"

	"github.com/STARRY-S/kube-helper-mcp/pkg/helper"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rancher/wrangler/v3/pkg/kubeconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type runCmd struct {
	*baseCmd

	sse    bool
	listen string
	port   int
}

func newRunCmd() *runCmd {
	cc := &runCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "run",
		Short: "MCP server to make kubernetes resources API calls",
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

func (cc *runCmd) run() error {
	cfg, err := kubeconfig.GetNonInteractiveClientConfig(cc.baseCmd.kubeConfig).ClientConfig()
	if err != nil {
		return fmt.Errorf("building kubeconfig: %w", err)
	}
	h := helper.NewKubeHelper(&helper.Options{
		Cfg: cfg,
	})
	s := h.Server()

	if cc.sse {
		u := fmt.Sprintf("http://%v:%v", cc.listen, cc.port)
		sseServer := server.NewSSEServer(h.Server(), server.WithBaseURL(u))
		logrus.Printf("SSE server listening on %q", u)
		return sseServer.Start(fmt.Sprintf(":%d", cc.port))
	}
	return server.ServeStdio(s)
}
