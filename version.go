package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type versionRequest struct {
	Version Version `json:"version"`
}

type versionResult struct {
	Version Version `json:"version"`
}

type versionsResult struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	Id           int            `json:"id"`
	Project      IdName         `json:"project"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Status       string         `json:"status"`
	DueDate      string         `json:"due_date"`
	CreatedOn    string         `json:"created_on"`
	UpdatedOn    string         `json:"updated_on"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

func (c *Client) Version(id int) (*Version, error) {
	req, err := c.NewRequest("GET", "/versions/"+strconv.Itoa(id)+".json", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}

	decoder := json.NewDecoder(res.Body)
	var r versionResult
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
	return &r.Version, nil
}

func (c *Client) Versions(projectId int) ([]Version, error) {
	req, err := c.NewRequest("GET", "/projects/"+strconv.Itoa(projectId)+"/versions.json?"+c.getPaginationClause(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}

	decoder := json.NewDecoder(res.Body)
	var r versionsResult
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
	return r.Versions, nil
}

func (c *Client) CreateVersion(version Version) (*Version, error) {
	var ir versionRequest
	ir.Version = version
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", "/projects/"+strconv.Itoa(version.Project.Id)+"/versions.json", strings.NewReader(string(s)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}

	decoder := json.NewDecoder(res.Body)
	var r versionRequest
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
	return &r.Version, err
}

func (c *Client) UpdateVersion(version Version) error {
	var ir versionRequest
	ir.Version = version
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := c.NewRequest("PUT", "/versions/"+strconv.Itoa(version.Id)+".json", strings.NewReader(string(s)))
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
		var er errorsResult
		err = json.NewDecoder(res.Body).Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}

func (c *Client) DeleteVersion(id int) error {
	req, err := c.NewRequest("DELETE", "/versions/"+strconv.Itoa(id)+".json", strings.NewReader(""))
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
		var er errorsResult
		err = json.NewDecoder(res.Body).Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}
