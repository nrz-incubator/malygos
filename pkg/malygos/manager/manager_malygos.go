package manager

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/nrz-incubator/malygos/pkg/api"
	"github.com/nrz-incubator/malygos/pkg/malygos/clustermanager"
	"github.com/nrz-incubator/malygos/pkg/malygos/clusterregistrar"
	"github.com/nrz-incubator/malygos/pkg/malygos/rbac"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MalygosManager struct {
	kubeConfig     *rest.Config
	clusterManager api.ClusterRegistrarManager
	logger         logr.Logger
	rbac           rbac.RBAC
}

func NewMalygosManager(logger logr.Logger, kubeconfig string, namespace string) (*MalygosManager, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s config: %v", err)
	}

	clusterManager, err := clusterregistrar.NewInKubeClusterManager(logger, config, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster manager: %v", err)
	}

	return &MalygosManager{
		kubeConfig:     config,
		clusterManager: clusterManager,
		logger:         logger,
		rbac:           rbac.NewNoop(),
	}, nil
}

func (m *MalygosManager) GetKubeconfig() *rest.Config {
	return m.kubeConfig
}

func (m *MalygosManager) GetClusterRegistrar() api.ClusterRegistrarManager {
	return m.clusterManager
}

func (m *MalygosManager) GetRBAC() rbac.RBAC {
	return m.rbac
}

func (m *MalygosManager) GetClusterManager(region string) (api.ClusterManager, error) {
	mgmtCluster, err := m.clusterManager.Get(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get management cluster for region %s: %v", region, err)
	}

	mgmtClusterClient, err := mgmtCluster.CreateClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client for management cluster: %v", err)
	}

	return m.InstanciateClusterManager(m.logger, mgmtClusterClient), nil
}

func (m *MalygosManager) InstanciateClusterManager(logger logr.Logger, clientset *kubernetes.Clientset) api.ClusterManager {
	return clustermanager.NewKamajiClusterManager(logger, clientset)
}
