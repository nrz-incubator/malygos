package managementclustermanager

type ManagementClusterManager interface {
	Create(cluster *ManagementCluster) (*ManagementCluster, error)
	Delete(id string) error
	List() ([]*ManagementCluster, error)
	Get(id string) (*ManagementCluster, error)
}

type ManagementCluster struct{}
