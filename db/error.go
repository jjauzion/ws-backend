package db

import "fmt"

var (
	ErrNotFound    = fmt.Errorf("not found")
	ErrTooManyHits = fmt.Errorf("found too many hits")
)

type ErrAlreadyExist string

func (e ErrAlreadyExist) Error() string {
	return string(e)
}
func (e ErrAlreadyExist) Ptr() *ErrAlreadyExist {
	return &e
}
