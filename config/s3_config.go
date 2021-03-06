package config

import (
	"fmt"

	"github.com/smartystreets/raptr/storage"
)

type s3Config struct {
	RegionName string `json:"region"`
	BucketName string `json:"bucket"`
	PathPrefix string `json:"prefix"`
	LayoutName string `json:"layout"`
	MaxRetries int    `json:"retries"`
	Timeout    int    `json:"timeout"`
}

func (this s3Config) validate() error {
	if this.BucketName == "" {
		return fmt.Errorf("The bucket name is missing.")
	} else {
		return nil
	}
}

func (this s3Config) buildStorage() (storage.Storage, error) {
	actual := storage.NewS3Storage(this.RegionName, this.BucketName, this.PathPrefix)
	if this.MaxRetries <= 0 {
		this.MaxRetries = defaultMaxRetries
	}

	inner := storage.Storage(actual)
	inner = storage.NewIntegrityStorage(inner)
	inner = storage.NewRetryStorage(inner, this.MaxRetries)
	inner = storage.NewConcurrentStorage(inner)
	return inner, nil
}

const defaultMaxRetries = 3
