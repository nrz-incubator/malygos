package api

type CatalogManager interface {
	ListComponents() ([]CatalogComponent, error)
	AddComponent(component CatalogComponent) error
	DeleteComponent(componentName string) error
	GetComponent(componentName string) (CatalogComponent, error)
	AddComponentVersion(componentName string, version *CatalogComponentVersion) error
	DeleteComponentVersion(componentName string, version string) error
	GetComponentVersion(componentName string, version string) (CatalogComponentVersion, error)
	SubscribeComponentVersion(region string, clusterId string, componentName string, componentVersion string) error
	ListComponentVersionSubscriptions(componentName string, componentVersion string) ([]SubscribedClusters, error)
	UnsubscribeComponentVersion(region string, clusterId string, componentName string, componentVersion string) error
}
