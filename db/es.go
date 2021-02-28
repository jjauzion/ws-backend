package db

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/jjauzion/ws-backend/conf"
	logger "github.com/jjauzion/ws-backend/internal/logger"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
)

const (
	userIndex = "user"
	taskIndex = "task"
)

type esHandler struct {
	client  *elasticsearch7.Client
	log     *logger.Logger
	cf      conf.Configuration
	elastic *elastic.Client
}

type esSearchResponse struct {
	Took     int  `json:"took"`
	Time_out bool `json:"time_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float32 `json:"max_score"`
		Hits     []struct {
			Index  string      `json:"_index"`
			Type   string      `json:"_type"`
			Id     string      `json:"_id"`
			Score  float32     `json:"_score"`
			Source interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func NewDBHandler(log *logger.Logger, cf conf.Configuration, elastic *elastic.Client) DatabaseHandler {
	var dbh DatabaseHandler
	dbh = &esHandler{
		log:     log,
		cf:      cf,
		elastic: elastic,
	}
	if err := dbh.new(); err != nil {
		log.Fatal("", zap.Error(err))
	}

	return dbh
}

func (es *esHandler) Info() string {
	info, _ := es.client.Info()
	s := fmt.Sprint(info)
	return s
}

func (es *esHandler) Bootstrap() (err error) {
	es.log.Info("Initializing Elasticsearch...")
	if err = es.createIndex(taskIndex, es.cf.ES_TASK_MAPPING); err != nil {
		es.log.Error("failed to create '"+taskIndex+"' index: ", zap.Error(err))
		return
	}
	es.log.Info("'" + taskIndex + "' index created")
	if err = es.createIndex(userIndex, es.cf.ES_USER_MAPPING); err != nil {
		es.log.Error("failed to create '"+userIndex+"' index: ", zap.Error(err))
		return
	}
	es.log.Info("'" + userIndex + "' index created")
	es.log.Info("Elasticsearch successfully initialized !")
	return
}

func (es *esHandler) new() error {
	es.log.Info("connexion to ES cluster...")
	address := es.cf.WS_ES_HOST + ":" + es.cf.WS_ES_PORT
	esConfig := elasticsearch7.Config{
		Addresses: []string{address},
	}
	var err error
	es.client, err = elasticsearch7.NewClient(esConfig)
	if err != nil {
		return err
	}
	_, err = es.client.Info()
	if err != nil {
		return err
	}
	es.log.Info("successfully connected to ES", zap.String("host", address))
	return err
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
	req := esapi.SearchRequest{
		Index: index,
		Body:  request,
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client); err != nil {
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
		es.log.Info("some shards failed during the search", zap.Int("number", r.Shards.Failed))
	}
	for i := 0; i < len(r.Hits.Hits); i++ {
		data = append(data, r.Hits.Hits[i].Source)
	}
	return
}

func (es *esHandler) createIndex(name, mappingFile string) (err error) {
	var mapping []byte
	if mapping, err = ioutil.ReadFile(mappingFile); err != nil {
		return
	}
	req := esapi.IndicesCreateRequest{
		Index: name,
		Body:  bytes.NewReader(mapping),
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client); err != nil {
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
		Index: index,
		Body:  reader,
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client); err != nil {
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
	if bulk, err = ioutil.ReadFile(file); err != nil {
		return
	}
	req := esapi.BulkRequest{
		Index:   index,
		Body:    bytes.NewReader(bulk),
		Refresh: refresh,
	}
	var res *esapi.Response
	if res, err = req.Do(context.Background(), es.client); err != nil {
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		err = es.parseError(res)
		return
	}
	return
}

type Itr func(*elastic.SearchHit) error

func (es *esHandler) elasticSearch(ctx context.Context, index string, query *elastic.MatchQuery, itr Itr) (*elastic.SearchResult, error) {
	searchSource := elastic.NewSearchSource()
	searchSource.Query(query)
	searchService := es.elastic.Search().Index(index).SearchSource(searchSource)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		err = itr(hit)
		if err != nil {
			return nil, err
		}
	}

	return searchResult, nil
}
