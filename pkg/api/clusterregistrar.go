package api

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterRegistrarManager interface {
	Create(cluster *ClusterRegistrar) (*ClusterRegistrar, error)
	Delete(region string) error
	List() ([]*ClusterRegistrar, error)
	Get(region string) (*ClusterRegistrar, error)
}

type ClusterRegistrar struct {
	Id         string
	Name       string
	Region     string
	restConfig *rest.Config
	Kubeconfig string
}

func (m *ClusterRegistrar) buildConfig() error {
	clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(m.Kubeconfig))
	if err != nil {
		return fmt.Errorf("failed to build k8s config: %v", err)
	}
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return fmt.Errorf("failed to get k8s client config: %v", err)
	}

	m.restConfig = config
	return nil
}

func (m *ClusterRegistrar) CreateClient() (*kubernetes.Clientset, error) {
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

func (m *ClusterRegistrar) CreateDynamicClient() (*dynamic.DynamicClient, error) {
	if m.restConfig == nil {
		if err := m.buildConfig(); err != nil {
			return nil, err
		}
	}

	client, err := dynamic.NewForConfig(m.restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	return client, nil
}
