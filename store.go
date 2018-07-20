package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type fileStore struct {
	basePath string
	files    map[string]*file
}

type file struct {
	name string
	mu   sync.Mutex
}

func (fs *fileStore) createNewFile(name string) error {
	return nil
}

func (fs *fileStore) readFile(name string, fn func(line string)) error {
	return nil
}

func (fs *fileStore) appendToFile(name, toWrite string) error {
	return nil
}

func (fs *fileStore) deleteFile(name string) error {
	return nil
}

// nees to have basePath ending with "/"
func (fs *fileStore) writeAsJSON(name string, obj interface{}) error {
	json, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fs.basePath+name, json, 0644)
}

func (fs *fileStore) readFromJSON(name string)

// import (
// 	"errors"
// 	"os"
// 	"sync"
// 	"time"
// )

// type store interface {
// 	addChecker(*Checker) error
// 	removeChecker(id string) error
// 	findChecker(id string) (*Checker, error)
// 	allCheckers() []*Checker
// 	checksSince(ch *Checker, interval *time.Duration, since *time.Time) ([]*Check, error)
// }

// const (
// 	CheckerPrefix    = "ch-"
// 	UptimeDataPrefix = "up-"
// )

// type fileStore struct {
// 	path string
// 	mu   *sync.Mutex
// }

// var (
// 	ErrPathDoesNotExist = errors.New("Path does not exist")
// )

// // init and sync with this method
// func newFileStore(path string) (*fileStore, error) {
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		return nil, ErrPathDoesNotExist
// 	}

// 	// files, err := ioutil.ReadDir(path)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// for _, f := range files {
// 	// 	log.Println(f.Name())
// 	// }

// 	return &fileStore{
// 		path:     path,
// 		checkers: make([]*Checker, 0, 1),
// 	}, nil
// }
