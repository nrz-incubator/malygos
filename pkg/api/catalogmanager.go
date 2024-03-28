package api

type CatalogManager interface {
	ListComponents() ([]CatalogComponent, error)
	AddComponent(component CatalogComponent) error
	DeleteComponent(componentName string) error
	GetComponent(componentName string) (CatalogComponent, error)
	AddComponentVersion(componentName string, version *CatalogComponentVersion) error
	DeleteComponentVersion(componentName string, version string) error
	GetComponentVersion(componentName string, version string) (CatalogComponentVersion, error)
}
