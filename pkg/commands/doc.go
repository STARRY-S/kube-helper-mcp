package commands

import (
	"github.com/spf13/cobra/doc"
)

func Doc(dir string) error {
	cmd := newHelperCmd()
	cmd.addCommands()

	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "1",
	}
	return doc.GenManTree(cmd.cmd, header, dir)
}
