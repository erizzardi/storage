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

// ListFiles list metadata, paged.
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

// WriteFile writes a file to disk, and updates metadata in DB.
// Returns 200, 400, 409, 500
func (ss *storageService) WriteFile(ctx context.Context, file io.Reader, metadata util.Metadata, storageFolder string) (string, error) {
	ss.logger.Debug("Method WriteFile invoked.")

	uuid := uuid.New().String()
	fileName := filepath.Join(storageFolder, uuid)

	if file == nil {
		ss.logger.Error("Error: no file in request")
		return "", util.BadRequestError{Message: "no file in request"}
	}

	// Check if file exists by querying the DB by fileName.
	// A filesystem check should not be necessary, since UUIDs are unique.
	if _, err := ss.db.RetrieveMetadata("filename", metadata.Name); errors.Is(err, base.NotFoundError) {
		ss.logger.Debug("Creating file " + fileName + "...")
		newFile, err := os.Create(fileName)
		if err != nil {
			ss.logger.Error("Error: " + err.Error())
			return "", util.InternalServerError{}
		}
		defer newFile.Close()
		ss.logger.Debug("Created file " + fileName)
		ss.logger.Debug("Copying file content to new destination...")
		if _, err := io.Copy(newFile, file); err != nil {
			ss.logger.Error("Error: " + err.Error())
			_ = os.Remove(fileName)
			return "", util.InternalServerError{}
		}
		ss.logger.Debug("File content copied")

		ss.logger.Debug(uuid, metadata.Name)

		// Write metadata to db
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
		ss.logger.Error("file already exists")
		return "", util.ConflictError{Message: "file already exists"}
	}

	return uuid, nil
}

// GetFile returns the metadata of a file from its Uuid.
// Returns 200, 404, 500
func (ss *storageService) GetFile(ctx context.Context, uuid string, storageFolder string) ([]byte, error) {
	ss.logger.Debug("Method GetFile invoked.")

	// Check db for entry corresponding to file
	row, err := ss.db.RetrieveMetadata("uuid", uuid)
	if err != nil {
		ss.logger.Error("Error: " + err.Error())
		return nil, util.InternalServerError{Message: err.Error()}
	}
	if row == (util.Row{}) {
		ss.logger.Errorf("Error: file %s not found", uuid)
		return nil, util.NotFoundError{Message: "file not found"}
	}

	fileName := filepath.Join(storageFolder, uuid)
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		ss.logger.Error("Error: " + err.Error())
		return nil, util.NotFoundError{Message: err.Error()}
	}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		ss.logger.Error("Error: " + err.Error())
		return nil, util.InternalServerError{Message: err.Error()}
	}
	ss.logger.Info("File " + uuid + " retrieved successfully")
	return file, nil
}

// DeleteFile deletes a file from disk by its Uuid.
// Returns 200, 404, 500
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
	if err := ss.db.DeleteMetadata("uuid", uuid); err != nil {
		ss.logger.Error("Error: " + err.Error())

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
