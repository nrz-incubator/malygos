package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (api *ApiImpl) CreateManagementCluster(c echo.Context) error {
	if !api.rbac.IsAllowed("TODO username", "create", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	// TODO
	_, err := api.manager.GetClusterRegistrar().Create(nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, nil)
}

func (api *ApiImpl) ListManagementClusters(c echo.Context) error {
	if !api.rbac.IsAllowed("TODO username", "list", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	clusters, err := api.manager.GetClusterRegistrar().List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := ListManagementClustersResponse{
		JSON200: &struct {
			Clusters []ManagementCluster "json:\"clusters\""
			Warnings *[]string           "json:\"warnings,omitempty\""
		}{
			Clusters: []ManagementCluster{},
			Warnings: nil,
		},
	}

	for _, cluster := range clusters {
		resp.JSON200.Clusters = append(resp.JSON200.Clusters, ManagementCluster{
			Id:         &cluster.ID,
			Name:       cluster.Name,
			Kubeconfig: &cluster.Kubeconfig,
			Region:     cluster.Region,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (api *ApiImpl) GetManagementCluster(c echo.Context, id string) error {
	if !api.rbac.IsAllowed("TODO username", "get", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	_, err := api.manager.GetClusterRegistrar().Get(id)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.JSON(http.StatusInternalServerError, nil)
	}

	// TODO
	return c.JSON(http.StatusOK, nil)
}

func (api *ApiImpl) DeleteManagementCluster(c echo.Context, id string) error {
	if !api.rbac.IsAllowed("TODO username", "delete", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	if err := api.manager.GetClusterRegistrar().Delete(id); err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}
