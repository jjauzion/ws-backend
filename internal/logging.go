package internal

import (
	"fmt"
	"go.uber.org/zap"
)

var logger *zap.Logger

func GetLogger() (log *zap.Logger) {
	if logger != nil {
		return logger
	}
	var err error
	if log, err = zap.NewDevelopment(zap.AddCaller()); err != nil {
		fmt.Println("couldn't create logger because:", err)
	}
	return log
}