package api

import (
	"github.com/go-logr/logr"
	"github.com/nrz-k8s-incubator/malygos/pkg/malygos/clustermanager"
	"github.com/nrz-k8s-incubator/malygos/pkg/malygos/managementclustermanager"
)

type ApiImpl struct {
	logger                   logr.Logger
	clusterManager           clustermanager.ClusterManager
	managementClusterManager managementclustermanager.ManagementClusterManager
}

func NewApiImpl(logger logr.Logger) *ApiImpl {
	return &ApiImpl{
		logger: logger,
	}
}
