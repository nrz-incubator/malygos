package managementclustermanager

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	managementClusterSecretType = "malygos.local/management-cluster"
)

type InKubeClusterManager struct {
	client       *kubernetes.Clientset
	cfgNamespace string
	logger       logr.Logger
}

func NewInKubeClusterManager(logger logr.Logger, client *kubernetes.Clientset) *InKubeClusterManager {
	return &InKubeClusterManager{
		logger: logger,
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
	secrets, err := m.client.CoreV1().Secrets(m.cfgNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	clusters := make([]*ManagementCluster, 0)

	for _, secret := range secrets.Items {
		if secret.Type != managementClusterSecretType {
			continue
		}

		if _, ok := secret.Data["kubeconfig"]; !ok {
			m.logger.Error(fmt.Errorf("kubeconfig not found in secret"), "secret", m.cfgNamespace, secret.Name)
			continue
		}

		if _, ok := secret.Data["region"]; !ok {
			m.logger.Error(fmt.Errorf("region not found in secret"), "secret, m.cfgNamespace, secret.Name")
			continue
		}

		// config, err := clientcmd.RESTConfigFromKubeConfig(secret.Data["kubeconfig"])
		// if err != nil {
		// 	m.logger.Error(err, "failed to load kubeconfig from secret", "secret", m.cfgNamespace, secret.Name)
		// 	continue
		// }

		clusters = append(clusters, &ManagementCluster{
			ID:         string(secret.UID),
			Name:       secret.Name,
			Kubeconfig: string(secret.Data["kubeconfig"]),
			Region:     string(secret.Data["region"]),
		})
	}
	return clusters, nil
}

func (m *InKubeClusterManager) Get(id string) (*ManagementCluster, error) {

	return nil, nil
}
