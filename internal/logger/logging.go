package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func ProvideLogger() (*Logger, error) {
	lg, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		return nil, err
	}

	return &Logger{lg}, nil
}
