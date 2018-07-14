package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

type endPoint struct {
	path           string
	method         string
	isHTTPS        bool
	allowRedirects bool
}

const (
	HTTPPrefix  = "http://"
	HTTPSPrefix = "https://"
	MinInterval = 2 * time.Second
	MaxTimeout  = 5 * time.Minute
)

var (
	AvailableMethods = [4]string{"GET", "POST", "PUT", "DELETE"}
)

func newEndPoint(path, method string, redirects bool) (*endPoint, error) {
	isHTTP := strings.Contains(path, HTTPPrefix)
	isHTTPs := strings.Contains(path, HTTPSPrefix)

	if !isHTTP && isHTTPs {
		return nil, ErrEndPointPathFormat
	}
	for _, m := range AvailableMethods {
		if m == method {
			return &endPoint{
				path,
				method,
				isHTTPs,
				redirects,
			}, nil
		}
	}
	return nil, ErrEndPointMethod
}

const (
	endPointTemplate = `` +
		`Path			[%s]` + "\n" +
		`Method			[%s]` + "\n" +
		`Allow Redirects		[%t]` + "\n"
)

func (ep *endPoint) String() string {
	return fmt.Sprintf(endPointTemplate, ep.path, ep.method, ep.allowRedirects)
}

type checker struct {
	id       string
	ep       *endPoint
	active   bool
	timeout  time.Duration
	interval time.Duration
	ticker   *time.Ticker
	stopped  chan bool
}

func newChecker(ep *endPoint, timeout, interval time.Duration) (*checker, error) {
	if ep == nil {
		return nil, ErrEndPointCantBeNil
	}
	if interval < MinInterval {
		return nil, ErrIntervalCantBeLower
	}
	if timeout > MaxTimeout {
		return nil, ErrTimeoutCantBeHigher
	}
	if timeout > interval {
		return nil, ErrTimeoutShouldBeLesser
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &checker{
		id.String(),
		ep,
		false,
		timeout,
		interval,
		&time.Ticker{},
		make(chan bool),
	}, nil
}

var (
	ErrEndPointPathFormat    = errors.New("End point is not valid")
	ErrEndPointMethod        = errors.New("End point method is not valid")
	ErrEndPointCantBeNil     = errors.New("End Point can't be nil")
	ErrIntervalCantBeLower   = errors.New("Interval is below threshold")
	ErrTimeoutCantBeHigher   = errors.New("Timeout is above threshold")
	ErrTimeoutShouldBeLesser = errors.New("Timeout can't be higher than interval")
	ErrCheckerAlreadyStarted = errors.New("Checker is already started")
	ErrCheckerAlreadyStopped = errors.New("Checker is already stopped")
)

func (ch *checker) start(fn func(ep *endPoint) error) error {
	if ch.active {
		return ErrCheckerAlreadyStarted
	}
	ch.active = true
	ch.ticker = time.NewTicker(ch.interval)
	log.Printf("checker starting for \n%s", ch.ep)

	go func() {
		defer log.Printf("checker stopped for \n%s", ch.ep)

		for {
			select {
			case <-ch.ticker.C:
				if err := fn(ch.ep); err != nil {
					log.Printf("checker error \n%s", err)
				}
			case <-ch.stopped:
				return
			}
		}
	}()
	return nil
}

func (ch *checker) stop() error {
	if !ch.active {
		return ErrCheckerAlreadyStopped
	}
	ch.active = false
	ch.ticker.Stop()
	ch.stopped <- true

	return nil
}
