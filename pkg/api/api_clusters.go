package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (api *ApiImpl) CreateCluster(c echo.Context) error {
	if !api.rbac.IsAllowed("TODO username", "create", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	_, err := api.clusterManager.Create(nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, nil)
}

func (api *ApiImpl) DeleteCluster(c echo.Context, id string) error {
	if !api.rbac.IsAllowed("TODO username", "delete", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	if err := api.clusterManager.Delete(id); err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (api *ApiImpl) GetCluster(c echo.Context, id string) error {
	if !api.rbac.IsAllowed("TODO username", "get", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	_, err := api.clusterManager.Get(id)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(200, nil)
}

func (api *ApiImpl) ListClusters(c echo.Context) error {
	if !api.rbac.IsAllowed("TODO username", "list", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	_, err := api.clusterManager.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(200, nil)
}
