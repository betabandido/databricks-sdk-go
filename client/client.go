package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/betabandido/databricks-sdk-go/models"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http    *http.Client
	baseUrl *url.URL
	header  http.Header
}

func NewClient(domain string, token string) (*Client, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("https://%s/api/2.0/", domain))
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseUrl: baseUrl,
		header:  header,
	}

	return &client, nil
}

func (c *Client) Query(method string, path string, data interface{}) ([]byte, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	queryUrl := c.baseUrl.ResolveReference(u)

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, queryUrl.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header = c.header

	glog.Infof("HTTP request: %v", request)

	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}

	glog.Infof("HTTP response: %v", response)

	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		errorResponse := models.ErrorResponse{}
		err = json.Unmarshal(responseBytes, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, Error{ErrorResponse: errorResponse}
	}

	return responseBytes, nil
}

type Error struct {
	ErrorResponse models.ErrorResponse
}

func (e Error) Error() string {
	return e.ErrorResponse.Message
}

func (e Error) Code() string {
	return e.ErrorResponse.ErrorCode
}
