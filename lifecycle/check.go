package lifecycle

import (
	"fmt"
	"strings"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/manager"
)

func confDir(confPath string) string {
	confArr := strings.Split(confPath, "/")
	var res string
	for i := range confArr[0 : len(confArr)-1] {
		res = fmt.Sprintf("%s%s", res, confArr[i])
	}
	return res
}

// Check checks the running queue for an available spot for the next entry from the waiting queue
func Check(conf *config.Config, queue *Queue, running *Running, manager manager.Driver, errChan chan error) {
	next := queue.Pop(running, conf.Server.MaxBuilds)
	if next != nil {
		for i := range conf.Projects {
			if next.Project == conf.Projects[i].Name {
				shouldDeploy := false
				for j := range conf.Projects[i].DeployBranches {
					if next.Branch == conf.Projects[i].DeployBranches[j] {
						shouldDeploy = true
						break
					}
				}
				confPath := conf.Projects[i].MaestroConfPath
				runErr := manager.Run(next.Project, confDir(confPath), confDir(confPath), []string{
					"maestro",
					fmt.Sprintf("--branch=%s", next.Branch),
					fmt.Sprintf("--deploy=%v", shouldDeploy),
					fmt.Sprintf("--prev-commit=%s", next.Commit),
					fmt.Sprintf("--config=%s", conf.Projects[i].MaestroConfPath),
					fmt.Sprintf("--clone-path=%s", conf.Server.WorkspaceDir),
				})
				errChan <- runErr
			}
		}
	}
}
