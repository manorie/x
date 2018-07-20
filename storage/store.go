package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type fileStore struct {
	basePath string
	files    map[string]*file
}

type file struct {
	mu *sync.Mutex
}

// NewFileStore ...
func NewFileStore(base string) (*fileStore, error) {
	if base[len(base)-1] != '/' {
		base += "/"
	}

	files, err := ioutil.ReadDir(base)
	if err != nil {
		return nil, err
	}

	fs := &fileStore{
		basePath: base,
		files:    make(map[string]*file),
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fs.files[f.Name()] = &file{
			mu: &sync.Mutex{},
		}
	}
	return fs, nil
}

func (fs *fileStore) createFile(name string) error {
	newFilePath := fs.basePath + name

	if _, err := os.Stat(newFilePath); err == nil {
		return os.ErrExist
	}

	if _, err := os.Create(newFilePath); err != nil {
		return err
	}

	fs.files[name] = &file{
		mu: &sync.Mutex{},
	}
	return nil
}

var (
	// ErrFSFileNotPresentError ...
	ErrFSFileNotPresentError = errors.New("file is not present in store")
)

func (fs *fileStore) deleteFile(name string) error {
	if !fs.filePresent(name) {
		return ErrFSFileNotPresentError
	}

	if err := os.Remove(fs.basePath + name); err != nil {
		return err
	}
	delete(fs.files, name)
	return nil
}

func (fs *fileStore) filePresent(name string) bool {
	if _, present := fs.files[name]; !present {
		return false
	}
	return true
}

func (fs *fileStore) WriteAsJSON(name string, obj interface{}) error {
	if !fs.filePresent(name) {
		if err := fs.createFile(name); err != nil {
			return err
		}
	}
	fl := fs.files[name]
	fl.mu.Lock()
	defer fl.mu.Unlock()

	json, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fs.basePath+name, json, 0644)
}

func (fs *fileStore) ReadAsJSON(name string, obj interface{}) error {
	if !fs.filePresent(name) {
		return ErrFSFileNotPresentError
	}

	fl := fs.files[name]
	fl.mu.Lock()
	defer fl.mu.Unlock()

	data, err := ioutil.ReadFile(fs.basePath + name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}

func (fs *fileStore) AppendToFile(name, line string) error {
	if !fs.filePresent(name) {
		if err := fs.createFile(name); err != nil {
			return err
		}
	}
	line += "\n"

	fl := fs.files[name]
	fl.mu.Lock()
	defer fl.mu.Unlock()

	fileToWrite, err := os.OpenFile(fs.basePath+name, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer fileToWrite.Close()

	if _, err := fileToWrite.WriteString(line); err != nil {
		return err
	}
	return nil
}

func (fs *fileStore) ReadLineByLine(name string, fn func(string) bool) error {
	if !fs.filePresent(name) {
		return ErrFSFileNotPresentError
	}

	fl := fs.files[name]
	fl.mu.Lock()
	defer fl.mu.Unlock()

	fileToRead, err := os.Open(fs.basePath + name)
	if err != nil {
		return err
	}
	defer fileToRead.Close()

	scanner := bufio.NewScanner(fileToRead)
	for scanner.Scan() {
		text := scanner.Text()
		if !fn(text) {
			break
		}

	}
	return nil
}
