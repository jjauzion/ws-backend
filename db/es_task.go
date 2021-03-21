package db

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func (es *esHandler) CreateTask(ctx context.Context, task Task) error {
	es.log.Debug("creating new task...")
	_, err := es.elastic.Index().Index(taskIndex).BodyJson(task).Do(ctx)
	if err != nil {
		es.log.Error("failed to create task", zap.Error(err))
		return err
	}

	es.log.Info("new task successfully created", zap.String("id", task.ID))
	return nil
}

func (es *esHandler) GetTasksByUserID(ctx context.Context, id string) ([]Task, error) {
	es.log.Debug("get tasks for user", zap.String("user_id", id))
	query := elastic.NewMatchQuery("user_id", id)
	s := elastic.NewSearchSource().Query(query)
	tasks, err := es.searchTasks(ctx, s)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (es *esHandler) DeleteTask(ctx context.Context, id string) error {
	es.log.Debug("delete task", zap.String("id", id))
	q := elastic.NewMatchQuery("id", id)
	_, err := es.elastic.DeleteByQuery().Index(taskIndex).Query(q).Do(ctx)

	if err != nil {
		es.log.Error("cannot delete task", zap.String("id", id), zap.Error(err))
	}

	return nil
}

func (es *esHandler) DeleteUserTasks(ctx context.Context, userID string) error {
	q := elastic.NewMatchQuery("user_id", userID)
	_, err := es.elastic.DeleteByQuery().Index(taskIndex).Query(q).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (es *esHandler) searchTasks(ctx context.Context, source *elastic.SearchSource) ([]Task, error) {
	var tasks []Task

	_, err := es.elasticSearch(ctx, taskIndex, source, func(hit *elastic.SearchHit) error {
		var task Task
		err := json.Unmarshal(hit.Source, &task)
		if err != nil {
			es.log.Error("json failed", zap.Error(err))
			return err
		}

		es.log.Debug("task found", zap.String(task.UserId, "ok"))

		tasks = append(tasks, task)
		return nil
	})

	return tasks, err
}
