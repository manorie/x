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
	check := Call("https://www.example.com", "PATCH", nil, 2*time.Second)
	assert.NotNil(t, check)
	assert.NotNil(t, check.Err)
	assert.Equal(t, ErrCallUnsupportedMethod, check.Err)

	check = Call("dumy", MethodGET, nil, 4*time.Second)
	assert.NotNil(t, check)
	assert.NotNil(t, check.Err)
	assert.Equal(t, "parse dumy: invalid URI for request", check.Err.Error())

	check = Call("", MethodPOST, nil, 3*time.Second)
	assert.NotNil(t, check)
	assert.NotNil(t, check.Err)
	assert.Equal(t, "parse : empty url", check.Err.Error())
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

	timeoutCall := Call(ts0.URL, MethodGET, nil, time.Millisecond*400)
	assert.NotNil(t, timeoutCall.Err)

	ts1 := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts1.Close()

	successCall := Call(ts1.URL, MethodGET, nil, time.Second*1)
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

	notFoundCall := Call(ts2.URL, MethodGET, nil, time.Second*1)
	assert.Nil(t, notFoundCall.Err)
	assert.Equal(t, 404, notFoundCall.StatusCode)
	assert.Equal(t, len("Page not found"), notFoundCall.Bytes)
}
