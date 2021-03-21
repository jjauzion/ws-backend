package db

import "fmt"

var (
	ErrNotFound    = fmt.Errorf("not found")
	ErrTooManyHits = fmt.Errorf("found too many hits")
)
