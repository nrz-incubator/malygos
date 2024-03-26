package managementclustermanager

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ManagementClusterManager interface {
	Create(cluster *ManagementCluster) (*ManagementCluster, error)
	Delete(region string) error
	List() ([]*ManagementCluster, error)
	Get(region string) (*ManagementCluster, error)
}

type ManagementCluster struct {
	ID         string
	Name       string
	Region     string
	restConfig *rest.Config
	Kubeconfig string
}

func (m *ManagementCluster) buildConfig() error {
	config, err := clientcmd.BuildConfigFromFlags("", m.Kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build k8s config: %v", err)
	}

	m.restConfig = config
	return nil
}

func (m *ManagementCluster) CreateClient() (*kubernetes.Clientset, error) {
	if m.restConfig == nil {
		if err := m.buildConfig(); err != nil {
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(m.restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %v", err)
	}

	_, err = client.DiscoveryClient.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to k8s cluster: %v", err)
	}

	return client, nil
}
