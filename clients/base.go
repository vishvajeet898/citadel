package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/utils"
)

type ApiClient struct {
	BaseURL    string
	HttpClient *http.Client
	Ctx        context.Context
}

func NewClient() *ApiClient {
	return &ApiClient{
		HttpClient: &http.Client{
			Timeout: constants.DefaultHTTPClientTimeout,
		},
	}
}

func buildURLWithParams(requestURL string, data map[string]interface{}) (string, error) {
	URL, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}

	parameters := url.Values{}
	for k, v := range data {
		switch val := v.(type) {
		case []int64:
			for _, value := range val {
				parameters.Add(k, fmt.Sprintf("%v", value))
			}
		case []string:
			for _, value := range val {
				parameters.Add(k, fmt.Sprintf("%v", value))
			}
		default:
			parameters.Add(k, fmt.Sprintf("%v", v))
		}
	}

	URL.RawQuery = parameters.Encode()

	return URL.String(), nil
}

func (c *ApiClient) addRequestHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func (c *ApiClient) doRequestWithRetry(req *http.Request, out *interface{}, maxRetries int, baseDelay time.Duration) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = c.doRequest(req, out)
		if err == nil {
			return nil
		}

		if i < maxRetries-1 {
			// Calculate exponential backoff duration with jitter
			delay := baseDelay * time.Duration(math.Pow(2, float64(i)))
			jitter := time.Duration(rand.Int63n(int64(delay) / 2))
			delay += jitter

			time.Sleep(delay)
		}
	}

	return err
}

func (c *ApiClient) doRequest(req *http.Request, out *interface{}) error {
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusOK &&
		response.StatusCode < http.StatusMultipleChoices {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(body, &out); err != nil {
			return err
		}
		return nil
	} else {
		responseData, err1 := io.ReadAll(response.Body)
		if err1 != nil {
			utils.AddLog(c.Ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{"request": req}, err1)
			return errors.New(constants.HTTP_REQUEST_FAILED + ", status code: " + response.Status)
		}
		utils.AddLog(c.Ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), map[string]interface{}{"request": req, "response": string(responseData)}, nil)
		return errors.New(constants.HTTP_REQUEST_FAILED + ", status code: " + response.Status)
	}
}

func (c *ApiClient) Get(
	ctx context.Context,
	out *interface{},
	path string,
	queryParams map[string]interface{},
	data interface{},
	headers map[string]string,
	maxRetries int,
	baseDelay time.Duration,
) error {
	var err error
	c.Ctx = ctx
	URL := fmt.Sprintf("%s%s", c.BaseURL, path)
	if URL, err = buildURLWithParams(URL, queryParams); err != nil {
		return err
	}

	var bodyBytes []byte

	if strings.Contains(headers["Content-Type"], "text/plain") {
		bodyBytes = []byte(data.(string))
	} else {
		bodyBytes, _ = json.Marshal(data)
	}
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		URL,
		bytes.NewBuffer(bodyBytes),
	)

	c.addRequestHeaders(req, headers)
	return c.doRequestWithRetry(req, out, maxRetries, baseDelay)
}

func (c *ApiClient) Post(
	ctx context.Context,
	out *interface{},
	path string,
	queryParams map[string]interface{},
	data interface{},
	headers map[string]string,
	maxRetries int,
	baseDelay time.Duration,
) error {
	var err error
	c.Ctx = ctx
	URL := fmt.Sprintf("%s%s", c.BaseURL, path)
	if URL, err = buildURLWithParams(URL, queryParams); err != nil {
		return err
	}
	jsonValue, _ := json.Marshal(data)
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		URL,
		bytes.NewBuffer(jsonValue),
	)
	c.addRequestHeaders(req, headers)
	return c.doRequestWithRetry(req, out, maxRetries, baseDelay)
}

func (c *ApiClient) Put(
	ctx context.Context,
	out *interface{},
	path string,
	queryParams map[string]interface{},
	data interface{},
	headers map[string]string,
	maxRetries int,
	baseDelay time.Duration,
) error {
	var err error
	c.Ctx = ctx
	URL := fmt.Sprintf("%s%s", c.BaseURL, path)
	if URL, err = buildURLWithParams(URL, queryParams); err != nil {
		return err
	}
	jsonValue, _ := json.Marshal(data)
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		URL,
		bytes.NewBuffer(jsonValue),
	)

	c.addRequestHeaders(req, headers)
	return c.doRequestWithRetry(req, out, maxRetries, baseDelay)
}
