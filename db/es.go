package db

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/jjauzion/ws-backend/internal"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
)

const (
	userIndex = "user"
	taskIndex = "task"
)


type esHandler struct {
	client *elasticsearch7.Client
}

type esSearchResponse struct {
	Took int						`json:"took"`
	Time_out bool					`json:"time_out"`
	Shards struct{
		Total int					`json:"total"`
		Successful int				`json:"sucessful"`
		Skipped int					`json:"skipped"`
		Failed int					`json:"failed"`
	}								`json:"_shards"`
	Hits struct{
		Total struct{
			Value int				`json:"value"`
			Relation string			`json:"relation"`
		}							`json:"total"`
		MaxScore float32			`json:"max_score"`
		Hits []struct{
			Index string			`json:"_index"`
			Type string				`json:"_type"`
			Id string				`json:"_id"`
			Score float32			`json:"_score"`
			Source interface{}		`json:"_source"`
		}							`json:"hits"`
	}								`json:"hits"`
}

func NewDBHandler() DatabaseHandler {
	logger := internal.GetLogger()
	var dbh DatabaseHandler
	dbh = &esHandler{}
	if err := dbh.new(); err != nil {
		logger.Fatal("", zap.Error(err))
	}
	return dbh
}

func (es *esHandler) Info() string {
	info, _ := es.client.Info()
	s := fmt.Sprint(info)
	return s
}

func (es *esHandler) Bootstrap() (err error) {
	logger := internal.GetLogger()
	cf := internal.GetConfig()
	logger.Info("Initializing Elasticsearch...")
	if err = es.createIndex(taskIndex, cf.ES_TASK_MAPPING); err != nil {
		logger.Error("failed to create '" + taskIndex + "' index: ", zap.Error(err))
		return
	}
	logger.Info("'" + taskIndex + "' index created")
	if err = es.createIndex(userIndex, cf.ES_USER_MAPPING); err != nil {
		logger.Error("failed to create '" + userIndex + "' index: ", zap.Error(err))
		return
	}
	logger.Info("'" + userIndex + "' index created")
	if err = es.bulkIngest(userIndex, cf.BOOTSTRAP_FILE, "true"); err != nil {
		logger.Error("failed to bulk ingest: ", zap.Error(err))
		return
	}
	logger.Info("Elasticsearch successfully initialized !")
	return
}

func (es *esHandler) new() (err error) {
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

func (es *esHandler) parseError(res *esapi.Response) (err error) {
	var e map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
		err = errors.New(fmt.Sprintf("error parsing the response body: %s", err))
	} else {
		err = errors.New(fmt.Sprintf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		))
	}
	return
}

func (es *esHandler) search(index []string, request io.Reader) (data []interface{}, err error) {
	logger := internal.GetLogger()
	req := esapi.SearchRequest{
		Index: index,
		Body: request,
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client);  err != nil {
		return
	}
	if res.IsError() {
		err = es.parseError(res)
		return
	}

	var r esSearchResponse
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		err = errors.New(fmt.Sprintf("Error parsing the response body: %s", err))
		return
	}
	if r.Shards.Failed > 0 {
		logger.Info("some shards failed during the search", zap.Int("number", r.Shards.Failed))
	}
	for i := 0; i < len(r.Hits.Hits); i++ {
		data = append(data, r.Hits.Hits[i].Source)
	}
	return
}

func (es *esHandler) createIndex(name, mappingFile string) (err error) {
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
		err = es.parseError(res)
		return
	}
	return
}

func (es *esHandler) indexNewDoc(index string, reader io.Reader) (err error) {
	req := esapi.IndexRequest{
		Index:		index,
		Body: 		reader,
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client);  err != nil {
		return
	}
	if res.IsError() {
		err = es.parseError(res)
		return
	}
	defer res.Body.Close()
	return
}


// WARNING: no error return if some doc failed
// Should use esutil to get a control on the nb of success / fail:
// https://github.com/elastic/go-elasticsearch/blob/master/_examples/bulk/indexer.go#L199
func (es *esHandler) bulkIngest(index, file, refresh string) (err error) {
	var bulk []byte
	if bulk, err = ioutil.ReadFile(file);  err != nil {
		return
	}
	req := esapi.BulkRequest{
		Index: index,
		Body: bytes.NewReader(bulk),
		Refresh: refresh,
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client);  err != nil {
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		err = es.parseError(res)
		return
	}
	return
}

