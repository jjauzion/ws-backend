package db

import (
	"fmt"
	"log"
	"strings"
	"context"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go.uber.org/zap"

	"github.com/jjauzion/ws-backend/internal"
)

type DatabaseHandler interface {
	New() error
	String() string
	Bootstrap() error
}


func NewDBHandler() DatabaseHandler {
	var dbh DatabaseHandler
	dbh = &esHandler{}
	if err := dbh.New(); err != nil {
		log.Fatal("", zap.Error(err))
	}
	return dbh
}

type esHandler struct {
	client *elasticsearch7.Client
}

func (es *esHandler) New() (err error) {
	log := internal.GetLogger()
	cfg := internal.GetConfig()
	log.Info("connexion to ES cluster...")
	address := cfg.WS_ES_HOST + ":" + cfg.WS_ES_PORT
	esConfig := elasticsearch7.Config{
		Addresses: []string{address},
	}
	es.client, err = elasticsearch7.NewClient(esConfig)
	if err != nil {
		log.Error("couldn't connect to ES:", zap.Error(err))
		return
	}
	log.Info("successfully connected to ES", zap.String("host", address))
	return
}

func (es *esHandler) String() string {
	info, _ := es.client.Info()
	s := fmt.Sprint(info)
	return s
}

func (es *esHandler) Bootstrap() error {
	log := internal.GetLogger()
	req := esapi.IndexRequest{
		Index:		"test",
		Body: 		strings.NewReader(`{"title": "Test1"`),
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		log.Error("failed to index doc:", zap.Error(err))
	}
	defer res.Body.Close()
	log.Info("index document", zap.String("status", res.Status()))
	return err
}

