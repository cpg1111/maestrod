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

package k8s

import (
	"fmt"
	"os"
	"testing"

	"github.com/cpg1111/maestrod/config"
)

var K8S = os.Getenv("TEST_K8S_URL")

var conf, confErr = config.Load(fmt.Sprintf("%s/src/github.com/cpg1111/maestrod/example.conf.toml", os.Getenv("GOPATH")))

var driver = New(os.Getenv("TEST_MAESTRO_VER"), conf)

func TestRun(t *testing.T) {
	if confErr != nil {
		t.Error(confErr)
	}
	conf.Server.ClientCertPath = os.Getenv("TEST_CLIENT_CERT")
	conf.Server.ClientKeyPath = os.Getenv("TEST_CLIENT_KEY")
	branch := os.Getenv("TEST_BRANCH")
	confPath := os.Getenv("TEST_CONF_PATH")
	prevCommit := os.Getenv("TEST_PREV_COMMIT")
	currCommit := os.Getenv("TEST_CURR_COMMIT")
	clonePath := "/tmp/test/"
	nsErr := driver.CreateNamespace("maestro")
	if nsErr != nil {
		t.Error(nsErr)
	}
	saErr := driver.CreateSvcAccnt("default")
	if saErr != nil {
		t.Error(saErr)
	}
	wd, wdErr := os.Getwd()
	if wdErr != nil {
		t.Error(wdErr)
	}
	runErr := driver.Run("test", confPath, wd, []string{
		"maestro",
		fmt.Sprintf("--branch=%s", branch),
		fmt.Sprintf("--deploy=%v", false),
		fmt.Sprintf("--prev-commit=%s", prevCommit),
		fmt.Sprintf("--curr-commit=%s", currCommit),
		fmt.Sprintf("--config=%s", "/etc/maestro/maestrod.toml"),
		fmt.Sprintf("--clone-path=%s", clonePath),
	})
	if runErr != nil {
		t.Error(runErr)
	}
	podURL := fmt.Sprintf("%s/namespaces/maestro/pods/test", K8S)
	resp, getErr := driver.Client.Get(podURL)
	if getErr != nil {
		t.Error(getErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error("maestro test pod not found")
	}
}
