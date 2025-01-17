package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type newsResult struct {
	News []News `json:"news"`
}

type News struct {
	Id          int    `json:"id"`
	Project     IdName `json:"project"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	CreatedOn   string `json:"created_on"`
}

func (c *Client) News(projectId int) ([]News, error) {
	req, err := c.NewRequest("GET", "/projects/"+strconv.Itoa(projectId)+"/news.json?"+c.getPaginationClause(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r newsResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
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
	return r.News, nil
}
