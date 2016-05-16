package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"

	maestroConfig "github.com/cpg1111/maestro/config"
)

type server struct {
	Runtime  string
	Host     string
	Port     uint
	CloneDir string
}

// Config is the struct of the config file
type Config struct {
	Server   server
	Projects []maestroConfig.Project
}

// Load loads a config file and returns a pointer to a config struct
func Load(path string) (*Config, error) {
	var conf Config
	confData, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return nil, readErr
	}
	if _, pErr := toml.Decode((string)(confData), &conf); pErr != nil {
		return nil, pErr
	}
	return &conf, nil
}
