package redmine

import (
	"encoding/json"
	"errors"
	"strings"
)

type issuePrioritiesResult struct {
	IssuePriorities []IssuePriority `json:"issue_priorities"`
}

type IssuePriority struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

func (c *Client) IssuePriorities() ([]IssuePriority, error) {
	req, err := c.NewRequest("GET", "/enumerations/issue_priorities.json?"+c.getPaginationClause(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issuePrioritiesResult
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
	return r.IssuePriorities, nil
}
