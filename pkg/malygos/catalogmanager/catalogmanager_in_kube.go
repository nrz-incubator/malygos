package catalogmanager

import (
	"context"
	"fmt"

	malygosv1 "github.com/nrz-incubator/malygos-controller/api/v1"
	"github.com/nrz-incubator/malygos/pkg/api"
	"github.com/nrz-incubator/malygos/pkg/errors"
	"github.com/nrz-incubator/malygos/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
)

type InKubeCatalogManager struct {
	client       *dynamic.DynamicClient
	cfgNamespace string
}

func NewInKubeCatalogManager(client *dynamic.DynamicClient, namespace string) (*InKubeCatalogManager, error) {
	return &InKubeCatalogManager{
		client:       client,
		cfgNamespace: namespace,
	}, nil
}

func (m *InKubeCatalogManager) ListComponents() ([]api.CatalogComponent, error) {
	unstructuredList, err := m.client.Resource(malygosv1.GroupVersion.WithResource("components")).
		Namespace(m.cfgNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var components []api.CatalogComponent
	for _, item := range unstructuredList.Items {
		var component api.CatalogComponent
		err := util.ConvertUnstructured(&item, &component)
		if err != nil {
			return nil, err
		}

		versions, err := m.listComponentVersions(component.Name)
		if err != nil {
			return nil, err
		}

		versionList := []string{}
		for _, componentVersionCRD := range versions {
			versionList = append(versionList, componentVersionCRD.Spec.Version)
		}
		component.Versions = &versionList

		components = append(components, component)
	}

	return components, nil
}

func (m *InKubeCatalogManager) AddComponent(component *api.CatalogComponent) error {
	componentCRD := malygosv1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      component.Name,
			Namespace: m.cfgNamespace,
		},
		Spec: malygosv1.ComponentSpec{
			Description:      component.Description,
			Icon:             component.Icon,
			ShortDescription: component.ShortDescription,
			URL:              component.Url,
			PublishedRegions: *component.PublishedRegions,
		},
	}

	unstructuredComponent, err := util.ConvertToUnstructured(componentCRD)
	if err != nil {
		return err
	}

	if component, err := m.GetComponent(component.Name); err != nil && !errors.IsNotFound(err) {
		return err
	} else if component != nil {
		return errors.NewConflictError("component", component.Name)
	}

	_, err = m.client.Resource(malygosv1.GroupVersion.WithResource("components")).
		Namespace(m.cfgNamespace).Create(context.Background(), unstructuredComponent, metav1.CreateOptions{})
	return err
}

func (m *InKubeCatalogManager) DeleteComponent(componentName string) error {
	err := m.client.Resource(malygosv1.GroupVersion.WithResource("components")).
		Namespace(m.cfgNamespace).Delete(context.Background(), componentName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.NewNotFoundError("component", componentName)
		}

		return err
	}

	return nil
}

func (m *InKubeCatalogManager) GetComponent(componentName string) (*api.CatalogComponent, error) {
	unstructuredComponent, err := m.client.Resource(malygosv1.GroupVersion.WithResource("components")).
		Namespace(m.cfgNamespace).Get(context.Background(), componentName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewNotFoundError("component", componentName)
		}

		return nil, err
	}

	componentCRD := &malygosv1.Component{}
	if err := util.ConvertUnstructured(unstructuredComponent, componentCRD); err != nil {
		return nil, err
	}

	component := &api.CatalogComponent{
		Name:             componentCRD.Name,
		Description:      componentCRD.Spec.Description,
		Icon:             componentCRD.Spec.Icon,
		ShortDescription: componentCRD.Spec.ShortDescription,
		Url:              componentCRD.Spec.URL,
		PublishedRegions: &componentCRD.Spec.PublishedRegions,
	}

	versions, err := m.listComponentVersions(componentName)
	if err != nil {
		return nil, err
	}

	versionList := []string{}

	for _, componentVersionCRD := range versions {
		versionList = append(versionList, componentVersionCRD.Spec.Version)
	}

	component.Versions = &versionList
	return component, nil
}

