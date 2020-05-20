package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	// TODO: replace with gitlab.com and move it to config
	defaultBaseURL          = "https://gitlab.com/"
	defaultVersionedAPIPath = "api/v4/"
	defaultUserAgent        = "go-lazylab"
)

type Client struct {
	client *http.Client

	baseURL *url.URL
	token   string

	UserAgent string

	MergeRequests *MergeRequestsService
}

type ClientConfig struct {
	User   string `json:"user"`
	UserID int    `json:"user_id"`
	Server string `json:"server"`
	Token  string `json:"token"`
}

func LoadConfig(cfg io.ReadCloser) (ClientConfig, error) {
	defer cfg.Close()
	parsedCfg := ClientConfig{}
	if err := json.NewDecoder(cfg).Decode(&parsedCfg); err != nil {
		return parsedCfg, err
	}
	return parsedCfg, nil
}

type RequestOptionFunc func(*http.Request) error

type ClientOptionsFunc func(client *Client) error

func NewClient(config ClientConfig, options ...ClientOptionsFunc) (*Client, error) {
	c := &Client{
		token:     config.Token,
		UserAgent: defaultUserAgent,
	}

	base, err := url.Parse(config.Server)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	if !strings.HasSuffix(base.Path, "/"+defaultVersionedAPIPath) {
		base.Path += defaultVersionedAPIPath
	}

	c.client = &http.Client{
		Transport: &http.Transport{},
	}
	c.baseURL = base

	if options != nil {
		for _, opt := range options {
			if opt == nil {
				continue
			}
			if err := opt(c); err != nil {
				return nil, err
			}
		}
	}

	c.MergeRequests = &MergeRequestsService{c}

	return c, nil
}

func (c *Client) NewRequest(method, path string, opt interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	switch {
	case method == "POST" || method == "PUT":
		return nil, fmt.Errorf("method %s currently is not supported", method)
	case opt != nil:
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

type Response struct {
	*http.Response

	TotalItems   int
	TotalPages   int
	ItemsPerPage int
	CurrentPage  int
	NextPage     int
	PreviousPage int
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}
	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if w, ok := v.(io.Writer); ok {
		_, err = io.Copy(w, resp.Body)
	} else {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return response, err
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

func CheckResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data
		errorResponse.Message = string(data)
	}

	return errorResponse
}
