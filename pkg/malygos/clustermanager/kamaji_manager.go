package clustermanager

import (
	"context"
	"fmt"

	kamaji "github.com/clastix/kamaji/api/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/nrz-incubator/malygos/pkg/api"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/ptr"
)

const (
	regionalClusterLabel = "malygos.local/region"
	clusterIDLabel       = "malygos.local/cluster-id"
)

type KamajiClusterManager struct {
	client *kubernetes.Clientset
	logger logr.Logger
}

func NewKamajiClusterManager(logger logr.Logger, client *kubernetes.Clientset) *KamajiClusterManager {
	return &KamajiClusterManager{
		client: client,
		logger: logger,
	}
}

func (m *KamajiClusterManager) Create(cluster *api.Cluster) (*api.Cluster, error) {
	kamajiCluster := &kamaji.TenantControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name: *cluster.Id,
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
					AdditionalMetadata: kamaji.AdditionalMetadata{
						Labels: map[string]string{
							clusterIDLabel:       *cluster.Id,
							regionalClusterLabel: cluster.Region,
						},
					},
				},
			},
			Kubernetes: kamaji.KubernetesSpec{
				Version: "1.28.1",
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

	err := m.client.RESTClient().Post().Resource("tenantcontrolplanes").Body(kamajiCluster).Do(context.TODO()).Into(kamajiCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create kamaji cluster: %v", err)
	}

	cluster.Status = &api.ClusterStatus{
		Phase:  "Pending",
		Online: false,
	}

	return cluster, nil
}

func (m *KamajiClusterManager) Delete(id string) error {
	err := m.client.RESTClient().Delete().Resource("tenantcontrolplanes").Name(id).Do(context.TODO()).Into(&kamaji.TenantControlPlane{})
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
