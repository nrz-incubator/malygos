package manager

import (
	"github.com/nrz-incubator/malygos/pkg/malygos/clustermanager"
	"github.com/nrz-incubator/malygos/pkg/malygos/clusterregistrar"
	"k8s.io/client-go/rest"
)

type Manager interface {
	GetKubeconfig() *rest.Config
	GetClusterRegistrar() clusterregistrar.ClusterRegistrarManager
	GetClusterManager(region string) (clustermanager.ClusterManager, error)
}
