package managementclustermanager

import "k8s.io/client-go/rest"

type ManagementClusterManager interface {
	Create(cluster *ManagementCluster) (*ManagementCluster, error)
	Delete(id string) error
	List() ([]*ManagementCluster, error)
	Get(id string) (*ManagementCluster, error)
}

type ManagementCluster struct {
	ID         string
	Name       string
	Region     string
	restConfig *rest.Config
	Kubeconfig string
}
