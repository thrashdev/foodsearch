package logging

import "go.uber.org/zap"

func NewLogger(logFilePath string) (*zap.SugaredLogger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"stdout",
		logFilePath,
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, nil
	}
	return logger.Sugar(), nil
}
