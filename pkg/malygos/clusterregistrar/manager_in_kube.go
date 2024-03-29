package clusterregistrar

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	malygosv1 "github.com/nrz-incubator/malygos/controller/api/v1"
	"github.com/nrz-incubator/malygos/pkg/api"
	"github.com/nrz-incubator/malygos/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type InKubeClusterManager struct {
	client       *kubernetes.Clientset
	cfgNamespace string
	logger       logr.Logger
}

func NewInKubeClusterManager(logger logr.Logger, config *rest.Config, namespace string) (*InKubeClusterManagerLegacy, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %v", err)
	}

	return &InKubeClusterManagerLegacy{
		logger:       logger,
		client:       client,
		cfgNamespace: namespace,
	}, nil
}

func (m *InKubeClusterManager) Create(cluster *api.ClusterRegistrar) (*api.ClusterRegistrar, error) {
	mc, err := m.Get(cluster.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to get management cluster: %v", err)
	}

	if mc != nil {
		return nil, errors.NewConflictError("region", cluster.Region)
	}

	registrar := malygosv1.Registrar{}

	secretName := generateSecretName()
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
			Labels: map[string]string{
				managementClusterLabelRegion: cluster.Region,
				managementClusterLabelName:   cluster.Name,
			},
		},
		Type: managementClusterSecretType,
		Data: map[string][]byte{
			"kubeconfig": []byte(cluster.Kubeconfig),
		},
	}

	secret, err = m.client.CoreV1().Secrets(m.cfgNamespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	cluster.Id = secret.Name
	return cluster, err
}

func (m *InKubeClusterManager) Delete(region string) error {
	cluster, err := m.Get(region)
	if err != nil {
		return err
	}

	if cluster == nil {
		return errors.NewNotFoundError("region", region)
	}

	err = m.client.CoreV1().Secrets(m.cfgNamespace).Delete(context.TODO(), cluster.Id, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (m *InKubeClusterManager) List() ([]*api.ClusterRegistrar, error) {
	secrets, err := m.client.CoreV1().Secrets(m.cfgNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	clusters := make([]*api.ClusterRegistrar, 0)

	for _, secret := range secrets.Items {
		mc, err := decodeSecret(&secret)
		if err != nil {
			if err, ok := err.(*secretTypeError); !ok {
				m.logger.Error(err, "failed to decode secret", "secret", secret.Name)
			}
			continue
		}

		clusters = append(clusters, mc)
	}
	return clusters, nil
}

func (m *InKubeClusterManager) Get(region string) (*api.ClusterRegistrar, error) {
	clusters, err := m.List()
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {
		if cluster.Region == region {
			return cluster, nil
		}
	}

	return nil, nil
}
