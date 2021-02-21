package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/jjauzion/ws-backend/internal"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"strings"
)

func (es *esHandler) GetUserByEmail(email string) (user *model.User, err error) {
	logger := internal.GetLogger()
	logger.Debug("searching user by email...")
	search := fmt.Sprintf(`{
		"query": {
			"match": {
			  "email.keyword": {
				"query": "%s"
			  }
			}
		}
    }`, email)
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
	user = &model.User{}
	if err = mapstructure.Decode(res[0], user); err != nil {
		logger.Error("can't decode user", zap.Error(err))
		return nil, err
	}
	logger.Info("search successfully completed")
	return
}

func (es *esHandler) CreateUser(user model.User) (err error) {
	logger := internal.GetLogger()
	logger.Debug("creating new user...")
	if tmp, err := es.GetUserByEmail(user.Email); err !=  nil {
		return err
	} else if tmp != nil {
		fmt.Println("here:", tmp)
		logger.Info("can't create user, user already exists", zap.String("email", user.Email))
		err = errors.New("user already exists")
		return err
	}
	var b []byte
	if b, err = json.Marshal(user); err != nil {
		return
	}
	if err = es.indexNewDoc(userIndex, bytes.NewReader(b)); err != nil {
		return
	}
	logger.Info("new user successfully created", zap.String("email", user.Email))
	return
}
