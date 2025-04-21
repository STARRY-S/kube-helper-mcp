package main

import (
	"context"
	"errors"
	"os"

	"github.com/STARRY-S/kube-helper-mcp/pkg/commands"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := commands.Execute(os.Args[1:]); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		logrus.Warn(err)
	}
}
