package api

import "github.com/go-logr/logr"

type ApiImpl struct {
	logger logr.Logger
}

func NewApiImpl(logger logr.Logger) *ApiImpl {
	return &ApiImpl{
		logger: logger,
	}
}
