package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCallWithWrongArguments(t *testing.T) {
	cl, err := NewClient("https://www.example.com", "PATCH", nil, 2*time.Second)
	assert.NotNil(t, err)
	assert.Nil(t, cl)
	assert.Equal(t, ErrCallUnsupportedMethod, err)

	cl, err = NewClient("dumyURL", "GET", nil, 2*time.Second)
	assert.NotNil(t, err)
	assert.Nil(t, cl)
	assert.Equal(t, "parse dumyURL: invalid URI for request", err.Error())

	cl, err = NewClient("", MethodPOST, nil, 3*time.Second)
	assert.NotNil(t, err)
	assert.Nil(t, cl)
	assert.Equal(t, "parse : empty url", err.Error())
}

func TestCall(t *testing.T) {
	sleepingHandler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 500)
	}

	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r)
	}

	ts0 := httptest.NewServer(http.HandlerFunc(sleepingHandler))
	defer ts0.Close()

	cl, err := NewClient(ts0.URL, MethodGET, nil, time.Millisecond*400)
	assert.Nil(t, err)
	timeoutCall := cl.Call()
	assert.NotNil(t, timeoutCall.Err)

	ts1 := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts1.Close()

	cl, err = NewClient(ts1.URL, MethodGET, nil, time.Second*1)
	assert.Nil(t, err)

	successCall := cl.Call()
	assert.Nil(t, successCall.Err)
	assert.Equal(t, 200, successCall.StatusCode)
	assert.Equal(t, 167, successCall.Bytes)
	assert.True(t, successCall.TimeElapsed < (time.Millisecond*400).Seconds())

	notFoundHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "Page not found")
	}

	ts2 := httptest.NewServer(http.HandlerFunc(notFoundHandler))
	defer ts2.Close()

	cl, err = NewClient(ts2.URL, MethodGET, nil, time.Second*1)
	assert.Nil(t, err)

	notFoundCall := cl.Call()
	assert.Nil(t, notFoundCall.Err)
	assert.Equal(t, 404, notFoundCall.StatusCode)
	assert.Equal(t, len("Page not found"), notFoundCall.Bytes)
}
