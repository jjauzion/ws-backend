package db

import "fmt"

var (
	ErrNotFound    = fmt.Errorf("user not found")
	ErrTooManyHits = fmt.Errorf("found too many hits")
)
