package config

import (
	"fmt"
	"testing"

	"github.com/fatih/structs"
)

var expected = &Config{
	Server: Server{
		Runtime:      "docker",
		Host:         "0.0.0.0",
		Port:         8080,
		WorkspaceDir: "/tmp/",
	},
	Projects: []Project{
		Project{
			Name:            "maestrod",
			MaestroConfPath: "../maestro/test_conf.toml",
			DeployBranches:  []string{"master"},
		},
	},
}

func TestLoad(t *testing.T) {
	conf, loadErr := Load("../test_conf.toml")
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
		}
		if testMap[i] == nil || testMap[i] != expectedMap[i] {
			t.Error(fmt.Errorf("Expected %v found %v for field %s", expectedMap[i], testMap[i], i))
		}
	}
}
