package db

import (
	"bytes"
	"encoding/json"
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/jjauzion/ws-backend/internal"
	"go.uber.org/zap"
)

func (es *esHandler) CreateTask(task model.Task) (err error) {
	logger := internal.GetLogger()
	logger.Debug("creating new task...")
	var b []byte
	if b, err = json.Marshal(task); err != nil {
		return
	}
	if err = es.indexNewDoc(taskIndex, bytes.NewReader(b)); err != nil {
		logger.Error("failed to create task", zap.Error(err))
		return
	}
	logger.Info("new task successfully created", zap.String("id", task.ID))
	return
}
