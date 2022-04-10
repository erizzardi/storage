package util

// Config struct definition
type Config struct {
	StorageFolder string
}

// Config constructor
func SetConfig(storageFolder string) *Config { return &Config{StorageFolder: storageFolder} }
