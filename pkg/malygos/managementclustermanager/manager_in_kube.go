package managementclustermanager

type InKubeClusterManager struct{}

func NewInKubeClusterManager() *InKubeClusterManager {
	return &InKubeClusterManager{}
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
