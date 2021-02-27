package conf

import (
	"github.com/jjauzion/ws-backend/internal/logger"
	"github.com/spf13/viper"
)

type Configuration struct {
	WS_ES_HOST      string
	WS_ES_PORT      string
	ES_USER_MAPPING string
	ES_TASK_MAPPING string
	BOOTSTRAP_FILE  string
	WS_API_PORT     string
}

func GetConfig(log *logger.Logger) (Configuration, error) {
	cf := Configuration{}
	viper.AddConfigPath("conf")
	if err := viper.MergeInConfig(); err != nil {
		return cf, err
	}
	viper.SetConfigFile(".env")
	if err := viper.MergeInConfig(); err != nil {
		return cf, err
	}
	if err := viper.Unmarshal(&cf); err != nil {
		return cf, err
	}
	log.Info("Configuration loaded")
	return cf, nil
}
