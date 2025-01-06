package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type membershipsResult struct {
	Memberships []Membership `json:"memberships"`
}

type membershipResult struct {
	Membership Membership `json:"membership"`
}

type membershipRequest struct {
	Membership Membership `json:"membership"`
}

type Membership struct {
	Id      int      `json:"id"`
	Project IdName   `json:"project"`
	User    IdName   `json:"user"`
	Roles   []IdName `json:"roles"`
	Groups  []IdName `json:"groups"`
}

func (c *Client) Memberships(projectId int) ([]Membership, error) {
	req, err := c.NewRequest("GET", "/projects/"+strconv.Itoa(projectId)+"/memberships.json?"+c.getPaginationClause(), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r membershipsResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		err = errorFromResp(decoder, res.StatusCode)
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.Memberships, nil
}

func (c *Client) Membership(id int) (*Membership, error) {
	req, err := c.NewRequest("GET", "/memberships/"+strconv.Itoa(id)+".json", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r membershipResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		err = errorFromResp(decoder, res.StatusCode)
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.Membership, nil
}

func (c *Client) CreateMembership(membership Membership) (*Membership, error) {
	var ir membershipRequest
	ir.Membership = membership
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("POST", "/memberships.json", strings.NewReader(string(s)))
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
	var r membershipRequest
	if res.StatusCode != 201 {
		err = errorFromResp(decoder, res.StatusCode)
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.Membership, nil
}

func (c *Client) UpdateMembership(membership Membership) error {
	var ir membershipRequest
	ir.Membership = membership
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := c.NewRequest("PUT", "/memberships/"+strconv.Itoa(membership.Id)+".json", strings.NewReader(string(s)))
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
		err = errorFromResp(decoder, res.StatusCode)
	}
	if err != nil {
		return err
	}
	return err
}

func (c *Client) DeleteMembership(id int) error {
	req, err := c.NewRequest("DELETE", "/memberships/"+strconv.Itoa(id)+".json", strings.NewReader(""))
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
		err = errorFromResp(decoder, res.StatusCode)
	}
	return err
}
