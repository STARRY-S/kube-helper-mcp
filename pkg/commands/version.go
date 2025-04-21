package commands

import (
	"fmt"

	"github.com/STARRY-S/kube-helper-mcp/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	*baseCmd
}

func newVersionCmd() *versionCmd {
	cc := &versionCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:     "version",
		Short:   "Show version",
		Example: "  helper version",
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.SetupLogrus(cc.hideLogTime)
			if cc.debug {
				logrus.SetLevel(logrus.DebugLevel)
				logrus.Debugf("Debug output enabled")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version %s\n", getVersion())
		},
	})

	return cc
}

func (cc *versionCmd) getCommand() *cobra.Command {
	return cc.cmd
}

func getVersion() string {
	if utils.Commit != "" {
		return fmt.Sprintf("%s - %s", utils.Version, utils.Commit)
	}
	return utils.Version
}
