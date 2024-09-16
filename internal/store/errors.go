package store

import "errors"

var (
	ErrRecordNotFound      = errors.New("no such record")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrStartingTransaction = errors.New("unable to start transaction")
	ErrConnClosed          = errors.New("connection closed")
	ErrUserNotFound        = errors.New("no such username in db")
)
