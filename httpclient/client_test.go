package httpclient

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/flute/flute"
)

func TestError_Error(t *testing.T) {
	data := []struct {
		title string
		exp   string
		err   *Error
	}{
		{
			title: "normal",
			err: &Error{
				statusCode: 500,
				bodyByte:   []byte(`{"error": "Internal Server Error"}`),
				err:        errors.New("status code >= 300"),
			},
			exp: `status code: 500, {"error": "Internal Server Error"}: status code >= 300`,
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			require.Equal(t, d.exp, d.err.Error())
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	exp := errors.New("foo")
	err := &Error{
		statusCode: 500,
		bodyByte:   []byte(`{"error": "Internal Server Error"}`),
		err:        exp,
	}
	if e := errors.Unwrap(err); e != exp {
		t.Errorf(`errors.Unwrap(err) = %v; want %v`, e, exp)
	}
}

func TestClient_Call(t *testing.T) {
	client := New("http://example.com/api")
	token := "xxx"
	client.SetRequest = func(req *http.Request) error {
		req.Header.Add("Authorization", "token "+token)
		return nil
	}
	client.HTTPClient = &http.Client{
		Transport: &flute.Transport{
			T: t,
			Services: []flute.Service{
				{
					Endpoint: "http://example.com",
					Routes: []flute.Route{
						{
							Name: "create a user",
							Matcher: &flute.Matcher{
								Method: "POST",
								Path:   "/api/users",
							},
							Tester: &flute.Tester{
								BodyJSONString: `{
										  "name": "foo",
										  "email": "foo@example.com"
										}`,
								Header: http.Header{
									"Authorization": []string{"token " + token},
								},
							},
							Response: &flute.Response{
								Base: http.Response{
									StatusCode: 201,
								},
								BodyString: `{
										  "id": 10,
										  "name": "foo",
										  "email": "foo@example.com"
										}`,
							},
						},
						{
							Name: "get a group",
							Matcher: &flute.Matcher{
								Method: "GET",
								Path:   "/api/groups/foo",
							},
							Response: &flute.Response{
								Base: http.Response{
									StatusCode: 404,
								},
								BodyString: `{
										  "error": "group foo isn't found"
										}`,
							},
						},
					},
				},
			},
		},
	}
	ctx := context.Background()
	data := []struct {
		title            string
		params           *CallParams
		exp              interface{}
		expErrorResponse interface{}
		isErr            bool
	}{
		{
			title: "request body is struct",
			params: &CallParams{
				Method: "POST",
				Path:   "/users",
				RequestBody: struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				}{
					Name:  "foo",
					Email: "foo@example.com",
				},
				ResponseBody: &map[string]interface{}{},
			},
			exp: &map[string]interface{}{
				"id":    10.0,
				"name":  "foo",
				"email": "foo@example.com",
			},
		},
		{
			title: "request body is string",
			params: &CallParams{
				Method:       "POST",
				Path:         "/users",
				RequestBody:  `{"name": "foo", "email": "foo@example.com"}`,
				ResponseBody: &map[string]interface{}{},
			},
			exp: &map[string]interface{}{
				"id":    10.0,
				"name":  "foo",
				"email": "foo@example.com",
			},
		},
		{
			title: "request body is []byte",
			params: &CallParams{
				Method:       "POST",
				Path:         "/users",
				RequestBody:  []byte(`{"name": "foo", "email": "foo@example.com"}`),
				ResponseBody: &map[string]interface{}{},
			},
			exp: &map[string]interface{}{
				"id":    10.0,
				"name":  "foo",
				"email": "foo@example.com",
			},
		},
		{
			title: "error response",
			params: &CallParams{
				Method:            "GET",
				Path:              "/groups/foo",
				ResponseErrorBody: &map[string]interface{}{},
			},
			expErrorResponse: &map[string]interface{}{
				"error": "group foo isn't found",
			},
			isErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			err := client.Call(ctx, d.params)
			if d.isErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, d.exp, d.params.ResponseBody)
			}
			require.Equal(t, d.expErrorResponse, d.params.ResponseErrorBody)
			if err != nil {
				var e *Error
				if errors.As(err, &e) {
					require.Equal(t, d.expErrorResponse, e.Body())
				}
			}
		})
	}
}
