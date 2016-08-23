package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/fatih/structs"
)

var expected = &Config{
	Server: Server{
		Runtime:          "docker",
		RuntimeTLSClient: true,
		RuntimeTLSServer: true,
		MaxBuilds:        6,
		DataStoreType:    "redis",
		DataStoreUser:    "",
		DataStorePWD:     "",
		DataStoreEnvIP:   "REDIS_SERVICE_HOST",
		DataStoreEnvPort: "REDIS_SERVICE_PORT",
		TargetHost:       "",
		MaestroVersion:   "0.1.1",
		ClientCertPath:   "/etc/maestrod/clientcrt.pem",
		ClientKeyPath:    "/etc/maestrod/clientkey.pem",
		ServerCertPath:   "/etc/maestrod/fullchain.pem",
		ServerKeyPath:    "/etc/maestrod/privkey.pem",
		Host:             "0.0.0.0",
		SecurePort:       8484,
		InsecurePort:     8585,
		WorkspaceDir:     "/tmp/",
	},
	Projects: []Project{
		Project{
			Name:            "maestrod",
			MaestroConfPath: "/etc/maestro/test_conf.toml",
			DeployBranches:  []string{"master"},
		},
	},
}

func TestLoad(t *testing.T) {
	conf, loadErr := Load(fmt.Sprintf("%ssrc/github.com/cpg1111/maestrod/example.conf.toml", os.Getenv("GOPATH")))
	if loadErr != nil {
		t.Error(loadErr)
	}
	expectedMap := structs.Map(*expected)
	testMap := structs.Map(*conf)
	for i := range expectedMap {
		if i == "Projects" {
			expectedArr := expectedMap[i].([]Project)
			testArr := testMap[i].([]Project)
			if len(expectedArr) != len(testArr) {
				t.Error("Did not load expected amount of projects")
			}
		} else if i == "Server" {
			expectedSrvMap := expectedMap[i].(map[string]interface{})
			testSrvMap := testMap[i].(map[string]interface{})
			for j := range expectedSrvMap {
				if testSrvMap[j] == nil || testSrvMap[j] != expectedSrvMap[j] {
					t.Error(fmt.Errorf("Expected %v found %v for field %s", expectedSrvMap[j], testSrvMap[j], j))
				}
			}
		}
	}
}
