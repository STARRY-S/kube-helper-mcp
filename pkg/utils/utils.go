package utils

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

var (
	Version = "0.1.0"
	Commit  = "UNKNOWN"
)

func Print(a any) string {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		logrus.Warnf("utils.Print: failed to json marshal (%T): %v", a, err)
	}
	return string(b)
}
