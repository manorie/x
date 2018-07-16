package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

func call(ep *endPoint, timeout time.Duration) *check {
	url, err := url.Parse(ep.path)
	if err != nil {
		return &check{
			err: err,
		}
	}

	req, err := http.NewRequest(ep.method, url.String(), nil)
	if err != nil {
		return &check{
			err: err,
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
		return &check{
			err: err,
		}
	}

	bodyRead, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &check{
			err:        err,
			statusCode: resp.StatusCode,
		}
	}
	t1 = time.Now()

	return &check{
		nil,
		resp.StatusCode,
		len(bodyRead),
		t1.Sub(t0).Seconds(),
	}
}
