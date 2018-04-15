package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// composite type
type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

// for twilio, it does not send json body...
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	httpClient := initClient(c)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func initClient(c *Client) *http.Client {
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	return c.httpClient
}
