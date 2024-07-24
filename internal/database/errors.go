package database

import "errors"

var (
	ErrRecordNotCreated = errors.New("record not created")
	ErrRecordNotFound   = errors.New("record not found")
)
