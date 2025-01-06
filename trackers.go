package redmine

import (
	"encoding/json"
	"errors"
	"strings"
)

type trackersResult struct {
	Trackers []IdName `json:"trackers"`
}

func (c *Client) Trackers() ([]IdName, error) {
	req, err := c.NewRequest("GET", "/trackers.json?"+c.getPaginationClause(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r trackersResult
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
	return r.Trackers, nil
}
