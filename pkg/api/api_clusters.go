package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrz-incubator/malygos/pkg/errors"
)

func (api *ApiImpl) CreateCluster(c echo.Context) error {
	logger := api.logger
	if !api.manager.GetRBAC().IsAllowed("TODO username", "create", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	if c.Request().Body == nil {
		logger.Error(fmt.Errorf("create cluster request error"), "request body is nil")
		return c.JSON(http.StatusBadRequest, nil)
	}

	req, err := io.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		logger.Error(err, "failed to read request body on create cluster")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	cluster := &Cluster{}
	if err := json.Unmarshal(req, cluster); err != nil {
		logger.Error(err, "failed to unmarshal request body on create cluster")
		return c.JSON(http.StatusBadRequest, nil)
	}

	if err := cluster.ValidateInputs(); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Error: err.Error()})
	}

	clusterManager, err := api.manager.GetClusterManager(cluster.Region)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, Error{Error: fmt.Errorf("region %s not found", cluster.Region).Error()})
		}

		logger.Error(err, "failed to get cluster manager")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	cluster, err = clusterManager.Create(cluster)
	if err != nil {
		logger.Error(err, "failed to create cluster")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	logger.WithValues("region", cluster.Region, "name", cluster.Name).Info("cluster created")
	return c.JSON(http.StatusCreated, cluster)
}

func (api *ApiImpl) DeleteCluster(c echo.Context, region string, id string) error {
	logger := api.logger.WithValues("region", region, "id", id)
	if !api.manager.GetRBAC().IsAllowed("TODO username", "delete", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	clusterManager, err := api.manager.GetClusterManager(region)
	if err != nil {
		logger.Error(err, "failed to get cluster manager")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if err = clusterManager.Delete(id); err != nil {
		logger.Error(err, "failed to delete cluster")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	logger.Info("cluster deleted")
	return c.JSON(http.StatusNoContent, nil)
}

func (api *ApiImpl) GetCluster(c echo.Context, region string, id string) error {
	logger := api.logger.WithValues("region", region, "id", id)
	if !api.manager.GetRBAC().IsAllowed("TODO username", "get", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	clusterManager, err := api.manager.GetClusterManager(region)
	if err != nil {
		logger.Error(err, "failed to get cluster manager")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	cluster, err := clusterManager.Get(id)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		logger.Error(err, "failed to get cluster")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, cluster)
}

func (api *ApiImpl) ListClusters(c echo.Context) error {
	if !api.manager.GetRBAC().IsAllowed("TODO username", "list", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	registars, err := api.manager.GetClusterRegistrar().List()
	if err != nil {
		api.logger.Error(err, "failed to list management clusters")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := ListClustersResponse{
		JSON200: &struct {
			Clusters []Cluster "json:\"clusters\""
			Warnings *[]string "json:\"warnings,omitempty\""
		}{
			Clusters: []Cluster{},
			Warnings: nil,
		},
	}

	for _, registrar := range registars {
		kubeClient, err := registrar.CreateClient()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		dynKubeClient, err := registrar.CreateDynamicClient()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		clusterManager := api.manager.InstanciateClusterManager(api.logger, kubeClient, dynKubeClient)
		clusters, err := clusterManager.List()
		if err != nil {
			api.logger.Error(err, "failed to list clusters")
			*resp.JSON200.Warnings = append(*resp.JSON200.Warnings, err.Error())
		}

		for _, cluster := range clusters {
			resp.JSON200.Clusters = append(resp.JSON200.Clusters, *cluster)
		}

	}

	return c.JSON(200, resp.JSON200)
}
