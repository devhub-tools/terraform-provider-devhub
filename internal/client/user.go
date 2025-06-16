package devhub

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetUser(identifier string, lookupBy string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/lookup?%s=%s", c.HostURL, lookupBy, identifier), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
