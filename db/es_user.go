package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func (es *esHandler) CreateUser(ctx context.Context, user User) (err error) {
	es.log.Debug("creating new user...")
	_, err = es.GetUserByEmail(ctx, user.Email)

	if err != nil && err != ErrNotFound {
		es.log.Error("cannot check if user exist", zap.Error(err))
		return fmt.Errorf("cannot create user: %w", err)
	} else if err == nil {
		es.log.Info("can't create user, user already exists", zap.String("email", user.Email))
		return errors.New("user already exists")
	}
	_, err = es.elastic.Index().Index(userIndex).BodyJson(user).Do(ctx)
	if err != nil {
		es.log.Error("cannot create user, db failure", zap.Error(err))
		return err
	}

	es.log.Info("new user successfully created", zap.String("email", user.Email))
	return nil
}

func (es *esHandler) GetUserByID(ctx context.Context, id string) (User, error) {
	es.log.Debug("searching user by id...")
	query := elastic.NewMatchQuery("id", id)
	var user User

	hit, err := es.elasticSearchOne(ctx, userIndex, query)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(hit.Source, &user)
	if err != nil {
		es.log.Error("json failed", zap.Error(err))
		return user, err
	}

	es.log.Info("search successfully completed")
	return user, err
}

func (es *esHandler) GetUserByEmail(ctx context.Context, email string) (User, error) {
	es.log.Debug("searching user by email...")
	query := elastic.NewMatchQuery("email", email)
	var user User

	hit, err := es.elasticSearchOne(ctx, userIndex, query)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(hit.Source, &user)
	if err != nil {
		es.log.Error("json failed", zap.Error(err))
		return user, err
	}

	es.log.Info("search successfully completed")
	return user, err
}

func (es *esHandler) DeleteUser(ctx context.Context, id string) error {
	es.log.Debug("delete user", zap.String("id", id))

	q := elastic.NewMatchQuery("id", id)
	_, err := es.elastic.DeleteByQuery().Index(userIndex).Query(q).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
