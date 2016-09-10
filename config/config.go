package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Mount struct {
	Kind      string
	ID        string
	FSType    string
	Server    string
	Endpoints string
	Path      string
	Name      string
	ReadOnly  bool
}

type Server struct {
	Runtime             string
	RuntimeTLSClient    bool
	RuntimeTLSServer    bool
	MaxBuilds           int
	DataStoreType       string
	DataStoreUser       string
	DataStorePWD        string
	DataStoreEnvIP      string
	DataStoreStaticIP   string
	DataStoreEnvPort    string
	DataStoreStaticPort string
	TargetHost          string
	ClientCertPath      string
	ClientKeyPath       string
	ServerCertPath      string
	ServerKeyPath       string
	MaestroVersion      string
	Host                string
	SecurePort          uint
	InsecurePort        uint
	StateComPort        uint
	WorkspaceDir        string
}

type Project struct {
	Name            string   `json:"name"`
	MaestroConfPath string   `json:"confPath"`
	DeployBranches  []string `json:"deployBranches"`
}

// Config is the struct of the config file
type Config struct {
	Server   Server
	Projects []Project
	Mounts   []Mount
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
