package entity

import "errors"

var (
	RepoNotFoundErr  = errors.New("No rows")
	RequestStatusErr = errors.New("Request status validation error")
)
