package manager

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/nrz-incubator/malygos/pkg/api"
	"github.com/nrz-incubator/malygos/pkg/errors"
	"github.com/nrz-incubator/malygos/pkg/malygos/clustermanager"
	"github.com/nrz-incubator/malygos/pkg/malygos/clusterregistrar"
	"github.com/nrz-incubator/malygos/pkg/malygos/rbac"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MalygosManager struct {
	kubeConfig       *rest.Config
	registrarManager api.ClusterRegistrarManager
	logger           logr.Logger
	rbac             api.RBAC
	namespace        string
}

func NewMalygosManager(logger logr.Logger, kubeconfig string, namespace string) (*MalygosManager, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s config: %v", err)
	}

	registarManager, err := clusterregistrar.NewInKubeClusterManagerLegacy(logger, config, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster manager: %v", err)
	}

	return &MalygosManager{
		kubeConfig:       config,
		registrarManager: registarManager,
		logger:           logger,
		rbac:             rbac.NewNoop(),
		namespace:        namespace,
	}, nil
}

func (m *MalygosManager) GetKubeconfig() *rest.Config {
	return m.kubeConfig
}

func (m *MalygosManager) GetClusterRegistrar() api.ClusterRegistrarManager {
	return m.registrarManager
}

func (m *MalygosManager) GetRBAC() api.RBAC {
	return m.rbac
}

func (m *MalygosManager) GetClusterManager(region string) (api.ClusterManager, error) {
	registar, err := m.registrarManager.Get(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get management cluster for region %s: %v", region, err)
	}

	if registar == nil {
		return nil, errors.NewNotFoundError("management cluster for region", region)
	}

	client, err := registar.CreateDynamicClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client for management cluster: %v", err)
	}

	return m.InstanciateClusterManager(m.logger, nil, client), nil
}

func (m *MalygosManager) InstanciateClusterManager(logger logr.Logger, _ *kubernetes.Clientset, client *dynamic.DynamicClient) api.ClusterManager {
	return clustermanager.NewKamajiClusterManager(logger, client, m.namespace)
}

func (m *MalygosManager) GetCatalog() api.CatalogManager {
	// TODO implement a sample for this
	return nil
}
