package db

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/jjauzion/ws-backend/graph/model"
	"github.com/jjauzion/ws-backend/internal"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
)

const (
	userIndex = "user"
	taskIndex = "task"
)


type esHandler struct {
	client *elasticsearch7.Client
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
	if err = es.bulkIngest(userIndex, "Elasticsearch/users.bulk", "true"); err != nil {
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

func (es *esHandler) search(index []string, request io.Reader) (res *esapi.Response, err error) {
	req := esapi.SearchRequest{
		Index: index,
		Body: request,
	}
	if res, err = req.Do(context.Background(), es.client);  err != nil {
		return
	}
	if res.IsError() {
		err = es.parseError(res)
		return
	}

	if res.IsError() {
		err = es.parseError(res)
		return
	}

	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return
	//}
	//fmt.Println("--------------------------------")
	//var data map[string]interface{}
	//if err = json.Unmarshal(body, &data); err != nil {
	//	return
	//}
	//fmt.Println("Here:", data)

	var r map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		err = errors.New(fmt.Sprintf("Error parsing the response body: %s", err))
		return
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data := hit.(map[string]interface{})["_source"]
		fmt.Println("----->", data)
		b, _ := json.Marshal(data)
		user := model.User{}
		err = json.Unmarshal(b, &user)
		if err != nil {
			return
		}
		fmt.Println("+++++", user)
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
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

