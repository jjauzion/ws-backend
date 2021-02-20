package db

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"bytes"

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
	logger := internal.GetLogger()
	var dbh DatabaseHandler
	dbh = &esHandler{}
	if err := dbh.New(); err != nil {
		logger.Fatal("", zap.Error(err))
	}
	return dbh
}

type esHandler struct {
	client *elasticsearch7.Client
}

func (es *esHandler) New() (err error) {
	logger := internal.GetLogger()
	cfg := internal.GetConfig()
	logger.Info("connexion to ES cluster...")
	address := cfg.WS_ES_HOST + ":" + cfg.WS_ES_PORT
	esConfig := elasticsearch7.Config{
		Addresses: []string{address},
	}
	es.client, err = elasticsearch7.NewClient(esConfig)
	if err != nil {
		logger.Error("couldn't connect to ES:", zap.Error(err))
		return
	}
	logger.Info("successfully connected to ES", zap.String("host", address))
	return
}

func (es *esHandler) String() string {
	info, _ := es.client.Info()
	s := fmt.Sprint(info)
	return s
}

func (es *esHandler) CreateIndex(name, mappingFile string) (err error) {
	var mapping []byte
	if mapping, err = ioutil.ReadFile(mappingFile);  err != nil {
		return
	}
	req := esapi.IndicesCreateRequest{
		Index: name,
		Body: bytes.NewReader(mapping),
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client);  err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New(res.Status())
		return
	}
	return
}

func (es *esHandler) Bootstrap() (err error) {
	logger := internal.GetLogger()
	cf := internal.GetConfig()
	logger.Info("Initializing Elasticsearch...")
	index := "task"
	if err = es.CreateIndex(index, cf.ES_TASK_MAPPING); err != nil {
		logger.Error("failed to create '" + index + "' index: ", zap.Error(err))
		return
	}
	logger.Info("'" + index + "' index created")
	index = "user"
	if err = es.CreateIndex(index, cf.ES_USER_MAPPING); err != nil {
		logger.Error("failed to create '" + index + "' index: ", zap.Error(err))
		return
	}
	logger.Info("'" + index + "' index created")

	//req := esapi.IndexRequest{
	//	Index:		"test",
	//	Body: 		strings.NewReader(`{"title": "Test1"}`),
	//}
	//if err != nil {
	//	logger.Error("failed to index doc:", zap.Error(err))
	//	return err
	//}
	//defer res.Body.Close()
	//logger.Info("index document", zap.String("status", res.Status()))
	logger.Info("Elasticsearch successfully initialized !")
	return
}

