package storage

import (
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
	name string
	mu   *sync.Mutex
}

func newFileStore(base string) (*fileStore, error) {
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
			name: f.Name(),
			mu:   &sync.Mutex{},
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
		name: name,
	}
	return nil
}

var (
	// ErrFSFileNotPresentError ...
	ErrFSFileNotPresentError = errors.New("file is not present in store")
)

func (fs *fileStore) deleteFile(name string) error {
	if _, present := fs.files[name]; !present {
		return ErrFSFileNotPresentError
	}

	if err := os.Remove(fs.basePath + name); err != nil {
		return err
	}
	delete(fs.files, name)
	return nil
}
