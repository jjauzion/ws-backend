package db

import "fmt"

var (
	ErrNotFound    = fmt.Errorf("not found")
	ErrTooManyRows = fmt.Errorf("found too many hits")
)
