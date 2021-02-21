package db

import (
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/jjauzion/ws-backend/internal"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"strings"
)

func (es *esHandler) GetUserByEmail(email string) (user model.User, err error) {
	logger := internal.GetLogger()
	logger.Debug("searching user by email...")
	search := `{
		"query": {
			"match": {
			  "email.keyword": {
				"query": "test@gmail.com"
			  }
			}
		}
    }`
	res, err := es.search([]string{userIndex}, strings.NewReader(search))
	if err != nil {
		logger.Error("", zap.Error(err))
	}
	if len(res) == 0 {
		logger.Info("user not found", zap.String("email", email))
		return
	}
	if len(res) > 1 {
		logger.Info("email matches more than one user", zap.String("email", email))
		return
	}
	if err = mapstructure.Decode(res[0], &user); err != nil {
		logger.Error("can't decode user")
		return
	}
	logger.Debug("search successfully completed")
	return
}
