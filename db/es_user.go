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

func (es *esHandler) searchUser(query, param string) (user *model.User, err error) {
	logger := internal.GetLogger()
	res, err := es.search([]string{userIndex}, strings.NewReader(query))
	if err != nil {
		logger.Error("", zap.Error(err))
	}
	if len(res) == 0 {
		logger.Info("user not found", zap.String("not found", param))
		return
	}
	if len(res) > 1 {
		logger.Info("more than one user found", zap.String("not unique", param))
		return
	}
	user = &model.User{}
	if err = mapstructure.Decode(res[0], user); err != nil {
		logger.Error("can't decode user", zap.Error(err))
		return nil, err
	}
	return
}

func (es *esHandler) GetUserByID(id string) (user *model.User, err error) {
	logger := internal.GetLogger()
	logger.Debug("searching user by id...")
	search := fmt.Sprintf(`{
		"query": {
			"match": {
			  "id": {
				"query": "%s"
			  }
			}
		}
    }`, id)
	if user, err = es.searchUser(search, id); err != nil {
		return
	}
	logger.Info("search successfully completed")
	return
}

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
	if user, err = es.searchUser(search, email); err != nil {
		return
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
		logger.Info("can't create user, user already exists", zap.String("email", user.Email))
		err = errors.New("user already exists")
		return err
	}
	var b []byte
	if b, err = json.Marshal(user); err != nil {
		logger.Error("failed to create user", zap.Error(err))
		return
	}
	if err = es.indexNewDoc(userIndex, bytes.NewReader(b)); err != nil {
		logger.Error("failed to create user", zap.Error(err))
		return
	}
	logger.Info("new user successfully created", zap.String("email", user.Email))
	return
}
