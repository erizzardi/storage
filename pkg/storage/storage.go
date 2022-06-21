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
)

// storageService implements the storage.Service interface
type storageService struct {
	// Pointer to a DB interface, that allows DB operations.
	db base.DB
	// Logger specific for the business logic layer
	logger *util.Logger
	// Map[layer]logger. To change logging level at run time
	layerLoggersMap map[string]*util.Logger
}

func NewService(db base.DB, logger *util.Logger, layerLoggersMap map[string]*util.Logger) Service {
	return &storageService{db: db, logger: logger, layerLoggersMap: layerLoggersMap}
}

//===================================================================================
// This is where the API methods are implemented
//-----------------------------------------------------------------------------------
// Error handling: this layer doesn't know anything about response codes or messages.
// This just returns the type of error.
//===================================================================================

//
// ListFiles list metadata, paged
// returns 200, 500
func (ss *storageService) ListFiles(ctx context.Context, limit uint, offset uint) ([]util.Row, error) {

	ss.logger.Debug("Method ListFiles invoked.")
	rows, err := ss.db.ListAllPaged(limit, offset)
	if err != nil {
		ss.logger.Error(err.Error())
		return nil, util.InternalServerError{}
	}
	return rows, nil
}

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
			return "", util.InternalServerError{}
		}
		defer newFile.Close()
		ss.logger.Debug("Created file " + fileName)
		ss.logger.Debug("Copying file content to new destination...")
		if _, err := io.Copy(newFile, file); err != nil { // why copy?
			ss.logger.Error("Error: " + err.Error())
			return "", util.InternalServerError{}
		}
		ss.logger.Debug("File content copied")

		ss.logger.Debug(uuid, metadata.Name)

		// write metadata to db
		err = ss.db.InsertMetadata(util.Row{
			Uuid:     uuid,
			FileName: metadata.Name,
		})
		if err != nil {
			ss.logger.Error("Error: " + err.Error())
			return "", util.InternalServerError{}
		}

		ss.logger.Info("File " + uuid + " created successfully")

	} else {
		ss.logger.Error("Error: file already exists")
		return "", util.ConflictError{}
	}
	return uuid, nil
}

func (ss *storageService) GetFile(ctx context.Context, uuid string, storageFolder string) ([]byte, error) {
	ss.logger.Debug("Method GetFile invoked.")

	// Check db for entry corresponding to file
	row, err := ss.db.RetrieveMetadata(util.Row{Uuid: uuid, FileName: ""})
	if err != nil {
		ss.logger.Error("Error: " + err.Error())
		return nil, util.InternalServerError{}
	}
	if row == (util.Row{}) {
		ss.logger.Error("Error: file not found")
		return nil, util.NotFoundError{}
	}

	fileName := filepath.Join(storageFolder, uuid)
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.logger.Error("Error: " + err.Error())
		return nil, util.NotFoundError{}
	}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		ss.logger.Error("Error: " + err.Error())
		return nil, util.InternalServerError{}
	}
	ss.logger.Info("File " + uuid + " retrieved successfully")
	return file, nil
}

func (ss *storageService) DeleteFile(ctx context.Context, uuid string, storageFolder string) error {
	ss.logger.Debug("Method DeleteFile invoked.")
	fileName := filepath.Join(storageFolder, uuid)

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.logger.Error("Error: " + err.Error())
		return util.NotFoundError{}
	}
	if err := os.Remove(fileName); err != nil {
		ss.logger.Error("Error: " + err.Error())
		return util.InternalServerError{}
	}
	ss.logger.Info("File " + fileName + "deleted successfully")
	return nil
}

func (ss *storageService) SetLogLevel(ctx context.Context, layer string, level string) error {
	ss.logger.Debug("Method SetLogLevel invoked")
	logLevel, err := util.LogLevelMapping(level)
	if err != nil {
		return util.BadRequestError{Message: err.Error()}
	}
	if logger, ok := ss.layerLoggersMap[layer]; ok {
		logger.SetLevel(logLevel)
	} else {
		return util.BadRequestError{Message: "invalid layer: " + layer}
	}
	return nil
}
