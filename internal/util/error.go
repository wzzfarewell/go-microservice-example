package util

import "errors"

var (
	ErrPathParamNotFound = errors.New("unknown argument passed")

	ErrInvalidArgument = errors.New("invalid argument passed")
)
