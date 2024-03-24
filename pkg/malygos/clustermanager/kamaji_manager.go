package clustermanager

import (
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

type KamajiClusterManager struct {
	client *kubernetes.Clientset
	logger logr.Logger
}

func NewKamajiClusterManager(logger logr.Logger, client *kubernetes.Clientset) *KamajiClusterManager {
	return &KamajiClusterManager{
		client: client,
		logger: logger,
	}
}

func (m *KamajiClusterManager) Create(cluster *Cluster) (*Cluster, error) {
	return cluster, nil
}

func (m *KamajiClusterManager) Delete(id string) error {
	return nil
}

func (m *KamajiClusterManager) List() ([]*Cluster, error) {
	return nil, nil
}

func (m *KamajiClusterManager) Get(id string) (*Cluster, error) {
	return nil, nil
}
