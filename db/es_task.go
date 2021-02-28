package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"strings"
)

func (es *esHandler) CreateTask(task Task) (err error) {
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

func (es *esHandler) GetTasksByUserID(ctx context.Context, id string) ([]Task, error) {
	es.log.Debug("get tasks for user", zap.String("user_id", id))
	query := elastic.NewMatchQuery("user_id", id)
	tasks, err := es.searchTasks(ctx, query)
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

func (es *esHandler) searchTasks(ctx context.Context, query *elastic.MatchQuery) ([]Task, error) {
	var tasks []Task

	_, err := es.elasticSearch(ctx, taskIndex, query, func(hit *elastic.SearchHit) error {
		var task Task
		err := json.Unmarshal(hit.Source, &task)
		if err != nil {
			es.log.Warn("json failed", zap.Error(err))
			return err
		}

		es.log.Debug("task found", zap.String(task.UserId, "ok"))

		tasks = append(tasks, task)
		return nil
	})

	return tasks, err
}
