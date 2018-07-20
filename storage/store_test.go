package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	fixturesBase     = "fixtures"
	tempTestDataBase = "tempTestData"
)

func TestNewFileStore(t *testing.T) {
	_, err := newFileStore("dummypath")
	assert.NotNil(t, err)

	fs, err := newFileStore(fixturesBase)
	assert.Nil(t, err)
	assert.NotNil(t, fs)
	assert.Equal(t, fixturesBase+"/", fs.basePath)

	fs, err = newFileStore(fixturesBase + "/")
	assert.Nil(t, err)
	assert.NotNil(t, fs)
	assert.Equal(t, fixturesBase+"/", fs.basePath)

	assert.Equal(t, len(fs.files), 3)
	for k, v := range fs.files {
		switch k {
		case "t0":
			assert.Equal(t, "t0", v.name)
		case "t1":
			assert.Equal(t, "t1", v.name)
		case "t2":
			assert.Equal(t, "t2", v.name)
		}
	}
}

func removeTempTestFiles() error {
	if err := os.RemoveAll(tempTestDataBase); err != nil {
		return err
	}
	return os.Mkdir(tempTestDataBase, os.ModePerm)
}

func TestFileStoreCreateFile(t *testing.T) {
	assert.Nil(t, removeTempTestFiles())

	fs, err := newFileStore(tempTestDataBase)
	assert.Nil(t, err)
	assert.NotNil(t, fs)

	assert.Nil(t, fs.createFile("alfa"))
	assert.Equal(t, 1, len(fs.files))
	assert.FileExists(t, tempTestDataBase+"/alfa")
	assert.Equal(t, "alfa", fs.files["alfa"].name)

	assert.Nil(t, fs.createFile("beta"))
	assert.Equal(t, 2, len(fs.files))
	assert.FileExists(t, tempTestDataBase+"/beta")
	assert.Equal(t, "beta", fs.files["beta"].name)

	// already created file
	assert.NotNil(t, fs.createFile("beta"))
}

func TestFileStoreDeleteFile(t *testing.T) {
	assert.Nil(t, removeTempTestFiles())

	fs, _ := newFileStore(tempTestDataBase)
	assert.Nil(t, fs.createFile("teta"))
	assert.Nil(t, fs.deleteFile("teta"))
	assert.Equal(t, 0, len(fs.files))
	assert.NotNil(t, fs.deleteFile("notpresent"))
}
