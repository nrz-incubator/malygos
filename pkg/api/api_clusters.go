package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrz-incubator/malygos/pkg/malygos/clustermanager"
)

func (api *ApiImpl) CreateCluster(c echo.Context) error {
	logger := api.logger
	if !api.rbac.IsAllowed("TODO username", "create", "cluster") {
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

	if cluster.Kubeconfig == nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("kubeconfig field is required"))
	}

	if cluster.Name == "" {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("name field is required"))
	}

	clusterManager, err := api.manager.GetClusterManager(cluster.Region)
	if err != nil {
		logger.Error(err, "failed to get cluster manager")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	newCluster, err := clusterManager.Create(&clustermanager.Cluster{
		Name:       cluster.Name,
		Region:     cluster.Region,
		Kubeconfig: *cluster.Kubeconfig,
	})

	if err != nil {
		logger.Error(err, "failed to create cluster")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	logger.WithValues("region", cluster.Region, "name", cluster.Name).Info("cluster created")
	return c.JSON(http.StatusCreated, CreateClusterResponse{
		JSON201: &Cluster{
			Id:         &newCluster.ID,
			Name:       newCluster.Name,
			Kubeconfig: &newCluster.Kubeconfig,
			Region:     newCluster.Region,
			Status: &ClusterStatus{
				Phase:  "Pending",
				Online: false,
			},
		},
	})
}

func (api *ApiImpl) DeleteCluster(c echo.Context, region string, id string) error {
	logger := api.logger.WithValues("region", region, "id", id)
	if !api.rbac.IsAllowed("TODO username", "delete", "cluster") {
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
	if !api.rbac.IsAllowed("TODO username", "get", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	clusterManager, err := api.manager.GetClusterManager(region)
	if err != nil {
		logger.Error(err, "failed to get cluster manager")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	cluster, err := clusterManager.Get(id)
	if err != nil {
		logger.Error(err, "failed to get cluster")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if cluster == nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	resp := GetClusterResponse{
		JSON200: &Cluster{
			Id:         &cluster.ID,
			Name:       cluster.Name,
			Kubeconfig: &cluster.Kubeconfig,
			Region:     cluster.Region,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

func (api *ApiImpl) ListClusters(c echo.Context) error {
	if !api.rbac.IsAllowed("TODO username", "list", "cluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	mgmtClusters, err := api.manager.GetClusterRegistrar().List()
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

	for _, mgmtCluster := range mgmtClusters {
		kubeClient, err := mgmtCluster.CreateClient()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		// TODO, find a way to create this in a modular way
		clusterManager := clustermanager.NewKamajiClusterManager(api.logger, kubeClient)
		clusters, err := clusterManager.List()
		if err != nil {
			api.logger.Error(err, "failed to list clusters")
			*resp.JSON200.Warnings = append(*resp.JSON200.Warnings, err.Error())
		}

		for _, cluster := range clusters {
			resp.JSON200.Clusters = append(resp.JSON200.Clusters, Cluster{
				Id:         &cluster.ID,
				Name:       cluster.Name,
				Kubeconfig: &cluster.Kubeconfig,
				Region:     cluster.Region,
			})
		}

	}

	return c.JSON(200, resp)
}
