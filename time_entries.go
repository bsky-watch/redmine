package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type timeEntriesResult struct {
	TimeEntries []TimeEntry `json:"time_entries"`
}

type timeEntryResult struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

type timeEntryRequest struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

type TimeEntry struct {
	Id           int            `json:"id"`
	Project      IdName         `json:"project"`
	Issue        Id             `json:"issue"`
	User         IdName         `json:"user"`
	Activity     IdName         `json:"activity"`
	Hours        float32        `json:"hours"`
	Comments     string         `json:"comments"`
	SpentOn      string         `json:"spent_on"`
	CreatedOn    string         `json:"created_on"`
	UpdatedOn    string         `json:"updated_on"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

// TimeEntriesWithFilter send query and return parsed result
func (c *Client) TimeEntriesWithFilter(filter Filter) ([]TimeEntry, error) {
	uri, err := c.URLWithFilter("/time_entries.json", filter)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntriesResult
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
	return r.TimeEntries, nil
}

func (c *Client) TimeEntries(projectId int) ([]TimeEntry, error) {
	req, err := c.NewRequest("GET", "/projects/"+strconv.Itoa(projectId)+"/time_entries.json"+c.getPaginationClause(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntriesResult
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
	return r.TimeEntries, nil
}

func (c *Client) TimeEntry(id int) (*TimeEntry, error) {
	req, err := c.NewRequest("GET", "/time_entries/"+strconv.Itoa(id)+".json", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntryResult
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
	return &r.TimeEntry, nil
}

func (c *Client) CreateTimeEntry(timeEntry TimeEntry) (*TimeEntry, error) {
	var ir timeEntryRequest
	ir.TimeEntry = timeEntry
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", "/time_entries.json", strings.NewReader(string(s)))
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
	var r timeEntryResult
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
	return &r.TimeEntry, nil
}

func (c *Client) UpdateTimeEntry(timeEntry TimeEntry) error {
	var ir timeEntryRequest
	ir.TimeEntry = timeEntry
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := c.NewRequest("PUT", "/time_entries/"+strconv.Itoa(timeEntry.Id)+".json", strings.NewReader(string(s)))
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

func (c *Client) DeleteTimeEntry(id int) error {
	req, err := c.NewRequest("DELETE", "/time_entries/"+strconv.Itoa(id)+".json", strings.NewReader(""))
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
