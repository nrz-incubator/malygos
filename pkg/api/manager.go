package api

import (
	"github.com/go-logr/logr"
	"github.com/nrz-incubator/malygos/pkg/malygos/rbac"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Manager interface {
	GetKubeconfig() *rest.Config
	GetClusterRegistrar() ClusterRegistrarManager

	InstanciateClusterManager(logr.Logger, *kubernetes.Clientset) ClusterManager
	GetClusterManager(region string) (ClusterManager, error)

	GetRBAC() rbac.RBAC
}
