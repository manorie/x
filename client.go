package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

func call(ep *endPoint, timeout time.Duration) (*check, error) {
	url, err := url.Parse(ep.path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(ep.method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	var t0, t1 time.Time

	trace := &httptrace.ClientTrace{
		ConnectStart: func(_, _ string) {
			t0 = time.Now()
		},
		ConnectDone: func(_, _ string, errR error) {
			err = errR
			t1 = time.Now()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(resp)
	}

	log.Println(t1.Sub(t0))

	return nil, nil
}
