package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrz-incubator/malygos/pkg/errors"
)

func (api *ApiImpl) ListCatalogComponents(c echo.Context) error {
	// rbac check
	if !api.manager.GetRBAC().IsAllowed("TODO username", "list", "catalog_component") {
		return c.JSON(http.StatusForbidden, nil)
	}

	catalog, err := api.manager.GetCatalog().ListComponents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
	}

	return c.JSON(http.StatusNotImplemented, catalog)
}

func (api *ApiImpl) AddCatalogComponent(c echo.Context) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "create", "catalog_component") {
		return c.JSON(http.StatusForbidden, nil)
	}

	var component CatalogComponent
	if err := c.Bind(&component); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	// TODO validate component
	err := api.manager.GetCatalog().AddComponent(component)
	if err != nil {
		if errors.IsConflict(err) {
			return c.JSON(http.StatusConflict, nil)
		}

		return c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
	}

	return c.JSON(http.StatusNotImplemented, component)
}

func (api *ApiImpl) DeleteCatalogComponent(c echo.Context, componentName string) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "delete", "catalog_component") {
		return c.JSON(http.StatusForbidden, nil)
	}

	err := api.manager.GetCatalog().DeleteComponent(componentName)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
	}

	return c.JSON(http.StatusNotImplemented, CatalogComponent{})
}

func (api *ApiImpl) GetCatalogComponent(c echo.Context, componentName string) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "get", "catalog_component") {
		return c.JSON(http.StatusForbidden, nil)
	}

	component, err := api.manager.GetCatalog().GetComponent(componentName)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
	}

	return c.JSON(http.StatusNotImplemented, component)
}

func (api *ApiImpl) AddCatalogComponentVersion(c echo.Context, componentName string) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "create", "catalog_component_version") {
		return c.JSON(http.StatusForbidden, nil)
	}

	var componentVersion CatalogComponentVersion
	if err := c.Bind(&componentVersion); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	// TODO validate componentVersion
	err := api.manager.GetCatalog().AddComponentVersion(componentName, &componentVersion)
	if err != nil {
		if errors.IsConflict(err) {
			return c.JSON(http.StatusConflict, nil)
		}

		return c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
	}

	return c.JSON(http.StatusNotImplemented, componentVersion)
}

func (api *ApiImpl) GetCatalogComponentVersion(c echo.Context, componentName string, componentVersion string) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "get", "catalog_component_version") {
		return c.JSON(http.StatusForbidden, nil)
	}

	version, err := api.manager.GetCatalog().GetComponentVersion(componentName, componentVersion)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
	}

	return c.JSON(http.StatusNotImplemented, version)
}

func (api *ApiImpl) DeleteCatalogComponentVersion(c echo.Context, componentName string, componentVersion string) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "delete", "catalog_component_version") {
		return c.JSON(http.StatusForbidden, nil)
	}

	err := api.manager.GetCatalog().DeleteComponentVersion(componentName, componentVersion)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, Error{Error: err.Error()})
	}

	return c.JSON(http.StatusNotImplemented, nil)
}
