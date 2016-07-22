package k8s

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

var K8S = os.Getenv("TEST_K8S_URL")

var driver = New(K8S, os.Getenv("TEST_MAESTRO_VER"))

func TestRun(t *testing.T) {
	branch := os.Getenv("TEST_BRANCH")
	confPath := os.Getenv("TEST_CONF_PATH")
	prevCommit := os.Getenv("TEST_PREV_COMMIT")
	currCommit := os.Getenv("TEST_CURR_COMMIT")
	clonePath := "/tmp/test/"
	runErr := driver.Run("test", confPath, confPath, []string{
		"maestro",
		fmt.Sprintf("--branch=%s", branch),
		fmt.Sprintf("--deploy=%v", false),
		fmt.Sprintf("--prev-commit=%s", prevCommit),
		fmt.Sprintf("--curr-commit=%s", currCommit),
		fmt.Sprintf("--config=%s", confPath),
		fmt.Sprintf("--clone-path=%s", clonePath),
	})
	if runErr != nil {
		t.Error(runErr)
	}
	podURL := fmt.Sprintf("%s/namespaces/maestro/pods/test", K8S)
	resp, getErr := http.Get(podURL)
	if getErr != nil {
		t.Error(getErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error("maestro test pod not found")
	}
}
