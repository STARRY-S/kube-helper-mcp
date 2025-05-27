package utils

import (
	"encoding/json"
	"io"
	"os"

	"github.com/STARRY-S/simple-logrus-formatter/pkg/formatter"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var (
	Version     = "0.1.0"
	Commit      = "UNKNOWN"
	LogFileName = ".helper.log"
)

type MCPProtocol string

const (
	ProtocolStdio MCPProtocol = "stdio"
	ProtocolSSE   MCPProtocol = "sse"
	ProtocolHTTP  MCPProtocol = "http"
)

func SetupLogrus(hideTime bool) {
	formatter := &formatter.Formatter{
		NoColors: false,
	}
	if hideTime {
		formatter.TimestampFormat = "-"
	}
	logrus.SetFormatter(formatter)
	logrus.SetOutput(io.Discard)
	logrus.AddHook(&writer.Hook{
		// Send logs with level higher than warning to stderr.
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
	logrus.AddHook(&writer.Hook{
		// Send info, debug and trace logs to stdout.
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.TraceLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})
}

func Print(a any) string {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		logrus.Warnf("utils.Print: failed to json marshal (%T): %v", a, err)
	}
	return string(b)
}
