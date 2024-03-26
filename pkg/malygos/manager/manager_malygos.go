package manager

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/nrz-incubator/malygos/pkg/malygos/clustermanager"
	"github.com/nrz-incubator/malygos/pkg/malygos/managementclustermanager"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MalygosManager struct {
	kubeConfig     *rest.Config
	clusterManager managementclustermanager.ManagementClusterManager
	logger         logr.Logger
}

func NewMalygosManager(logger logr.Logger, kubeconfig string) (*MalygosManager, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s config: %v", err)
	}

	clusterManager, err := managementclustermanager.NewInKubeClusterManager(logger, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster manager: %v", err)
	}

	return &MalygosManager{
		kubeConfig:     config,
		clusterManager: clusterManager,
		logger:         logger,
	}, nil
}

func (m *MalygosManager) GetKubeconfig() *rest.Config {
	return m.kubeConfig
}

func (m *MalygosManager) GetManagementClusterManager() managementclustermanager.ManagementClusterManager {
	return m.clusterManager
}

func (m *MalygosManager) GetClusterManager(region string) (clustermanager.ClusterManager, error) {
	mgmtCluster, err := m.clusterManager.Get(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get management cluster for region %s: %v", region, err)
	}

	mgmtClusterClient, err := mgmtCluster.CreateClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client for management cluster: %v", err)
	}

	return clustermanager.NewKamajiClusterManager(m.logger, mgmtClusterClient), nil
}
