package api

import (
	"github.com/go-logr/logr"
)

type ApiImpl struct {
	logger  logr.Logger
	manager Manager
}

func NewApiImpl(logger logr.Logger,
	manager Manager) *ApiImpl {
	return &ApiImpl{
		logger:  logger,
		manager: manager,
	}
}
