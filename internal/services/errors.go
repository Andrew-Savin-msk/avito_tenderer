package services

import "errors"

var (
	ErrNoSuchUser                  = errors.New("user doesn't exists")
	ErrUserExists                  = errors.New("user already exists")
	ErrConnectionLost              = errors.New("db connection lost")
	ErrServiceDatabaseDisconnected = errors.New("service database not available")
	ErrNothingToChange             = errors.New("nothing to change")
	ErrNoPermitions                = errors.New("user donn't have permitions")
	ErrNoSuchTender                = errors.New("tender doesn't exists")
	ErrNoSucnResource              = errors.New("no resource with such identifier")
	ErrNoSuchBid                   = errors.New("bid doesn't exists")
)
