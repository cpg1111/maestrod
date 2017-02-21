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

package driver

import (
	"fmt"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	driver, err := New("v1.23", "latest")
	if err != nil {
		t.Error(err)
	}
	err = driver.Run("test", "/etc/maestro/", fmt.Sprintf("%s/src/github.com/cpg1111/maestro", os.Getenv("GOPATH")), []string{
		"maestro",
		"--branch=master",
		"--deploy=true",
		"--prev-commit=f0dfac3dd5efdb0c80a2321f5a2a69c0bc3cb67f",
		"--config=/etc/maestro/test_conf.toml",
		"--clone-path=./clonetest/",
	})
	if err != nil {
		t.Error(err)
	}
}
