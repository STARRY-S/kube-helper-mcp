package utils

import (
	"encoding/json"
	"io"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var (
	Version     = "0.1.0"
	Commit      = "UNKNOWN"
	LogFileName = ".helper.log"
)

func SetupLogrus(hideTime bool) {
	// f, err := os.Create(LogFileName)
	// if err != nil {
	// 	logrus.Fatalf("failed to create log file %s: %v", LogFileName, err)
	// }
	formatter := &nested.Formatter{
		HideKeys:        false,
		TimestampFormat: "[15:04:05]", // hour, time, sec only
		FieldsOrder:     []string{"IMG"},
		NoColors:        true,
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
