package db

import (
	"context"
	"fmt"
	"github.com/jjauzion/ws-backend/conf"
	logger "github.com/jjauzion/ws-backend/internal/logger"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"io/ioutil"
)

const (
	userIndex = "user"
	taskIndex = "task"
)

type esHandler struct {
	log     logger.Logger
	conf    conf.Configuration
	elastic *elastic.Client
}

func NewDatabaseAbstractedLayerImplemented(log logger.Logger, cf conf.Configuration) (Dbal, error) {
	var dbal Dbal
	dbal = &esHandler{
		log:  log,
		conf: cf,
	}

	err := dbal.newConnection(cf.WS_ES_HOST + ":" + cf.WS_ES_PORT)
	if err != nil {
		return nil, fmt.Errorf("cannot create new connection: %w", err)
	}

	return dbal, nil
}

func (es *esHandler) newConnection(url string) error {
	es.log.Info("connexion to ES cluster...")
	var err error
	es.elastic, err = elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		return err
	}

	es.log.Info("successfully connected to ES")
	return nil
}

func (es *esHandler) Info() string {

	return ""
}

func (es *esHandler) CreateIndexes(ctx context.Context) error {
	es.log.Info("Initializing Elasticsearch...")
	_, err := es.elastic.CreateIndex(taskIndex).Body(es.readMappingFile(es.conf.ES_TASK_MAPPING)).Do(ctx)
	if err != nil {
		es.log.Error("failed to create '"+taskIndex+"' index", zap.Error(err))
		return err
	}

	es.log.Info("'" + taskIndex + "' index created")
	_, err = es.elastic.CreateIndex(userIndex).Body(es.readMappingFile(es.conf.ES_USER_MAPPING)).Do(ctx)
	if err != nil {
		es.log.Error("failed to create '"+userIndex+" index", zap.Error(err))
		return err
	}

	es.log.Info("'" + userIndex + "' index created")
	es.log.Info("Elasticsearch successfully initialized !")
	return nil
}

type Itr func(*elastic.SearchHit) error

func (es *esHandler) elasticSearchOne(ctx context.Context, index string, source *elastic.SearchSource) (*elastic.SearchHit, error) {
	searchService := es.elastic.Search().Index(index).SearchSource(source)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return nil, err
	}

	if searchResult.TotalHits() <= 0 {
		return nil, ErrNotFound
	}
	if searchResult.TotalHits() > 1 {
		return nil, ErrTooManyRows
	}
	return searchResult.Hits.Hits[0], nil
}

func (es *esHandler) elasticSearch(ctx context.Context, index string, source *elastic.SearchSource, itr Itr) (*elastic.SearchResult, error) {
	searchService := es.elastic.Search().Index(index).SearchSource(source)
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

func (es *esHandler) readMappingFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		es.log.Fatal("cannot read mapping", zap.String("file", file), zap.Error(err))
		return ""
	}

	return string(content)
}
