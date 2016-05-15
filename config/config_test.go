package config

import (
	"fmt"
	"testing"

	"github.com/fatih/structs"

	maestroConfig "github.com/cpg1111/maestro/config"
)

var expected = &Config{
	Runtime:  "Native",
	Bind:     "0.0.0.0",
	Port:     8080,
	CloneDir: "/tmp/",
	Projects: []maestroConfig.Project{
		maestroConfig.Project{
			RepoURL:        "git@github.com:cpg1111/maestro.git",
			CloneCMD:       "git clone",
			AuthType:       "SSH",
			SSHPrivKeyPath: "~/.ssh/id_rsa",
			SSHPubKeyPath:  "~/.ssh/id_rsa.pub",
			Username:       "git",
			Password:       "",
			PromptForPWD:   false,
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
			expectedArr := expectedMap[i].([]maestroConfig.Project)
			testArr := testMap[i].([]maestroConfig.Project)
			if len(expectedArr) != len(testArr) {
				t.Error("Did not load expected amount of projects")
			}
		}
		if testMap[i] == nil || testMap[i] != expectedMap[i] {
			t.Error(fmt.Errorf("Expected %v found %v for field %s", expectedMap[i], testMap[i], i))
		}
	}
}
