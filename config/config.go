package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Runtime          string
	RuntimeTLSClient bool
	RuntimeTLSServer bool
	ClientCertPath   string
	ClientKeyPath    string
	ServerCertPath   string
	ServerKeyPath    string
	MaestroVersion   string
	Host             string
	Port             uint
	WorkspaceDir     string
}

type Project struct {
	Name            string
	MaestroConfPath string
	DeployBranches  []string
}

// Config is the struct of the config file
type Config struct {
	Server   Server
	Projects []Project
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
