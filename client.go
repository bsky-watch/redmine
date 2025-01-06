package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type Client struct {
	endpoint string
	apikey   string
	*http.Client
	Limit  int
	Offset int
}

var DefaultLimit int = -1  // "-1" means "No setting"
var DefaultOffset int = -1 //"-1" means "No setting"

func NewClient(endpoint, apikey string) *Client {
	return &Client{endpoint, apikey, http.DefaultClient, DefaultLimit, DefaultOffset}
}

// URLWithFilter return string url by concat endpoint, path and filter
// err != nil when endpoin can not parse
func (c *Client) URLWithFilter(path string, f Filter) (string, error) {
	var fullURL *url.URL
	fullURL, err := url.Parse(c.endpoint)
	if err != nil {
		return "", err
	}
	fullURL.Path += path
	if c.Limit > -1 {
		f.AddPair("limit", strconv.Itoa(c.Limit))
	}
	if c.Offset > -1 {
		f.AddPair("offset", strconv.Itoa(c.Offset))
	}
	fullURL.RawQuery = f.ToURLParams()
	return fullURL.String(), nil
}

func (c *Client) getPaginationClause() string {
	clauses := []string{}
	if c.Limit > -1 {
		clauses = append(clauses, fmt.Sprintf("limit=%v", c.Limit))
	}
	if c.Offset > -1 {
		clauses = append(clauses, fmt.Sprintf("offset=%v", c.Offset))
	}
	return strings.Join(clauses, "&")
}

type errorsResult struct {
	Errors []string `json:"errors"`
}

func errorFromResp(bodyDecoder *json.Decoder, httpStatus int) (err error) {
	var er errorsResult
	err = bodyDecoder.Decode(&er)
	if err == nil { /* error from redmine */
		return errors.New(strings.Join(er.Errors, "\n"))
	}

	if err == io.EOF { /* empty body */
		err = errors.New(http.StatusText(httpStatus))
	}

	return
}

type IdName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Id struct {
	Id int `json:"id"`
}

func (c *Client) NewRequest(method string, urlPath string, body io.Reader) (*http.Request, error) {

	// Hack to avoid changing how URLWithFilter works.
	if !strings.HasPrefix("http://", urlPath) && !strings.HasPrefix("https://", urlPath) {
		urlPath = path.Join(c.endpoint, urlPath)
	}

	r, err := http.NewRequest(method, urlPath, body)
	if err != nil {
		return nil, err
	}
	if c.apikey != "" {
		r.Header.Set("X-Redmine-API-Key", c.apikey)
	}
	return r, nil
}
