package mock

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

/*
  Usage:

  resp := &http.Response{
    	// specify your desired response here
    }
    client := NewMockHTTPClient(resp)
    oa := &OAuth{
    	// initialize your other fields
    	HTTPClient: client,
    }

*/

type MockHTTPClient struct {
	response *http.Response
}

func NewMockHTTPClient(resp *http.Response) *MockHTTPClient {
	return &MockHTTPClient{response: resp}
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, nil
}

type FakeService func(*http.Request) (*http.Response, error)

func (s FakeService) RoundTrip(req *http.Request) (*http.Response, error) {
	return s(req)
}

// take json string and return *http.Client which returns the string
func HTTPClientJSONString(body string) *http.Client {
	return &http.Client{
		Transport: FakeService(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(string(body))),
			}, nil
		}),
	}
}

// take map[string]interface{} and return *http.Client
func HTTPClientJSONBody(body map[string]interface{}) *http.Client {
	jsonString, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return &http.Client{
		Transport: FakeService(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(string(jsonString))),
			}, nil
		}),
	}
}
