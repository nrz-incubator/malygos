package manager

import (
	"github.com/nrz-incubator/malygos/pkg/malygos/clustermanager"
	"github.com/nrz-incubator/malygos/pkg/malygos/managementclustermanager"
	"k8s.io/client-go/rest"
)

type Manager interface {
	GetKubeconfig() *rest.Config
	GetManagementClusterManager() managementclustermanager.ManagementClusterManager
	GetClusterManager(region string) (clustermanager.ClusterManager, error)
}
