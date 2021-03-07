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
	//client  *elasticsearch7.Client
	log     *logger.Logger
	cf      conf.Configuration
	elastic *elastic.Client
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
	return ""
}

func (es *esHandler) Bootstrap(ctx context.Context) error {
	es.log.Info("Initializing Elasticsearch...")
	_, err := es.elastic.CreateIndex(taskIndex).Body(es.getMapping(es.cf.ES_TASK_MAPPING)).Do(ctx)
	if err != nil {
		es.log.Error("failed to create '"+taskIndex+"' index", zap.Error(err))
		return err
	}

	es.log.Info("'" + taskIndex + "' index created")
	_, err = es.elastic.CreateIndex(userIndex).Body(es.getMapping(es.cf.ES_USER_MAPPING)).Do(ctx)
	if err != nil {
		es.log.Error("failed to create '"+userIndex+" index", zap.Error(err))
		return err
	}

	es.log.Info("'" + userIndex + "' index created")
	es.log.Info("Elasticsearch successfully initialized !")
	return nil
}

func (es *esHandler) new() error {
	es.log.Info("connexion to ES cluster...")

	es.log.Info("successfully connected to ES")
	return nil
}

type Itr func(*elastic.SearchHit) error

func (es *esHandler) elasticSearchOne(ctx context.Context, index string, query *elastic.MatchQuery) (*elastic.SearchHit, error) {
	searchSource := elastic.NewSearchSource()
	searchSource.Query(query)
	searchService := es.elastic.Search().Index(index).SearchSource(searchSource)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return nil, err
	}
	if searchResult.TotalHits() == 0 {
		return nil, ErrNotFound
	}
	if searchResult.TotalHits() > 1 {
		return nil, ErrTooManyRows
	}
	if searchResult.TotalHits() == 1 {
		return searchResult.Hits.Hits[0], nil
	}
	return nil, fmt.Errorf("something wrong happened")
}

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

func (es *esHandler) getMapping(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		es.log.Fatal("cannot read mapping", zap.String("file", file), zap.Error(err))
		return ""
	}

	return string(content)
}
