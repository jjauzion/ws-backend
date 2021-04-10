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
	userIndex = "ws_user" // Prefix all index with "ws_" as the workstation user have only access to those index
	taskIndex = "ws_task"
)

type esHandler struct {
	log    logger.Logger
	conf   conf.Configuration
	client *elastic.Client
}

func NewDatabaseAbstractedLayerImplemented(log logger.Logger, cf conf.Configuration) (Dbal, error) {
	var dbal Dbal
	dbal = &esHandler{
		log:  log,
		conf: cf,
	}

	err := dbal.NewConnection(cf.WS_ES_HOST + ":" + cf.WS_ES_PORT)
	if err != nil {
		return nil, fmt.Errorf("cannot create new connection: %w", err)
	}

	err = dbal.Ping()
	if err != nil {
		return nil, fmt.Errorf("elastic cluster is offline: %w", err)
	}

	return dbal, nil
}

func (es *esHandler) NewConnection(url string) error {
	es.log.Info("connexion to ES cluster...")
	var err error
	es.client, err = elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetBasicAuth(es.conf.WS_ES_USERNAME, es.conf.WS_ES_PWD))
	if err != nil {
		return err
	}

	return nil
}

func (es *esHandler) Ping() error {
	res, err := es.client.NodesInfo().Do(context.Background())
	if err == nil {
		es.log.Info("connected to ES cluster", zap.String("cluster_name", res.ClusterName))
	}
	return err
}

func (es *esHandler) CreateIndexes(ctx context.Context) error {
	es.log.Info("Initializing Elasticsearch...")
	mappingBody, err := es.readMappingFile(es.conf.ES_TASK_MAPPING)
	if err != nil {
		return fmt.Errorf("cannot read mapping file: %w", err)
	}

	_, err = es.client.CreateIndex(taskIndex).Body(mappingBody).Do(ctx)
	if err != nil {
		es.log.Error("failed to create '"+taskIndex+"' index", zap.Error(err))
		return err
	}

	mappingBody, err = es.readMappingFile(es.conf.ES_USER_MAPPING)
	if err != nil {
		return fmt.Errorf("cannot read mapping file: %w", err)
	}

	es.log.Info("'" + taskIndex + "' index created")
	_, err = es.client.CreateIndex(userIndex).Body(mappingBody).Do(ctx)
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
	searchService := es.client.Search().Index(index).SearchSource(source)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return nil, err
	}

	if searchResult.TotalHits() <= 0 {
		es.log.Debug(ErrNotFound.Error(), zap.Int64("hits", searchResult.TotalHits()))
		return nil, ErrNotFound
	}
	if searchResult.TotalHits() > 1 {
		es.log.Debug(ErrTooManyHits.Error(), zap.Int64("hits", searchResult.TotalHits()))
		for i, hit := range searchResult.Hits.Hits {
			es.log.Debug("hit", zap.Any(fmt.Sprintf("%d", i), hit.Source))
		}
		return nil, ErrTooManyHits
	}
	return searchResult.Hits.Hits[0], nil
}

func (es *esHandler) elasticSearch(ctx context.Context, index string, source *elastic.SearchSource, itr Itr) (*elastic.SearchResult, error) {
	searchService := es.client.Search().Index(index).SearchSource(source)
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

func (es *esHandler) readMappingFile(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		es.log.Error("cannot read mapping", zap.String("file", file), zap.Error(err))
		return "", err
	}

	return string(content), nil
}
