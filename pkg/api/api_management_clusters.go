package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (api *ApiImpl) CreateManagementCluster(c echo.Context) error {
	_, err := api.managementClusterManager.Create(nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, nil)
}

func (api *ApiImpl) ListManagementClusters(c echo.Context) error {
	_, err := api.managementClusterManager.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func (api *ApiImpl) GetManagementCluster(c echo.Context, id string) error {
	_, err := api.managementClusterManager.Get(id)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func (api *ApiImpl) DeleteManagementCluster(c echo.Context, id string) error {
	if err := api.managementClusterManager.Delete(id); err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
