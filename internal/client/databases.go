package devhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetDatabase(databaseId string) (*Database, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/querydesk/databases/%s", c.HostURL, databaseId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	database := Database{}
	err = json.Unmarshal(body, &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

func (c *Client) CreateDatabase(input Database) (*Database, error) {
	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/querydesk/databases", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	database := Database{}
	err = json.Unmarshal(body, &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

func (c *Client) UpdateDatabase(databaseId string, input Database) (*Database, error) {
	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/querydesk/databases/%s", c.HostURL, databaseId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	database := Database{}
	err = json.Unmarshal(body, &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

func (c *Client) DeleteDatabase(databaseId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/querydesk/databases/%s", c.HostURL, databaseId), nil)
	if err != nil {
		return err
	}

	if _, err := c.doRequest(req); err != nil {
		return err
	}

	return nil
}
