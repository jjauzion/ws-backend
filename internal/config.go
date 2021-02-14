package internal

import (
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

type configuration struct {
	WS_ES_HOST string
	WS_ES_PORT string
}

var cf *configuration

func GetConfig() *configuration {
	log := GetLogger()
	if cf != nil {
		return cf
	}
	cf = &configuration{}
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("", zap.Error(err))
	}
	if err := viper.Unmarshal(cf); err != nil {
		log.Fatal("unable to unmarshall config into struc:", zap.Error(err))
	}
	log.Info("configuration loaded")
	return cf
}
