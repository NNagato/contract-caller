package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// RestClient ...
type RestClient struct {
	c *http.Client
}

// NewRestClient ...
func NewRestClient(c *http.Client) *RestClient {
	return &RestClient{
		c: c,
	}
}

// DoReq ...
func (rc *RestClient) DoReq(url, method string, data, result interface{}) error {
	var (
		httpMethod = strings.ToUpper(method)
		body       io.Reader
	)
	if httpMethod != http.MethodGet && data != nil {
		dataBody, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(dataBody)
	}
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	switch httpMethod {
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		req.Header.Add("Content-Type", "application/json")
	case http.MethodGet:
	default:
		return errors.Errorf("invalid method %s", httpMethod)
	}
	rsp, err := rc.c.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do req")
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read body")
	}
	if err := rsp.Body.Close(); err != nil {
		return errors.Wrap(err, "failed to close body")
	}
	if rsp.StatusCode != http.StatusOK {
		var rspData interface{}
		if err := json.Unmarshal(rspBody, &rspData); err != nil {
			return errors.Wrap(err, "cannot unmarshal response data")
		}
		return errors.Errorf("receive unexpected code, actual code: %d, data: %+v", rsp.StatusCode, rspData)
	}
	if result != nil {
		return json.Unmarshal(rspBody, result)
	}
	return nil
}
