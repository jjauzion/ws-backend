package db

import (
	"context"
	"fmt"
	"github.com/jjauzion/ws-backend/internal"
	"go.uber.org/zap"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)


func NewClient(ctx context.Context) (es *elasticsearch7.Client, err error) {
	log := internal.GetLogger()
	cfg := internal.GetConfig()
	log.Info("connexion to ES cluster...")
	address := cfg.WS_ES_HOST + ":" + cfg.WS_ES_PORT
	esConfig := elasticsearch7.Config{
		Addresses: []string{address},
	}
	if es, err = elasticsearch7.NewClient(esConfig); err != nil {
		return
	}
		log.Error("couldn't connect to DB:", zap.Error(err))
	log.Info("successfully connected to ES", zap.String("host", address))
	info, _ := es.Info()
	log.Debug(fmt.Sprint(info))
	return
}