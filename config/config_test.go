/*
Copyright 2016 Christian Grabowski All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
		DataStoreEnvIP:   "ETCD_SERVICE_HOST",
		DataStoreEnvPort: "ETCD_SERVICE_PORT",
		TargetProtocol:   "https",
		TargetHost:       "",
		TargetPort:       "",
		TargetEnvHost:    "KUBERNETES_SERVICE_HOST",
		TargetEnvPort:    "KUBERNETES_SERVICE_PORT",
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
