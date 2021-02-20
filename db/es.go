package db

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"bytes"
	"strings"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
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
	if res.IsError() {
		err = errors.New(res.String())
		return
	}
	return
}

func (es *esHandler) IndexNewDoc(index, doc string) (err error) {
	req := esapi.IndexRequest{
		Index:		index,
		Body: 		strings.NewReader(doc),
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client);  err != nil {
		return
	}
	if res.IsError() {
		err = errors.New(res.String())
		return
	}
	defer res.Body.Close()
	return
}


// WARNING: no error return if some doc failed
// Should use esutil to get a control on the nb of success / fail:
// https://github.com/elastic/go-elasticsearch/blob/master/_examples/bulk/indexer.go#L199
func (es *esHandler) BulkIngest(index, file string) (err error) {
	var bulk []byte
	if bulk, err = ioutil.ReadFile(file);  err != nil {
		return
	}
	req := esapi.BulkRequest{
		Index: index,
		Body: bytes.NewReader(bulk),
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client);  err != nil {
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		err = errors.New(res.String())
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
	if err = es.BulkIngest("user", "Elasticsearch/users.bulk"); err != nil {
		logger.Error("failed to bulk ingest: ", zap.Error(err))
		return
	}
	logger.Info("Elasticsearch successfully initialized !")
	return
}

