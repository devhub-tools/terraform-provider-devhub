package devhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) CreateWorkflow(input Workflow) (*Workflow, error) {
	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/workflows", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var workflow Workflow
	err = json.Unmarshal(body, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

func (c *Client) GetWorkflow(id string) (*Workflow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/workflows/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var workflow Workflow
	err = json.Unmarshal(body, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

func (c *Client) UpdateWorkflow(workflowId string, input Workflow) (*Workflow, error) {
	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/workflows/%s", c.HostURL, workflowId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var workflow Workflow
	err = json.Unmarshal(body, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

func (c *Client) DeleteWorkflow(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/workflows/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
