package clustermanager

import (
	"context"
	"encoding/json"
	"fmt"

	kamaji "github.com/clastix/kamaji/api/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/nrz-incubator/malygos/pkg/api"
	"github.com/nrz-incubator/malygos/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/utils/ptr"
)

const (
	regionalClusterLabel    = "malygos.local/region"
	clusterIDLabel          = "malygos.local/cluster-id"
	clusterRandomNameLength = 10
)

type KamajiClusterManager struct {
	client *dynamic.DynamicClient
	logger logr.Logger
}

func NewKamajiClusterManager(logger logr.Logger, client *dynamic.DynamicClient) *KamajiClusterManager {
	return &KamajiClusterManager{
		client: client,
		logger: logger,
	}
}

func generateClusterID() string {
	return fmt.Sprintf("malygos-%s", util.GenerateRandomString(clusterRandomNameLength))
}

func (m *KamajiClusterManager) Create(cluster *api.Cluster) (*api.Cluster, error) {
	clusterID := generateClusterID()
	kamajiCluster := &kamaji.TenantControlPlane{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kamaji.GroupVersion.String(),
			Kind:       "TenantControlPlane",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterID,
			Labels: map[string]string{
				regionalClusterLabel: cluster.Region,
			},
		},
		Spec: kamaji.TenantControlPlaneSpec{
			DataStore: "default",
			ControlPlane: kamaji.ControlPlane{
				Deployment: kamaji.DeploymentSpec{
					Replicas: ptr.To(int32(3)),
					Resources: &kamaji.ControlPlaneComponentsResources{
						APIServer:         &v1.ResourceRequirements{},
						ControllerManager: &v1.ResourceRequirements{},
						Scheduler:         &v1.ResourceRequirements{},
						Kine:              &v1.ResourceRequirements{},
					},
				},
				Service: kamaji.ServiceSpec{
					ServiceType: "LoadBalancer",
					AdditionalMetadata: kamaji.AdditionalMetadata{
						Labels: map[string]string{
							clusterIDLabel:       clusterID,
							regionalClusterLabel: cluster.Region,
						},
					},
				},
			},
			Kubernetes: kamaji.KubernetesSpec{
				Version: "v1.28.1",
				Kubelet: kamaji.KubeletSpec{
					CGroupFS: "systemd",
				},
			},
			NetworkProfile: kamaji.NetworkProfileSpec{},
			Addons: kamaji.AddonsSpec{
				CoreDNS: &kamaji.AddonSpec{},
			},
		},
	}

	b, err := json.Marshal(kamajiCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal kamaji cluster: %v", err)
	}

	unstructuredObj := &unstructured.Unstructured{}
	err = unstructuredObj.UnmarshalJSON(b)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal kamaji cluster: %v", err)
	}

	_, err = m.client.Resource(kamaji.GroupVersion.WithResource("tenantcontrolplanes")).
		Namespace("default").
		Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to create kamaji cluster: %v", err)
	}

	cluster.Id = &clusterID
	cluster.Status = &api.ClusterStatus{
		Phase:  "Pending",
		Online: false,
	}

	return cluster, nil
}

func (m *KamajiClusterManager) Delete(id string) error {
	err := m.client.Resource(kamaji.GroupVersion.WithResource("tenantcontrolplanes")).Namespace("default").Delete(context.TODO(), id, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete kamaji cluster: %v", err)
	}
	return nil
}

func (m *KamajiClusterManager) List() ([]*api.Cluster, error) {
	return nil, nil
}

func (m *KamajiClusterManager) Get(id string) (*api.Cluster, error) {
	return nil, nil
}
