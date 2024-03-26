package api

import (
	"github.com/go-logr/logr"
	"github.com/nrz-incubator/malygos/pkg/malygos/rbac"
)

type ApiImpl struct {
	logger  logr.Logger
	manager Manager
	rbac    rbac.RBAC
}

func NewApiImpl(logger logr.Logger,
	manager Manager) *ApiImpl {
	return &ApiImpl{
		logger:  logger,
		manager: manager,
	}
}
