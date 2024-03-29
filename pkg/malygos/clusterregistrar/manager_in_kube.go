package clusterregistrar

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	malygosv1 "github.com/nrz-incubator/malygos-controller/api/v1"
	"github.com/nrz-incubator/malygos/pkg/api"
	"github.com/nrz-incubator/malygos/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type InKubeClusterManager struct {
	client       *dynamic.DynamicClient
	cfgNamespace string
	logger       logr.Logger
}

func NewInKubeClusterManager(logger logr.Logger, config *rest.Config, namespace string) (*InKubeClusterManager, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	return &InKubeClusterManager{
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

	registrar := malygosv1.Registrar{
		TypeMeta: metav1.TypeMeta{
			APIVersion: malygosv1.GroupVersion.String(),
			Kind:       "Registrar",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: m.cfgNamespace,
		},
		Spec: malygosv1.RegistrarSpec{
			Region:     cluster.Region,
			Kubeconfig: cluster.Kubeconfig,
		},
	}

	b, err := json.Marshal(registrar)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal registrar: %v", err)
	}

	unstructuredObj := &unstructured.Unstructured{}
	if err = unstructuredObj.UnmarshalJSON(b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal registrar cluster: %v", err)
	}

	if _, err := m.client.Resource(malygosv1.GroupVersion.WithResource("registrars")).
		Namespace(m.cfgNamespace).
		Create(context.TODO(), unstructuredObj, metav1.CreateOptions{}); err != nil {
		return nil, fmt.Errorf("failed to create registrar: %v", err)
	}

	cluster.Id = registrar.Name
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

	if err := m.client.Resource(malygosv1.GroupVersion.WithResource("registrars")).
		Namespace(m.cfgNamespace).
		Delete(context.TODO(), cluster.Name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete registrar: %v", err)
	}
	return nil
}

func (m *InKubeClusterManager) List() ([]*api.ClusterRegistrar, error) {
	unstructuredList, err := m.client.Resource(malygosv1.GroupVersion.WithResource("registrars")).
		Namespace(m.cfgNamespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to list registrars clusters: %v", err)
	}

	b, err := json.Marshal(unstructuredList)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal registrars cluster list: %v", err)
	}

	var registrars malygosv1.RegistrarList
	err = json.Unmarshal(b, &registrars)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal registrars cluster list: %v", err)
	}

	clusters := make([]*api.ClusterRegistrar, 0)

	for _, registar := range registrars.Items {
		clusters = append(clusters, &api.ClusterRegistrar{
			Id:         registar.Name,
			Name:       registar.Name,
			Region:     registar.Spec.Region,
			Kubeconfig: registar.Spec.Kubeconfig,
		})
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
