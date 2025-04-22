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

	addCommands(cc.cmd)
	return cc
}

func (cc *runCmd) run() error {
	cfg, err := kubeconfig.GetNonInteractiveClientConfig(cc.baseCmd.kubeConfig).ClientConfig()
	if err != nil {
		return fmt.Errorf("building kubeconfig: %w", err)
	}
	h := helper.NewKubeHelper(cfg)
	s := h.Server()
	return server.ServeStdio(s)
}
