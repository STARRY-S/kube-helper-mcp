package commands

import (
	"github.com/spf13/cobra"
)

func Execute(args []string) error {
	cmd := newHelperCmd()
	cmd.addCommands()
	cmd.cmd.SetArgs(args)

	_, err := cmd.cmd.ExecuteC()
	if err != nil {
		if signalContext.Err() != nil {
			return signalContext.Err()
		}
		return err
	}
	return nil
}

type helperCmd struct {
	*baseCmd
}

func newHelperCmd() *helperCmd {
	cc := &helperCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "helper",
		Short: "Kubernetes MCP helper",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	})
	cc.cmd.Version = getVersion()
	cc.cmd.SilenceUsage = true
	cc.cmd.SilenceErrors = true

	flags := cc.cmd.PersistentFlags()
	flags.StringVarP(&cc.baseCmd.kubeConfig, "kubeconfig", "c", "", "kube-config file (optional)")
	flags.BoolVarP(&cc.baseCmd.debug, "debug", "", false, "enable debug output")
	flags.BoolVar(&cc.baseCmd.hideLogTime, "hide-log-time", false, "hide log output timestamp")
	flags.MarkHidden("hide-log-time")

	return cc
}

func (cc *helperCmd) getCommand() *cobra.Command {
	return cc.cmd
}

func (cc *helperCmd) addCommands() {
	addCommands(
		cc.cmd,
		newVersionCmd(),
		newListCmd(),
	)
}
