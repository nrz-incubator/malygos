package clustermanager

type KamajiClusterManager struct{}

func NewKamajiClusterManager() *KamajiClusterManager {
	return &KamajiClusterManager{}
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
