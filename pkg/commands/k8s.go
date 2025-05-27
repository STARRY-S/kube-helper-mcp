package commands

import (
	"fmt"

	"github.com/STARRY-S/kube-helper-mcp/pkg/helper"
	"github.com/STARRY-S/kube-helper-mcp/pkg/signal"
	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/rancher/wrangler/v3/pkg/kubeconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type k8sCmd struct {
	*baseCmd

	protocol string
	listen   string
	port     int
}

func newK8sCmd() *k8sCmd {
	cc := &k8sCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "k8s",
		Short: "MCP server to make K8s API get/list operations",
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
	flags.StringVarP(&cc.protocol, "protocol", "p", "stdio", "MCP Protocol (stdio,http,sse)")
	flags.StringVarP(&cc.listen, "listen", "l", "", "Listen Address")
	flags.IntVarP(&cc.port, "port", "", 8000, "Listen Port")
	addCommands(cc.cmd)
	return cc
}

func (cc *k8sCmd) run() error {
	cfg, err := kubeconfig.GetNonInteractiveClientConfig(cc.baseCmd.kubeConfig).ClientConfig()
	if err != nil {
		return fmt.Errorf("building kubeconfig: %w", err)
	}
	h := helper.NewK8sGPTHelper(&helper.Options{
		Cfg:      cfg,
		Protocol: utils.MCPProtocol(cc.protocol),
		Listen:   cc.listen,
		Port:     cc.port,
	})
	signal.RegisterOnShutdown(func() error {
		return h.Shutdown(signalContext)
	})
	return h.Serve(signalContext)
}
