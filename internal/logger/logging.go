package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func ProvideLogger(dev bool) (Logger, error) {
	logger := Logger{}
	var err error

	if dev {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Mon 02 Jan 15:04:05 MST")
		logger.Logger, err = config.Build(zap.AddCaller())
		if err != nil {
			return Logger{}, err
		}
		logger.Info("logger initialized in development mode")
	} else {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.EpochTimeEncoder
		logger.Logger, err = config.Build()
		if err != nil {
			return Logger{}, err
		}
	}

	return logger, nil
}

func MockLogger() Logger {
	return Logger{zap.NewNop()}
}
