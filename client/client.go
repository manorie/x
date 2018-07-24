package client

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

// Check ..
type Check struct {
	Err         error   `json:"error"`
	StatusCode  int     `json:"statusCode"`
	Bytes       int     `json:"bytes"`
	TimeElapsed float64 `json:"timeElapsed"`
}

const (
	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodDELETE = "DELETE"
)

var (
	AvailableMethods = map[string]bool{
		MethodGET:    true,
		MethodPOST:   true,
		MethodPUT:    true,
		MethodDELETE: true,
	}
)

var (
	ErrCallUnsupportedMethod = errors.New("unsupported request method")
)

type Client struct {
	request *http.Request
	timeout time.Duration
}

func NewClient(path, method string, body io.Reader, timeout time.Duration) (*Client, error) {
	if _, present := AvailableMethods[method]; !present {
		return nil, ErrCallUnsupportedMethod
	}

	url, err := url.ParseRequestURI(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	return &Client{
		request: req,
		timeout: timeout,
	}, nil
}

func (cl *Client) Call() *Check {
	var t0, t1 time.Time
	trace := &httptrace.ClientTrace{
		ConnectStart: func(_, _ string) {
			t0 = time.Now()
		},
	}

	cl.request = cl.request.WithContext(httptrace.WithClientTrace(
		context.Background(), trace,
	))

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    cl.timeout,
			DisableCompression: true,
		},
		Timeout: cl.timeout,
	}

	resp, err := client.Do(cl.request)
	if err != nil {
		return &Check{
			Err: err,
		}
	}

	bodyRead, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &Check{
			Err:        err,
			StatusCode: resp.StatusCode,
		}
	}
	t1 = time.Now()

	return &Check{
		Err:         nil,
		StatusCode:  resp.StatusCode,
		Bytes:       len(bodyRead),
		TimeElapsed: t1.Sub(t0).Seconds(),
	}
}
