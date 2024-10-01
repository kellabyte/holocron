package logging

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/kellabyte/holocron/node"
	"github.com/rs/zerolog"
)

// var customStandardOutput = zerolog.ConsoleWriter{
// 	Out:                 os.Stdout,
// 	NoColor:             false,
// 	TimeFormat:          time.Stamp,
// 	TimeLocation:        nil,
// 	PartsOrder:          []string{"time", "level", "node", "function", "role", "epoch", "message"},
// 	PartsExclude:        nil,
// 	FieldsOrder:         nil,
// 	FieldsExclude:       []string{"node", "function", "role", "epoch", "message"},
// 	FormatTimestamp:     nil,
// 	FormatLevel:         func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%-6s|", i)) }, // INFO  |
// 	FormatCaller:        func(i interface{}) string { return filepath.Base(fmt.Sprintf("%s", i)) },
// 	FormatMessage:       nil,
// 	FormatFieldName:     nil,
// 	FormatFieldValue:    nil,
// 	FormatErrFieldName:  nil,
// 	FormatErrFieldValue: nil,
// 	FormatExtra:         nil,
// 	FormatPrepare:       nil,
// }

const (
	ColorBlack = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite

	ColorBold     = 1
	ColorDarkGray = 90

	UnknownLevel = "???"
)

func Colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func CreateLogger() zerolog.Logger {
	buildInfo, _ := debug.ReadBuildInfo()
	writer := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}

	logger := zerolog.New(writer).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	return logger
}

func CreateNodeLogger(logger zerolog.Logger, nodeId uuid.UUID) zerolog.Logger {
	nodeLogger := logger.With().
		Str("node", node.NodeIdShort(nodeId)).
		Logger()

	return nodeLogger
}
