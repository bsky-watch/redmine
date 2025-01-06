package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type customFieldsResult struct {
	CustomFields []CustomField `json:"custom_fields"`
}

// CustomFields consulta los campos personalizados
func (c *Client) CustomFields() ([]CustomField, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/custom_fields.json?%s", c.getPaginationClause()), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r customFieldsResult
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.CustomFields, nil
}
