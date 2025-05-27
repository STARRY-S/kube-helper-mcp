package commands

import (
	"github.com/STARRY-S/kube-helper-mcp/pkg/signal"
	"github.com/spf13/cobra"
)

var (
	signalContext = signal.SetupSignalContext()
)

type baseCmd struct {
	*baseOpts
	cmd *cobra.Command
}

func newBaseCmd(cmd *cobra.Command) *baseCmd {
	return &baseCmd{cmd: cmd, baseOpts: &globalOpts}
}

type baseOpts struct {
	kubeConfig  string
	debug       bool // Enable debug output
	hideLogTime bool // Hide log output time (used in validation test)
}

var globalOpts = baseOpts{}

func (cc *baseCmd) getCommand() *cobra.Command {
	return cc.cmd
}

type cmder interface {
	getCommand() *cobra.Command
}

func addCommands(root *cobra.Command, commands ...cmder) {
	for _, command := range commands {
		cmd := command.getCommand()
		if cmd == nil {
			continue
		}
		root.AddCommand(cmd)
	}
}
