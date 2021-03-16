package client

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	Uri          string
	AuthLoginPwd string
	UserAgent    string
	Interval     int
	Timeout      int
	Limit        int
}

func NewApiClient(conf Config) *ApiClient {
	return &ApiClient{
		endpoint: conf.Uri,

		headers: map[string]string{
			`Authorization`: `Basic ` + base64.StdEncoding.EncodeToString([]byte(conf.AuthLoginPwd)),
			`User-Agent`:    conf.UserAgent,
		},

		client: &http.Client{
			Timeout: time.Second * time.Duration(conf.Timeout),
		},

		limit: conf.Limit,
	}
}

type ApiClient struct {
	endpoint string

	client *http.Client

	headers map[string]string

	offsetSize *http.Client

	limit int
}

func (c *ApiClient) GetItems(offset int) (Response, error) {
	var response Response
	request, err := http.NewRequest(http.MethodGet, c.endpoint, nil)
	if err != nil {
		return response, err
	}

	for header, value := range c.headers {
		request.Header.Add(header, value)
	}
	q := request.URL.Query()
	q.Set(`p_offset`, strconv.Itoa(offset))
	q.Set(`p_limit`, strconv.Itoa(c.limit))
	request.URL.RawQuery = q.Encode()
	response.EndPoint = request.URL.String()

	resp, err := c.client.Do(request)
	if err != nil {
		return response, err
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	if len(bodyBytes) == 0 {
		return response, nil
	}

	err = json.Unmarshal(bodyBytes, &response)
	return response, err
}
