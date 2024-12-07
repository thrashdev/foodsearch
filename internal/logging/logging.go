package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logFilePath string) (*zap.SugaredLogger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"stdout",
		logFilePath,
	}
	cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	cfg.EncoderConfig.CallerKey = "caller"

	logger, err := cfg.Build()
	if err != nil {
		return nil, nil
	}
	return logger.Sugar(), nil
}
