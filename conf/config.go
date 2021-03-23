package conf

import (
	"github.com/jjauzion/ws-backend/internal/logger"
	"github.com/spf13/viper"
)

type Configuration struct {
	// ElasticSearch
	ES_USER_MAPPING string
	ES_TASK_MAPPING string
	BOOTSTRAP_FILE  string

	// WorkStation
	WS_API_PORT  string
	WS_GRPC_HOST string
	WS_GRPC_PORT string

	// WorkStation_ElasticSearch
	WS_ES_HOST string
	WS_ES_PORT string

	// JWT
	JWT_SIGNIN_KEY string
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
