package devhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetWorkspace(workspaceId string) (*TerradeskWorkspace, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/terradesk/workspaces/%s", c.HostURL, workspaceId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := TerradeskWorkspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (c *Client) CreateWorkspace(input TerradeskWorkspace) (*TerradeskWorkspace, error) {
	if input.EnvVars == nil {
		input.EnvVars = make([]EnvVar, 0)
	}

	if input.Secrets == nil {
		input.Secrets = make([]Secret, 0)
	}

	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/terradesk/workspaces", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := TerradeskWorkspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (c *Client) UpdateWorkspace(workspaceId string, input TerradeskWorkspace) (*TerradeskWorkspace, error) {
	if input.EnvVars == nil {
		input.EnvVars = make([]EnvVar, 0)
	}

	if input.Secrets == nil {
		input.Secrets = make([]Secret, 0)
	}

	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/terradesk/workspaces/%s", c.HostURL, workspaceId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workspace := TerradeskWorkspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (c *Client) DeleteWorkspace(workspaceId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/terradesk/workspaces/%s", c.HostURL, workspaceId), nil)
	if err != nil {
		return err
	}

	if _, err := c.doRequest(req); err != nil {
		return err
	}

	return nil
}
