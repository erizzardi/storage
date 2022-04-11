package storage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/erizzardi/storage/base"
	"github.com/erizzardi/storage/util"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type storageService struct {
	db     base.DB
	logger *logrus.Logger
}

func NewService(db base.DB, logger *logrus.Logger) Service {
	return &storageService{db: db, logger: logger}
}

//----------------------------------------------
// This is where the API methods are implemented
//----------------------------------------------
func (ss *storageService) WriteFile(ctx context.Context, file io.Reader, metadata util.Metadata, storageFolder string) (string, error) {
	ss.logger.Debug("Method WriteFile invoked.")

	uuid := uuid.New().String()
	fileName := filepath.Join(storageFolder, uuid)

	// check if file exists
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.logger.Debug("Creating file " + fileName + "...")
		newFile, err := os.Create(fileName)
		if err != nil {
			ss.logger.Error("Error: " + err.Error())
			return "", &util.ResponseError{
				StatusCode: 500,
				Err:        errors.New(err.Error()),
			}
		}
		defer newFile.Close()
		ss.logger.Debug("Created file " + fileName)
		ss.logger.Debug("Copying file content to new destination...")
		if _, err := io.Copy(newFile, file); err != nil {
			ss.logger.Error("Error: " + err.Error())
			return "", &util.ResponseError{
				StatusCode: 500,
				Err:        errors.New(err.Error()),
			}
		}
		ss.logger.Debug("File content copied")

		ss.logger.Debug(uuid, metadata.Name)

		// write metadata to db
		err = ss.db.InsertMetadata(base.Row{
			Uuid:     uuid,
			FileName: metadata.Name,
		})
		if err != nil {
			ss.logger.Error("Error: " + err.Error())
			return "", &util.ResponseError{
				StatusCode: 500,
				Err:        errors.New(err.Error()),
			}
		}

		ss.logger.Info("File " + uuid + " created successfully")

	} else {
		ss.logger.Error("Error: file already exists")
		return "", &util.ResponseError{
			StatusCode: 409,
			Err:        errors.New("file already exists"),
		}
	}
	return uuid, nil
}

func (ss *storageService) GetFile(ctx context.Context, uuid string, storageFolder string) ([]byte, error) {
	ss.logger.Debug("Method GetFile invoked.")

	fileName := filepath.Join(storageFolder, uuid)
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.logger.Error("Error: " + err.Error())
		return nil, &util.ResponseError{
			StatusCode: 404,
			Err:        errors.New(err.Error()),
		}
	}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		ss.logger.Error("Error: " + err.Error())
		return nil, &util.ResponseError{
			StatusCode: 500,
			Err:        errors.New(err.Error()),
		}
	}
	ss.logger.Info("File " + uuid + " retrieved successfully")
	return file, nil
}

func (ss *storageService) DeleteFile(ctx context.Context, uuid string, storageFolder string) error {
	ss.logger.Debug("Method DeleteFile invoked.")
	fileName := filepath.Join(storageFolder, uuid)

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.logger.Error("Error: " + err.Error())
		return &util.ResponseError{
			StatusCode: 404,
			Err:        errors.New(err.Error()),
		}
	}
	if err := os.Remove(fileName); err != nil {
		ss.logger.Error("Error: " + err.Error())
		return &util.ResponseError{
			StatusCode: 500,
			Err:        errors.New(err.Error()),
		}
	}
	ss.logger.Info("File " + fileName + "deleted successfully")
	return nil
}
