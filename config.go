package zapflatencoder

import (
	"strconv"

	"go.uber.org/zap/zapcore"
)

const (
	// EncoderName registered with zap for usage
	EncoderName            = "flat-encoder"
	tokenTab               = '\t'
	tokenReplacement       = `\ufffd`
	tokenLineEnding        = zapcore.DefaultLineEnding
	tokenNamespaceOpen     = '{'
	tokenNamespaceClose    = '}'
	tokenArrayOpen         = '['
	tokenArrayClose        = ']'
	tokenKeyValueSeparator = '='
	tokenStringEnclosed    = '"'
)

// ShortColorCallerEncoder encodes caller information with sort path filename and enable color.
func ShortColorCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	callerStr := caller.TrimmedPath() + ":" + strconv.Itoa(caller.Line)
	enc.AppendString(callerStr)
}

// DefaultConfig ...
var DefaultConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   ShortColorCallerEncoder,
}
