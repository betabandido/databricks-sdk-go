package clusters

import (
	"encoding/json"
	"github.com/betabandido/databricks-sdk-go/client"
	"github.com/betabandido/databricks-sdk-go/models"
)

type Endpoint struct {
	Client *client.Client
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
