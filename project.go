package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type projectRequest struct {
	Project Project `json:"project"`
}

type projectResult struct {
	Project Project `json:"project"`
}

type projectsResult struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	Id           int            `json:"id"`
	Parent       IdName         `json:"parent"`
	Name         string         `json:"name"`
	Identifier   string         `json:"identifier"`
	Description  string         `json:"description"`
	CreatedOn    string         `json:"created_on"`
	UpdatedOn    string         `json:"updated_on"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

func (c *Client) Project(id int) (*Project, error) {
	req, err := c.NewRequest("GET", "/projects/"+strconv.Itoa(id)+".json", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectResult
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
	return &r.Project, nil
}

func (c *Client) Projects() ([]Project, error) {
	req, err := c.NewRequest("GET", "/projects.json?"+c.getPaginationClause(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectsResult
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
	return r.Projects, nil
}

func (c *Client) CreateProject(project Project) (*Project, error) {
	var ir projectRequest
	ir.Project = project
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", "/projects.json", strings.NewReader(string(s)))
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
	var r projectRequest
	if res.StatusCode != 201 {
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
	return &r.Project, nil
}

func (c *Client) UpdateProject(project Project) error {
	var ir projectRequest
	ir.Project = project
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := c.NewRequest("PUT", "/projects/"+strconv.Itoa(project.Id)+".json", strings.NewReader(string(s)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		decoder := json.NewDecoder(res.Body)
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	if err != nil {
		return err
	}
	return err
}

func (c *Client) DeleteProject(id int) error {
	req, err := c.NewRequest("DELETE", "/projects/"+strconv.Itoa(id)+".json", strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}
