package storage

import (
	"errors"
)

var (
	ErrURLOpening   = errors.New("url is not correct")
	ErrDBConnection = errors.New("is not possible to make db connect")
)
