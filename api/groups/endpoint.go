package groups

import (
	"encoding/json"

	"github.com/betabandido/databricks-sdk-go/client"
	"github.com/betabandido/databricks-sdk-go/models"
)

type Endpoint struct {
	Client *client.Client
}

func (c *Endpoint) AddMember(request *models.GroupsAddMemberRequest) error {
	_, err := c.Client.Query("POST", "groups/add-member", request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Endpoint) Create(request *models.GroupsCreateRequest) (*models.GroupsCreateResponse, error) {
	bytes, err := c.Client.Query("POST", "groups/create", request)
	if err != nil {
		return nil, err
	}

	resp := models.GroupsCreateResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) ListMembers(request *models.GroupsListMembersRequest) (*models.GroupsListMembersResponse, error) {
	bytes, err := c.Client.Query("GET", "groups/list-members", request)
	if err != nil {
		return nil, err
	}

	resp := models.GroupsListMembersResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) List() (*models.GroupsListResponse, error) {
	bytes, err := c.Client.Query("GET", "groups/list", nil)
	if err != nil {
		return nil, err
	}

	resp := models.GroupsListResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) ListParents(request *models.GroupsListParentsRequest) (*models.GroupsListParentsResponse, error) {
	bytes, err := c.Client.Query("GET", "groups/list-parents", request)
	if err != nil {
		return nil, err
	}

	resp := models.GroupsListParentsResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) RemoveMember(request *models.GroupsRemoveMemberRequest) error {
	_, err := c.Client.Query("POST", "groups/remove-member", request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Endpoint) Delete(request *models.GroupsDeleteRequest) error {
	_, err := c.Client.Query("POST", "groups/delete", request)
	if err != nil {
		return err
	}

	return nil
}
