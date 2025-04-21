package commands

import (
	"fmt"

	"github.com/STARRY-S/kube-helper-mcp/pkg/mcp/lister"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rancher/wrangler/v3/pkg/kubeconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type listCmd struct {
	*baseCmd
}

func newListCmd() *listCmd {
	cc := &listCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "list",
		Short: "MCP server to list kubernetes resources",
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

	// flags := cc.baseCmd.cmd.Flags()

	addCommands(cc.cmd)
	return cc
}

func (cc *listCmd) run() error {
	cfg, err := kubeconfig.GetNonInteractiveClientConfig(cc.baseCmd.kubeConfig).ClientConfig()
	if err != nil {
		return fmt.Errorf("building kubeconfig: %w", err)
	}
	l := lister.NewLister(cfg)
	s := l.Server()
	return server.ServeStdio(s)
}
