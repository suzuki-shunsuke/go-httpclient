package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	Endpoint   string
	HTTPClient *http.Client
	SetRequest func(req *http.Request) error
}

func New(endpoint string) *Client {
	return &Client{
		Endpoint:   endpoint,
		HTTPClient: http.DefaultClient,
	}
}

type Error struct {
	statusCode int
	bodyByte   []byte
	body       interface{}
	err        error
}

func (e *Error) StatusCode() int {
	return e.statusCode
}

func (e *Error) BodyByte() []byte {
	return e.bodyByte
}

func (e *Error) Body() interface{} {
	return e.body
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	a := ""
	if e.err != nil {
		a = e.err.Error()
	}
	return "status code: " + strconv.Itoa(e.statusCode) + ", " + string(e.bodyByte) + ": " + a
}

func (e *Error) Unwrap() error {
	return e.err
}

type CallParams struct {
	Method            string
	Path              string
	RequestBody       interface{}
	ResponseBody      interface{}
	ResponseErrorBody interface{}
}

func (client *Client) Call(ctx context.Context, api *CallParams) error {
	if client.Endpoint == "" {
		return errors.New("endpoint is required")
	}
	if api.Method == "" {
		return errors.New("method is required")
	}

	var body io.Reader
	if api.RequestBody != nil {
		switch b := api.RequestBody.(type) {
		case string:
			body = strings.NewReader(b)
		case []byte:
			body = bytes.NewBuffer(b)
		default:
			buf := &bytes.Buffer{}
			if err := json.NewEncoder(buf).Encode(api.RequestBody); err != nil {
				return fmt.Errorf("failed to parse the request body as JSON: %w", err)
			}
			body = buf
		}
	}

	req, err := http.NewRequestWithContext(ctx, api.Method, client.Endpoint+api.Path, body)
	if err != nil {
		return fmt.Errorf("failed to create a request: %w", err)
	}
	if client.SetRequest != nil {
		if err := client.SetRequest(req); err != nil {
			return fmt.Errorf("failed to set a request: %w", err)
		}
	}

	httpClient := client.HTTPClient
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send a request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return &Error{
				statusCode: res.StatusCode,
				err:        fmt.Errorf("status code >= 300: failed to read a response body: %w", err),
			}
		}
		if api.ResponseErrorBody != nil {
			if err := json.Unmarshal(body, api.ResponseErrorBody); err != nil {
				return &Error{
					statusCode: res.StatusCode,
					bodyByte:   body,
					err:        fmt.Errorf("status code >= 300: failed to parse an error response body as JSON: %w", err),
				}
			}
			return &Error{
				statusCode: res.StatusCode,
				bodyByte:   body,
				body:       api.ResponseErrorBody,
				err:        errors.New("status code >= 300"),
			}
		}
		return &Error{
			statusCode: res.StatusCode,
			bodyByte:   body,
			err:        errors.New("status code >= 300"),
		}
	}

	if api.ResponseBody != nil {
		if err := json.NewDecoder(res.Body).Decode(api.ResponseBody); err != nil {
			return fmt.Errorf("failed to read a response body as JSON: %w", err)
		}
	}
	return nil
}
