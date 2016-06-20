package docker

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	driver, err := New(os.Getenv("DOCKER_HOST"), "v1.23", "0.1.1", "maestro")
	if err != nil {
		t.Error(err)
	}
	err = driver.Run([]string{
		"maestro",
		"--branch=master",
		"--deploy=true",
		"--prev-commit=f0dfac3dd5efdb0c80a2321f5a2a69c0bc3cb67f",
		"--config=./test_conf.toml",
		"--clone-path=./clonetest/",
	})
	if err != nil {
		t.Error(err)
	}
}
