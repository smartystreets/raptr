package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/smartystreets/raptr/conf"
)

type Configuration struct {
	repos map[string]RepositoryConfig
}

func (this Configuration) Open(repositoryName string) (RepositoryConfig, bool) {
	repo, found := this.repos[repositoryName]
	return repo, found
}
func LoadConfiguration(fullPath string) (Configuration, error) {
	deserializedFile := configFile{}
	if deserialized, err := readFile(fullPath); !os.IsNotExist(err) {
		return deserialized, err // the result wasn't a "file not found" issue
	} else if reader, err := conf.Read(".raptr", "raptr.conf"); err != nil {
		return Configuration{}, err // can't find any kind of config file
	} else if err := json.NewDecoder(reader).Decode(&deserializedFile); err != nil {
		return Configuration{}, err // unable to deserialize
	} else {
		return newConfiguration(deserializedFile)
	}
}
func readFile(fullPath string) (Configuration, error) {
	deserialized := configFile{}
	if len(strings.TrimSpace(fullPath)) == 0 {
		return Configuration{}, os.ErrNotExist
	} else if handle, err := os.Open(filepath.Clean(fullPath)); err != nil {
		return Configuration{}, err // file doesn't exist or access problems
	} else if contents, err := ioutil.ReadAll(handle); err != nil {
		handle.Close()
		return Configuration{}, err // couldn't read file
	} else if err := json.Unmarshal(contents, &deserialized); err != nil {
		handle.Close()
		return Configuration{}, err // malformed JSON
	} else {
		handle.Close()
		return newConfiguration(deserialized)
	}
}
func newConfiguration(file configFile) (Configuration, error) {
	repos := map[string]RepositoryConfig{}
	layouts := map[string]RepositoryLayout{}

	for key, item := range file.Layouts {
		if err := item.validate(); err != nil {
			return Configuration{}, fmt.Errorf("Layout '%s' has missing or corrupt values.", key)
		} else {
			layouts[key] = item
		}
	}

	for key, item := range file.S3 {
		if layout, found := layouts[item.LayoutName]; !found {
			return Configuration{}, fmt.Errorf("S3 store '%s' references not-existent layout '%s'.", key, item.LayoutName)
		} else if err := item.validate(); err != nil {
			return Configuration{}, fmt.Errorf("S3 store '%s' has missing or corrupt values.", key)
		} else if store, err := item.buildStorage(); err != nil {
			return Configuration{}, fmt.Errorf("S3 store '%s' cannot be initialized.", key)
		} else {
			repos[key] = RepositoryConfig{Storage: store, Layout: layout}
		}
	}

	return Configuration{repos: repos}, nil
}
