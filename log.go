package accumulator

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogToStdOut creates a new file logger
func NewLogToStdOut(tag, version string, prod bool) *zap.SugaredLogger {

	if prod {
		config := zap.NewProductionConfig()
		l, err := config.Build()
		if err != nil {
			panic("can't initialize zap logger: " + err.Error())
		}
		return l.Sugar().With("version", version).With("tag", tag)
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	l, err := config.Build()
	if err != nil {
		panic("can't initialize zap logger: " + err.Error())
	}
	return l.Sugar().With("version", version).With("tag", tag)
}
