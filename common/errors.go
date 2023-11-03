package common

import "errors"

var (
	DbInternalError = errors.New("db internal error")
	ForbiddenError  = errors.New("forbidden")
)
