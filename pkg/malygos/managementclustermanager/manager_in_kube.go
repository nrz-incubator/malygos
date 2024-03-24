package managementclustermanager

import (
	"k8s.io/client-go/kubernetes"
)

type InKubeClusterManager struct {
	client *kubernetes.Clientset
}

func NewInKubeClusterManager(client *kubernetes.Clientset) *InKubeClusterManager {
	return &InKubeClusterManager{
		client: client,
	}
}

func (m *InKubeClusterManager) Create(cluster *ManagementCluster) (*ManagementCluster, error) {
	return cluster, nil
}

func (m *InKubeClusterManager) Delete(id string) error {
	return nil
}

func (m *InKubeClusterManager) List() ([]*ManagementCluster, error) {
	return nil, nil
}

func (m *InKubeClusterManager) Get(id string) (*ManagementCluster, error) {
	return nil, nil
}
