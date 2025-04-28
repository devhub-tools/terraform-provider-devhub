package devhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) CreateDashboard(input Dashboard) (*Dashboard, error) {
	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/dashboards", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var dashboard Dashboard
	err = json.Unmarshal(body, &dashboard)
	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (c *Client) GetDashboard(id string) (*Dashboard, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/dashboards/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var dashboard Dashboard
	err = json.Unmarshal(body, &dashboard)
	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (c *Client) UpdateDashboard(dashboardId string, input Dashboard) (*Dashboard, error) {
	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/dashboards/%s", c.HostURL, dashboardId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var dashboard Dashboard
	err = json.Unmarshal(body, &dashboard)
	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (c *Client) DeleteDashboard(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/dashboards/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
