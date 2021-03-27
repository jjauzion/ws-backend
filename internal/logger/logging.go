package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func ProvideLogger() Logger {
	lg, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		panic(err)
	}

	return Logger{lg}
}
