package devhub

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetRole(name string) (*Role, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/roles/lookup?name=%s", c.HostURL, name), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	role := Role{}
	err = json.Unmarshal(body, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}
