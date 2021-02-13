package db

import (
	"context"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)


func NewClient(ctx context.Context) (es *elasticsearch7.Client, err error) {
	if es, err = elasticsearch7.NewDefaultClient(); err != nil {
		return
	}
	return
}