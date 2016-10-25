package lifecycle

import (
	"fmt"
	"log"
	"strings"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/manager"
)

func confDir(confPath string) string {
	confArr := strings.Split(confPath, "/")
	res := ""
	for i := range confArr[0 : len(confArr)-1] {
		if len(res) == 0 {
			res = confArr[i]
		} else {
			res = fmt.Sprintf("%s/%s", res, confArr[i])
		}
	}
	return res
}

// Check checks the running queue for an available spot for the next entry from the waiting queue
func Check(conf *config.Config, queue *Queue, running *Running, manager manager.Driver) error {
	log.Println("Checking for a job to run")
	log.Println("Queue: ", *queue)
	next := queue.Pop(running, conf.Server.MaxBuilds)
	if next != nil {
		log.Println("About to build this on: ", next)
		for i := range conf.Projects {
			log.Println(next.Project, conf.Projects[i].Name)
			if next.Project == conf.Projects[i].Name {
				shouldDeploy := false
				log.Println("Found a job to run")
				for j := range conf.Projects[i].DeployBranches {
					if next.Branch == conf.Projects[i].DeployBranches[j] {
						log.Println("Will Deploy")
						shouldDeploy = true
						break
					}
				}
				confPath := conf.Projects[i].MaestroConfPath
				log.Println("Running build")
				runErr := manager.Run(
					fmt.Sprintf("%s-%s", next.Project, next.Branch),
					confDir(confPath),
					confDir(confPath),
					[]string{
						"maestro",
						fmt.Sprintf("--branch=%s", next.Branch),
						fmt.Sprintf("--deploy=%v", shouldDeploy),
						fmt.Sprintf("--prev-commit=%s", next.PrevCommit),
						fmt.Sprintf("--curr-commit=%s", next.CurrCommit),
						fmt.Sprintf("--config=%s", confPath),
						fmt.Sprintf("--clone-path=%s", conf.Server.WorkspaceDir),
					},
				)
				return runErr
			}
		}
	}
	return nil
}
