package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	fixturesBase     = "fixtures"
	tempTestDataBase = "tempTestData"
)

func TestNewFileStore(t *testing.T) {
	_, err := NewFileStore("dummypath")
	assert.NotNil(t, err)

	fs, err := NewFileStore(fixturesBase)
	assert.Nil(t, err)
	assert.NotNil(t, fs)
	assert.Equal(t, fixturesBase+"/", fs.basePath)

	fs, err = NewFileStore(fixturesBase + "/")
	assert.Nil(t, err)
	assert.NotNil(t, fs)
	assert.Equal(t, fixturesBase+"/", fs.basePath)

	assert.Equal(t, len(fs.files), 3)
	for k, v := range fs.files {
		switch k {
		case "t0":
			assert.NotNil(t, v)
		case "t1":
			assert.NotNil(t, v)
		case "t2":
			assert.NotNil(t, v)
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

	fs, err := NewFileStore(tempTestDataBase)
	assert.Nil(t, err)
	assert.NotNil(t, fs)

	assert.Nil(t, fs.createFile("alfa"))
	assert.Equal(t, 1, len(fs.files))
	assert.FileExists(t, tempTestDataBase+"/alfa")
	assert.NotNil(t, fs.files["alfa"])

	assert.Nil(t, fs.createFile("beta"))
	assert.Equal(t, 2, len(fs.files))
	assert.FileExists(t, tempTestDataBase+"/beta")
	assert.NotNil(t, fs.files["beta"])

	// already created file
	assert.NotNil(t, fs.createFile("beta"))
}

func TestFileStoreDeleteFile(t *testing.T) {
	assert.Nil(t, removeTempTestFiles())

	fs, _ := NewFileStore(tempTestDataBase)
	assert.Nil(t, fs.createFile("teta"))
	assert.Nil(t, fs.deleteFile("teta"))
	assert.Equal(t, 0, len(fs.files))
	assert.NotNil(t, fs.deleteFile("notpresent"))
}

type MockForJSON struct {
	ID    string
	Value float64
}

func TestFileStoreWriteAsJSON(t *testing.T) {
	assert.Nil(t, removeTempTestFiles())

	fs, _ := NewFileStore(tempTestDataBase)
	mockObject0 := &MockForJSON{
		ID:    "1",
		Value: 1.2,
	}

	mockObject1 := &MockForJSON{
		ID:    "2",
		Value: 1.3,
	}

	assert.Nil(t, fs.WriteAsJSON("mock0", mockObject0))

	dat, err := ioutil.ReadFile(tempTestDataBase + "/mock0")
	assert.Nil(t, err)
	assert.Equal(t, "{\"ID\":\"1\",\"Value\":1.2}", string(dat))

	assert.Nil(t, fs.WriteAsJSON("mock0", mockObject1))

	dat, err = ioutil.ReadFile(tempTestDataBase + "/mock0")
	assert.Nil(t, err)
	assert.Equal(t, "{\"ID\":\"2\",\"Value\":1.3}", string(dat))
}

func TestFileStoreReadAsJSON(t *testing.T) {
	assert.Nil(t, removeTempTestFiles())

	fs, _ := NewFileStore(tempTestDataBase)
	mockObject0 := &MockForJSON{
		ID:    "3",
		Value: 1.4,
	}
	assert.Nil(t, fs.WriteAsJSON("mock1", mockObject0))

	mockObject1 := &MockForJSON{}
	assert.Nil(t, fs.ReadAsJSON("mock1", mockObject1))
	assert.Equal(t, "3", mockObject1.ID)
	assert.Equal(t, 1.4, mockObject1.Value)
}

func TestFileStoreAppendToFile(t *testing.T) {
	assert.Nil(t, removeTempTestFiles())

	fs, _ := NewFileStore(tempTestDataBase)
	assert.Nil(t, fs.AppendToFile("alfa", "beta"))
	assert.FileExists(t, tempTestDataBase+"/alfa")

	dat, err := ioutil.ReadFile(tempTestDataBase + "/alfa")
	assert.Nil(t, err)
	assert.Equal(t, "beta\n", string(dat))

	assert.Nil(t, fs.AppendToFile("alfa", "teta"))

	dat, err = ioutil.ReadFile(tempTestDataBase + "/alfa")
	assert.Nil(t, err)
	assert.Equal(t, "beta\nteta\n", string(dat))
}

func TestFileStoreReadLineByLine(t *testing.T) {
	assert.Nil(t, removeTempTestFiles())

	fs, _ := NewFileStore(tempTestDataBase)
	assert.Nil(t, fs.AppendToFile("alfa", "beta"))
	assert.Nil(t, fs.AppendToFile("alfa", "teta"))
	assert.Nil(t, fs.AppendToFile("alfa", "omega"))

	linesRead0 := make([]string, 0, 3)
	assert.Nil(t, fs.ReadLineByLine("alfa", func(line string) bool {
		linesRead0 = append(linesRead0, line)
		return true
	}))

	assert.Equal(t, "beta", linesRead0[0])
	assert.Equal(t, "teta", linesRead0[1])
	assert.Equal(t, "omega", linesRead0[2])

	linesRead1 := make([]string, 0, 3)
	assert.Nil(t, fs.ReadLineByLine("alfa", func(line string) bool {
		linesRead1 = append(linesRead1, line)
		return line != "teta"
	}))

	assert.Equal(t, 2, len(linesRead1))
	assert.Equal(t, "beta", linesRead1[0])
	assert.Equal(t, "teta", linesRead1[1])
}
