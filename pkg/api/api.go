package api

import (
	"github.com/go-logr/logr"
	"github.com/nrz-k8s-incubator/malygos/pkg/malygos/clustermanager"
	"github.com/nrz-k8s-incubator/malygos/pkg/malygos/managementclustermanager"
	"github.com/nrz-k8s-incubator/malygos/pkg/malygos/rbac"
)

type ApiImpl struct {
	logger                   logr.Logger
	clusterManager           clustermanager.ClusterManager
	managementClusterManager managementclustermanager.ManagementClusterManager
	rbac                     rbac.RBAC
}

func NewApiImpl(logger logr.Logger,
	clusterManager clustermanager.ClusterManager,
	managementClusterManager managementclustermanager.ManagementClusterManager) *ApiImpl {
	return &ApiImpl{
		logger:                   logger,
		clusterManager:           clusterManager,
		managementClusterManager: managementClusterManager,
	}
}
