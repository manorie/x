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

// set request ready with body to post
// enrich header
// -> then hit with command and return check result

// client -> NEW may return error -> have a UUID and stored
// client.hit to perform

func Call(path, method string, body io.Reader, timeout time.Duration) *Check {
	if _, present := AvailableMethods[method]; !present {
		return &Check{
			Err: ErrCallUnsupportedMethod,
		}
	}

	url, err := url.ParseRequestURI(path)
	if err != nil {
		return &Check{
			Err: err,
		}
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return &Check{
			Err: err,
		}
	}

	var t0, t1 time.Time
	trace := &httptrace.ClientTrace{
		ConnectStart: func(_, _ string) {
			t0 = time.Now()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(
		context.Background(), trace,
	))

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    timeout,
			DisableCompression: true,
		},
		Timeout: timeout,
	}

	resp, err := client.Do(req)
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
