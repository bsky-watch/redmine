package redmine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type CustomFieldDefinition struct {
	Id             int                        `json:"id"`
	Name           string                     `json:"name"`
	Description    *string                    `json:"description"`
	CustomizedType string                     `json:"customized_type"`
	FieldFormat    string                     `json:"field_format"`
	Regexp         string                     `json:"regexp"`
	MinLength      *int                       `json:"min_length"`
	MaxLength      *int                       `json:"max_length"`
	IsRequired     bool                       `json:"is_required"`
	IsFilter       bool                       `json:"is_filter"`
	Searchable     bool                       `json:"searchable"`
	Multiple       bool                       `json:"multiple"`
	DefaultValue   interface{}                `json:"default_value"`
	Visible        bool                       `json:"visible"`
	PossibleValues []CustomFieldPossibleValue `json:"possible_values"`
	Trackers       []IdName                   `json:"trackers"`
	Roles          []IdName                   `json:"roles"`
	EditTagStyle   *string                    `json:"edit_tag_style,omitempty"`
}

// Redmine (as of 6.0.2) has very weird handling of possible_values:
// it is returned as an array of JSON object with "label" and "value" fields,
// but parsed as an array of strings.

type CustomFieldPossibleValue struct {
	Value string
}

func (v *CustomFieldPossibleValue) UnmarshalJSON(b []byte) error {
	var input struct {
		Label string `json:"label"`
	}
	if err := json.Unmarshal(b, &input); err != nil {
		return err
	}
	v.Value = input.Label
	return nil
}

func (v *CustomFieldPossibleValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

type customFieldRequest struct {
	CustomField CustomFieldDefinition `json:"custom_field"`
}

type customFieldsResult struct {
	CustomFields []CustomFieldDefinition `json:"custom_fields"`
}

// CustomFields consulta los campos personalizados
func (c *Client) CustomFields() ([]CustomFieldDefinition, error) {
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

func (c *Client) UpdateCustomField(cf CustomFieldDefinition) error {
	b, err := json.Marshal(&customFieldRequest{CustomField: cf})
	if err != nil {
		return err
	}

	req, err := c.NewRequest("PUT", fmt.Sprintf("/custom_fields/%d.json", cf.Id), bytes.NewReader(b))
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
	if res.StatusCode != 204 {
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
