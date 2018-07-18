package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type store interface {
	addChecker(*Checker) error
	removeChecker(id string) error
	findChecker(id string) (*Checker, error)
	allCheckers() []*Checker
	checksSince(ch *Checker, interval *time.Duration, since *time.Time) ([]*Check, error)
}

const (
	CheckerPrefix    = "ch-"
	UptimeDataPrefix = "up-"
)

type fileStore struct {
	path     string
	checkers []*Checker
	mu       *sync.Mutex
}

var (
	ErrPathDoesNotExist = errors.New("Path does not exist")
)

// init and sync with this method
func newFileStore(path string) (*fileStore, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, ErrPathDoesNotExist
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		log.Println(f.Name())
	}

	return nil, nil
}