func (m *InKubeCatalogManager) GetComponentVersion(componentName string, version string) (*api.CatalogComponentVersion, error) {
	versions, err := m.listComponentVersions(componentName)
	if err != nil {
		return nil, err
	}

	for _, componentVersionCRD := range versions {
		if componentVersionCRD.Spec.Version == version {
			return &api.CatalogComponentVersion{
				Version:          componentVersionCRD.Spec.Version,
				DeprecationDate:  &componentVersionCRD.Spec.DeprecationDate,
				Description:      componentVersionCRD.Spec.Description,
				PublicationDate:  componentVersionCRD.Spec.PublicationDate,
				PublishedRegions: &componentVersionCRD.Spec.PublishedRegions,
				RemovalDate:      &componentVersionCRD.Spec.RemovalDate,
			}, nil
		}
	}

	return nil, errors.NewNotFoundError("component version", version)
}

func (m *InKubeCatalogManager) AddComponentVersion(componentName string, version *api.CatalogComponentVersion) error {
	if _, err := m.GetComponentVersion(componentName, version.Version); err != nil && !errors.IsNotFound(err) {
		return err
	}

	deprecationDate := ""
	if version.DeprecationDate != nil {
		deprecationDate = *version.DeprecationDate
	}

	removalDate := ""
	if version.RemovalDate != nil {
		removalDate = *version.RemovalDate
	}

	regions := []string{}
	if version.PublishedRegions != nil {
		regions = *version.PublishedRegions
	}

	componentVersionCRD := malygosv1.ComponentVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", componentName, version.Version),
			Namespace: m.cfgNamespace,
		},
		Spec: malygosv1.ComponentVersionSpec{
			ComponentRef:     componentName,
			Version:          version.Version,
			DeprecationDate:  deprecationDate,
			Description:      version.Description,
			PublicationDate:  version.PublicationDate,
			PublishedRegions: regions,
			RemovalDate:      removalDate,
		},
	}

	unstructuredComponentVersion, err := util.ConvertToUnstructured(componentVersionCRD)
	if err != nil {
		return err
	}

	_, err = m.client.Resource(malygosv1.GroupVersion.WithResource("componentversions")).
		Namespace(m.cfgNamespace).Create(context.Background(), unstructuredComponentVersion, metav1.CreateOptions{})
	return err
}

func (m *InKubeCatalogManager) DeleteComponentVersion(componentName string, version string) error {
	err := m.client.Resource(malygosv1.GroupVersion.WithResource("componentversions")).
		Namespace(m.cfgNamespace).Delete(context.Background(), fmt.Sprintf("%s-%s", componentName, version), metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.NewNotFoundError("component version", version)
		}

		return err
	}

	return nil
}

func (m *InKubeCatalogManager) listComponentVersions(componentName string) ([]malygosv1.ComponentVersion, error) {
	unstructuredList, err := m.client.Resource(malygosv1.GroupVersion.WithResource("componentversions")).
		Namespace(m.cfgNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var componentVersions []malygosv1.ComponentVersion
	for _, item := range unstructuredList.Items {
		var componentVersion malygosv1.ComponentVersion
		err := util.ConvertUnstructured(&item, &componentVersion)
		if err != nil {
			return nil, err
		}

		if componentVersion.Spec.ComponentRef != componentName {
			continue
		}

		componentVersions = append(componentVersions, componentVersion)
	}

	return componentVersions, nil
}

func (m *InKubeCatalogManager) SubscribeComponentVersion(region string, clusterId string, componentName string, componentVersion string) error {
	panic("implement me")
}

func (m *InKubeCatalogManager) ListComponentVersionSubscriptions(componentName string, componentVersion string) ([]api.SubscribedClusters, error) {
	panic("implement me")
}

func (m *InKubeCatalogManager) UnsubscribeComponentVersion(region string, clusterId string, componentName string, componentVersion string) error {
	panic("implement me")
}
