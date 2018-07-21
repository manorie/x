package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

type EndPoint struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type Check struct {
	Err         error   `json:"-"`
	StatusCode  int     `json:"statusCode"`
	Bytes       int     `json:"bytes"`
	TimeElapsed float64 `json:"timeElapsed"`
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

func newEndPoint(path, method string, redirects bool) (*EndPoint, error) {
	isHTTP := strings.Contains(path, HTTPPrefix)
	isHTTPS := strings.Contains(path, HTTPSPrefix)

	if !isHTTP && !isHTTPS {
		return nil, ErrEndPointPathFormat
	}
	for _, m := range AvailableMethods {
		if m == method {
			return &EndPoint{
				path,
				method,
			}, nil
		}
	}
	return nil, ErrEndPointMethod
}

const (
	endPointTemplate = `` +
		`Path			[%s]` + "\n" +
		`Method			[%s]` + "\n"
)

func (ep *EndPoint) String() string {
	return fmt.Sprintf(endPointTemplate, ep.Path, ep.Method)
}

type Checker struct {
	ID       string        `json:"id"`
	EP       *EndPoint     `json:"endPoint"`
	Active   bool          `json:"active"`
	Timeout  time.Duration `json:"timeout,timeunit:s"`
	Interval time.Duration `json:"inteval,timeunit:s"`
	ticker   *time.Ticker
	stopped  chan bool
	mu       *sync.Mutex
}

func newChecker(ep *EndPoint, timeout, interval time.Duration) (*Checker, error) {
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

	return &Checker{
		id.String(),
		ep,
		false,
		timeout,
		interval,
		&time.Ticker{},
		make(chan bool),
		&sync.Mutex{},
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

func (ch *Checker) start(fn func(ep *EndPoint) error) error {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ch.Active {
		return ErrCheckerAlreadyStarted
	}
	ch.Active = true
	ch.ticker = time.NewTicker(ch.Interval)
	log.Printf("Checker starting for \n%s", ch.EP)

	go func() {
		defer log.Printf("Checker stopped for \n%s", ch.EP)

		for {
			select {
			case <-ch.ticker.C:
				if err := fn(ch.EP); err != nil {
					log.Printf("Checker error \n%s", err)
				}
			case <-ch.stopped:
				return
			}
		}
	}()
	return nil
}

func (ch *Checker) stop() error {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if !ch.Active {
		return ErrCheckerAlreadyStopped
	}
	ch.Active = false
	ch.ticker.Stop()
	ch.stopped <- true

	return nil
}
