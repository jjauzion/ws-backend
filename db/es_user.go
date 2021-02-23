package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"strings"
)

var (
	ErrNotFound    = fmt.Errorf("user not found")
	ErrTooManyRows = fmt.Errorf("found too many rows")
)

func (es *esHandler) searchUser(query, param string) (model.User, error) {
	user := model.User{}
	res, err := es.search([]string{userIndex}, strings.NewReader(query))
	if err != nil {
		es.log.Error("", zap.Error(err))
	}
	if len(res) == 0 {
		es.log.Info("user not found", zap.String("not found", param))
		return user, ErrNotFound
	}
	if len(res) > 1 {
		es.log.Info("more than one user found", zap.String("not unique", param))
		return user, ErrTooManyRows
	}
	if err = mapstructure.Decode(res[0], user); err != nil {
		es.log.Error("can't decode user", zap.Error(err))
		return user, err
	}
	return user, nil
}

func (es *esHandler) GetUserByID(id string) (user model.User, err error) {
	es.log.Debug("searching user by id...")
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
	es.log.Info("search successfully completed")
	return
}

func (es *esHandler) GetUserByEmail(email string) (user model.User, err error) {
	es.log.Debug("searching user by email...")
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
	es.log.Info("search successfully completed")
	return
}

func (es *esHandler) CreateUser(user model.User) (err error) {
	es.log.Debug("creating new user...")
	_, err = es.GetUserByEmail(user.Email)
	if err != nil && err != ErrNotFound {
		es.log.Error("cannot check if user exist", zap.Error(err))
		return fmt.Errorf("cannot create user")
	} else if err == nil {
		es.log.Info("can't create user, user already exists", zap.String("email", user.Email))
		err = errors.New("user already exists")
		return
	}
	var b []byte
	if b, err = json.Marshal(user); err != nil {
		es.log.Error("failed to create user", zap.Error(err))
		return
	}
	if err = es.indexNewDoc(userIndex, bytes.NewReader(b)); err != nil {
		es.log.Error("failed to create user", zap.Error(err))
		return
	}
	es.log.Info("new user successfully created", zap.String("email", user.Email))
	return
}
