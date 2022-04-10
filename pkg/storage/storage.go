package storage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/erizzardi/storage/util"
	"github.com/sirupsen/logrus"
)

type storageService struct {
	Logger *logrus.Logger
}

func NewService(logger *logrus.Logger) Service { return &storageService{Logger: logger} }

//----------------------------------------------
// This is where the API methods are implemented
//----------------------------------------------
func (ss *storageService) WriteFile(ctx context.Context, file io.Reader, fileName string, storageFolder string) error {
	ss.Logger.Debug("Method WriteFile invoked.")
	fileName = filepath.Join(storageFolder, fileName)

	var err error

	// check if file exists
	if _, err = os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.Logger.Debug("Creating file " + fileName)
		newFile, err := os.Create(fileName)
		if err != nil {
			ss.Logger.Error("Error: " + err.Error())
			return &util.ResponseError{
				StatusCode: 500,
				Err:        errors.New(err.Error()),
			}
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			ss.Logger.Error("Error: " + err.Error())
			return &util.ResponseError{
				StatusCode: 500,
				Err:        errors.New(err.Error()),
			}
		}
		ss.Logger.Info("File " + fileName + " created successfully")

	} else {
		ss.Logger.Error("Error: file already exists")
		return &util.ResponseError{
			StatusCode: 409,
			Err:        errors.New("file already exists"),
		}
	}
	return nil
}

func (ss *storageService) GetFile(ctx context.Context, fileName string, storageFolder string) ([]byte, error) {
	fileName = filepath.Join(storageFolder, fileName)
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.Logger.Error("Error: " + err.Error())
		return nil, &util.ResponseError{
			StatusCode: 404,
			Err:        errors.New(err.Error()),
		}
	}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		ss.Logger.Error("Error: " + err.Error())
		return nil, &util.ResponseError{
			StatusCode: 500,
			Err:        errors.New(err.Error()),
		}
	}
	ss.Logger.Info("File " + fileName + " retrieved correctly")
	return file, nil
}
