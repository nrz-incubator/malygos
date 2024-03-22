package clustermanager

type ClusterManager interface {
	Create(cluster *Cluster) (*Cluster, error)
	Delete(id string) error
	List() ([]*Cluster, error)
	Get(id string) (*Cluster, error)
}

type Cluster struct{}
