package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrz-incubator/malygos/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
)

func (api *ApiImpl) CreateRegistrarCluster(c echo.Context) error {
	logger := api.logger
	if !api.manager.GetRBAC().IsAllowed("TODO username", "create", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	cluster := &RegistrarCluster{}
	if err := c.Bind(cluster); err != nil {
		logger.Error(err, "failed to bind request body on create cluster")
		return c.JSON(http.StatusBadRequest, nil)
	}

	if cluster.Kubeconfig == nil {
		return c.JSON(http.StatusBadRequest, Error{Error: "kubeconfig field is required"})
	}

	if len(*cluster.Kubeconfig) == 0 {
		return c.JSON(http.StatusBadRequest, Error{Error: "kubeconfig field must be non empty"})
	}

	if cluster.Name == "" {
		return c.JSON(http.StatusBadRequest, Error{Error: "name field is required"})
	}

	if cluster.Region == "" {
		return c.JSON(http.StatusBadRequest, Error{Error: "region field is required"})
	}

	// validate kubeconfig
	if _, err := clientcmd.NewClientConfigFromBytes([]byte(*cluster.Kubeconfig)); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Error: fmt.Errorf("kubeconfig is invalid: %v", err).Error()})
	}

	regCluster, err := api.manager.GetClusterRegistrar().Create(&ClusterRegistrar{
		Name:       cluster.Name,
		Region:     cluster.Region,
		Kubeconfig: *cluster.Kubeconfig,
	})
	if err != nil {
		if errors.IsConflict(err) {
			return c.JSON(http.StatusConflict, Error{Error: err.Error()})
		}

		logger.Error(err, "failed to create cluster")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, &RegistrarCluster{
		Id:         &regCluster.Id,
		Name:       regCluster.Name,
		Kubeconfig: &regCluster.Kubeconfig,
		Region:     regCluster.Region,
	})
}

func (api *ApiImpl) ListRegistrarClusters(c echo.Context) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "list", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	clusters, err := api.manager.GetClusterRegistrar().List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := ListRegistrarClustersResponse{
		JSON200: &struct {
			Clusters []RegistrarCluster "json:\"clusters\""
			Warnings *[]string          "json:\"warnings,omitempty\""
		}{
			Clusters: []RegistrarCluster{},
			Warnings: nil,
		},
	}

	for _, cluster := range clusters {
		resp.JSON200.Clusters = append(resp.JSON200.Clusters, RegistrarCluster{
			Id:         &cluster.Id,
			Name:       cluster.Name,
			Kubeconfig: &cluster.Kubeconfig,
			Region:     cluster.Region,
		})
	}

	return c.JSON(http.StatusOK, resp.JSON200)
}

func (api *ApiImpl) GetRegistrarCluster(c echo.Context, region string) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "get", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	cluster, err := api.manager.GetClusterRegistrar().Get(region)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, &RegistrarCluster{
		Id:         &cluster.Id,
		Name:       cluster.Name,
		Kubeconfig: &cluster.Kubeconfig,
		Region:     cluster.Region,
	})
}

func (api *ApiImpl) DeleteRegistrarCluster(c echo.Context, id string) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "delete", "managementcluster") {
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
