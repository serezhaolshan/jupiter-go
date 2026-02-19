package jupiter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"golang.org/x/time/rate"
)

const DefaultURL = "https://api.jup.ag"

const ApiKeyHeader = "x-api-key"

const RateLimitMilliseconds = 120

type Client struct {
	ApiUrl  string
	ApiKey  string
	Limiter *rate.Limiter
	c       *http.Client
}

func NewClient(url, key string) *Client {
	return &Client{
		ApiUrl:  url,
		ApiKey:  key,
		Limiter: rate.NewLimiter(rate.Every(RateLimitMilliseconds*time.Millisecond), 1),
		c:       http.DefaultClient,
	}
}

func (c *Client) Url(endpoint string) string {
	return fmt.Sprintf("%s%s", c.ApiUrl, endpoint)
}

func (c *Client) doCall(ctx context.Context, req *Request, response any) (*http.Response, error) {
	err := c.Limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	if reflect.TypeOf(response).Kind() != reflect.Pointer {
		return nil, fmt.Errorf("response struct is not a pointer")
	}
	httpRequest, err := req.NewHttpRequest(ctx, c.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("api call %v() on %v: %v", req.Method, req.Endpoint, err.Error())
	}
	httpResponse, err := c.c.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("api call %v() on %v: %v", req.Method, httpRequest.URL.String(), err.Error())
	}

	bodyBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"call %v() on %v status code: %v. could not decode body to response: %v",
			req.Method,
			httpRequest.URL.String(),
			httpResponse.StatusCode,
			err.Error())
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"call %v() on %v status code: %v.raw body %v",
			req.Method,
			httpRequest.URL.String(),
			httpResponse.StatusCode,
			string(bodyBytes))
	}

	err = json.Unmarshal(bodyBytes, response)

	if err != nil {
		return nil, fmt.Errorf(
			"call %v() on %v status code: %v. could not decode body to response model: %v",
			req.Method,
			httpRequest.URL.String(),
			httpResponse.StatusCode,
			err.Error())
	}
	if response == nil {
		return nil, fmt.Errorf("call %v() on %v status code: %v. response missing",
			req.Method,
			httpRequest.URL.String(),
			httpResponse.StatusCode)
	}

	return httpResponse, nil
}

type Request struct {
	Endpoint    string
	Method      string
	QueryParams url.Values
	Body        io.Reader
}

func NewRequest(endpoint string, queryParams url.Values, methods ...string) *Request {
	method := http.MethodGet
	if len(methods) != 0 {
		method = methods[0]
	}
	return &Request{
		Endpoint:    endpoint,
		Method:      method,
		QueryParams: queryParams,
	}
}

func NewPostRequest(endpoint string, body any) (*Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}
	return &Request{
		Endpoint: endpoint,
		Method:   http.MethodPost,
		Body:     bytes.NewReader(jsonBody),
	}, nil
}

func (r *Request) NewHttpRequest(ctx context.Context, apiKey string) (*http.Request, error) {
	fullURL := r.Endpoint
	if len(r.QueryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s", r.Endpoint, r.QueryParams.Encode())
	}

	request, err := http.NewRequest(r.Method, fullURL, r.Body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		request.Header.Set(ApiKeyHeader, apiKey)
	}
	request = request.WithContext(ctx)

	return request, nil
}