package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type journalRequest struct {
	Journal *Journal `json:"journal"`
}

func (c *Client) UpdateJournal(journal *Journal) error {
	var jr journalRequest
	jr.Journal = journal
	s, err := json.Marshal(jr)
	if err != nil {
		return err
	}
	req, err := c.NewRequest("PUT", "/journals/"+strconv.Itoa(journal.Id)+".json", strings.NewReader(string(s)))
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
	if res.StatusCode/100 != 2 {
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
