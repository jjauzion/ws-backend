package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"strings"
)

func (es *esHandler) CreateTask(task model.Task) (err error) {
	es.log.Debug("creating new task...")
	var b []byte
	if b, err = json.Marshal(task); err != nil {
		return
	}
	if err = es.indexNewDoc(taskIndex, bytes.NewReader(b)); err != nil {
		es.log.Error("failed to create task", zap.Error(err))
		return
	}
	es.log.Info("new task successfully created", zap.String("id", task.ID))
	return
}

func (es *esHandler) GetTasksByUserID(ctx context.Context, id string) ([]model.Task, error) {
	es.log.Debug("get tasks for user", zap.String("user_id", id))
	search := fmt.Sprintf(`{
		"query": {
			"match": {
			  "created_by": {
				"query": "%s"
			  }
			}
		}
    }`, id)
	tasks, err := es.searchTasks(nil, search)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (es *esHandler) DeleteTask(ctx context.Context, id string) error {
	q := fmt.Sprintf(`{
		"query": {
			"match": {
			  "id": {
				"query": "%s"
			  }
			}
		}
    }`, id)

	res, err := es.client.DeleteByQuery([]string{taskIndex}, strings.NewReader(q))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("%v", res)
	}
	return nil
}

func (es *esHandler) DeleteUserTasks(ctx context.Context, userID string) error {
	q := fmt.Sprintf(`{
		"query": {
			"match": {
			  "created_by": {
				"query": "%s"
			  }
			}
		}
    }`, userID)

	_, err := es.client.DeleteByQuery([]string{taskIndex}, strings.NewReader(q))
	if err != nil {
		return err
	}
	return nil
}

func (es *esHandler) searchTasks(ctx context.Context, query string) ([]model.Task, error) {
	tasks := []model.Task{}
	res, err := es.search([]string{taskIndex}, strings.NewReader(query))
	if err != nil {
		es.log.Error("es search", zap.Error(err))
		return nil, err
	}

	for _, re := range res {
		t := model.Task{}
		es.log.Debug("decode res", zap.Any("json", re))
		e := mapstructure.Decode(re, &t)
		if e != nil {
			err = multierr.Append(err, e)
		} else {
			tasks = append(tasks, t)
		}
		es.log.Debug("after decode", zap.Any("task", t))
	}

	return tasks, err
}
