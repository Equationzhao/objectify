package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetupZap sets up a new zap logger, and replace global zap-logger(zap.L()) with it
// actual OutputPaths and ErrorOutputPaths are stdout + OutputPaths and stderr + ErrorOutputPaths
func SetupZap(OutputPaths, ErrorOutputPaths []string, Level zapcore.Level) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(Level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      append([]string{"stdout"}, OutputPaths...),
		ErrorOutputPaths: append([]string{"stderr"}, ErrorOutputPaths...),
	}
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 - 15:04:05")
	config.EncoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[caller=" + caller.TrimmedPath() + "]")
	}
	config.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
}
