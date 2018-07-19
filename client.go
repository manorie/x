package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

func call(ep *EndPoint, timeout time.Duration) *Check {
	url, err := url.Parse(ep.Path)
	if err != nil {
		return &Check{
			Err: err,
		}
	}

	req, err := http.NewRequest(ep.Method, url.String(), nil)
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
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    timeout,
			DisableCompression: true,
		},
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
		nil,
		resp.StatusCode,
		len(bodyRead),
		t1.Sub(t0).Seconds(),
	}
}

// ss
