package storage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/erizzardi/storage/util"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type storageService struct {
	Logger *logrus.Logger
}

func NewService(logger *logrus.Logger) Service { return &storageService{Logger: logger} }

//----------------------------------------------
// This is where the API methods are implemented
//----------------------------------------------
func (ss *storageService) WriteFile(ctx context.Context, file io.Reader, storageFolder string) (string, error) {
	ss.Logger.Debug("Method WriteFile invoked.")

	uuid := uuid.New().String()
	fileName := filepath.Join(storageFolder, uuid)

	// check if file exists
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.Logger.Debug("Creating file " + fileName)
		newFile, err := os.Create(fileName)
		if err != nil {
			ss.Logger.Error("Error: " + err.Error())
			return "", &util.ResponseError{
				StatusCode: 500,
				Err:        errors.New(err.Error()),
			}
		}
		defer newFile.Close()

		if _, err := io.Copy(newFile, file); err != nil {
			ss.Logger.Error("Error: " + err.Error())
			return "", &util.ResponseError{
				StatusCode: 500,
				Err:        errors.New(err.Error()),
			}
		}
		ss.Logger.Info("File " + uuid + " created successfully")

	} else {
		ss.Logger.Error("Error: file already exists")
		return "", &util.ResponseError{
			StatusCode: 409,
			Err:        errors.New("file already exists"),
		}
	}
	return uuid, nil
}

func (ss *storageService) GetFile(ctx context.Context, uuid string, storageFolder string) ([]byte, error) {
	ss.Logger.Debug("Method GetFile invoked.")

	fileName := filepath.Join(storageFolder, uuid)
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
	ss.Logger.Info("File " + uuid + " retrieved successfully")
	return file, nil
}

func (ss *storageService) DeleteFile(ctx context.Context, uuid string, storageFolder string) error {
	ss.Logger.Debug("Method DeleteFile invoked.")
	fileName := filepath.Join(storageFolder, uuid)

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.Logger.Error("Error: " + err.Error())
		return &util.ResponseError{
			StatusCode: 404,
			Err:        errors.New(err.Error()),
		}
	}
	if err := os.Remove(fileName); err != nil {
		ss.Logger.Error("Error: " + err.Error())
		return &util.ResponseError{
			StatusCode: 500,
			Err:        errors.New(err.Error()),
		}
	}
	ss.Logger.Info("File " + fileName + "deleted successfully")
	return nil
}
