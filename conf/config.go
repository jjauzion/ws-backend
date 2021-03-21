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
	WS_GRPC_HOST    string
	WS_GRPC_PORT    string
}

func GetConfig(log logger.Logger) (Configuration, error) {
	cf := Configuration{}
	err := viper.Unmarshal(&cf)
	if err != nil {
		return cf, err
	}
	log.Info("configuration loaded")
	return cf, nil
}
