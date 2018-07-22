package clusters

import (
	"encoding/json"
	"github.com/betabandido/databricks-sdk-go/client"
	"github.com/betabandido/databricks-sdk-go/models"
)

type Endpoint struct {
	Client *client.Client
}

func (c *Endpoint) Create(request *models.ClustersCreateRequest) (*models.ClustersCreateResponse, error) {
	bytes, err := c.Client.Query("POST", "clusters/create", request)
	if err != nil {
		return nil, err
	}

	resp := models.ClustersCreateResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) List() (*models.ClustersListResponse, error) {
	bytes, err := c.Client.Query("GET", "clusters/list", nil)
	if err != nil {
		return nil, err
	}

	resp := models.ClustersListResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
