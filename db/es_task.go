package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func (es *esHandler) CreateTask(ctx context.Context, task Task) error {
	es.log.Debug("creating new task...")
	_, err := es.client.Index().Index(taskIndex).Id(task.ID).BodyJson(task).Do(ctx)
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

func (es *esHandler) GetOldestTask(ctx context.Context) (*Task, error) {
	es.log.Debug("searching most recent task")
	q := elastic.NewMatchAllQuery()
	s := elastic.NewSearchSource()
	s = s.Query(q)
	s = s.Sort(taskFieldCreatedAt, true)
	s = s.Size(2)
	tasks, err := es.searchTasks(ctx, s)
	if err != nil {
		return nil, err
	}
	if len(tasks) <= 0 {
		return nil, nil
	}
	es.log.Info("search successfully completed")
	return &tasks[0], err
}

func (es *esHandler) UpdateTaskTime(ctx context.Context, taskID string, status Status) error {
	es.log.Info("updating status", zap.String("task_id", taskID), zap.String("status", status.String()))
	if !status.IsValid() {
		panic("'" + status + "' is not a valid status.")
	}
	update, err := es.client.Update().Index(taskIndex).Id(taskID).Doc(map[string]string{taskFieldStatus: status.String()}).Do(ctx)
	if err != nil {
		return err
	}
	es.log.Info("status updated", zap.String("task_id", update.Id))
	return nil
}

func (es *esHandler) UpdateTaskStatus(ctx context.Context, taskID string, status Status) error {
	es.log.Info("updating status", zap.String("task_id", taskID), zap.String("status", status.String()))
	if !status.IsValid() {
		panic("'" + status + "' is not a valid status.")
	}
	var doc = map[string]interface{}{}
	doc[taskFieldStatus] = status.String()
	switch status {
	case StatusRunning:
		{
			doc[taskFieldStartedAt] = time.Now()
		}
	case StatusEnded:
		{
			doc[taskFieldEndedAt] = time.Now()
		}
	case StatusNotStarted:
		{
			doc[taskFieldStartedAt] = time.Unix(0, 0)
			doc[taskFieldEndedAt] =  time.Unix(0, 0)
		}
	}
	update, err := es.client.Update().Index(taskIndex).Id(taskID).Doc(doc).Do(ctx)
	if err != nil {
		return err
	}
	es.log.Info("status updated", zap.String("task_id", update.Id))
	return nil
}

func (es *esHandler) DeleteTaskByID(ctx context.Context, id string) error {
	es.log.Debug("delete task", zap.String("id", id))
	q := elastic.NewMatchQuery("id", id)
	_, err := es.client.DeleteByQuery().Index(taskIndex).Query(q).Do(ctx)

	if err != nil {
		es.log.Error("cannot delete task", zap.String("id", id), zap.Error(err))
	}

	return nil
}

func (es *esHandler) DeleteTasksBysUserID(ctx context.Context, userID string) error {
	q := elastic.NewMatchQuery("user_id", userID)
	_, err := es.client.DeleteByQuery().Index(taskIndex).Query(q).Do(ctx)
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
