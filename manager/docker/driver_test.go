package docker

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
